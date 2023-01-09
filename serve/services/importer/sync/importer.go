package sync

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bangwork/import-tools/serve/services/file"

	"github.com/bangwork/import-tools/serve/services/account"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services"
	"github.com/bangwork/import-tools/serve/services/cache"
	"github.com/bangwork/import-tools/serve/services/importer/resolve"
	"github.com/bangwork/import-tools/serve/services/importer/types"
	log2 "github.com/bangwork/import-tools/serve/services/log"
	"github.com/bangwork/import-tools/serve/utils"
	"github.com/bangwork/import-tools/serve/utils/timestamp"
)

const (
	sendPerDataCount  = 5000
	fileEndSign       = "FILE_END_SIGN\n"
	getLogPerCount    = 50
	getLogAtLastCount = 500
	retryCount        = 5
	retryIntervalSec  = 10
)

var (
	uploadAttachmentDoneChannel = make(chan bool)
	syncImportLogChannel        = make(chan bool)
)

type Importer struct {
	importTask       *types.ImportTask
	resourceUUID     string
	MapTagFile       map[string]*os.File
	importUUID       string
	logFile          *os.File
	mapCount         map[string]int64
	needCalculateTag map[string]bool
}

func (p *Importer) Resolve() error {
	startT := time.Now()
	log.Println("[start InitImportFile]")
	if err := p.initXmlPath(); err != nil {
		return err
	}
	resolver, err := InitImportFile(p.importTask)
	if err != nil {
		log.Println("init import file err", err)
		info, err := cache.GetCacheInfo()
		if err != nil {
			log.Println("get cache fail", err)
			return err
		}
		info.ResolveStatus = common.ResolveStatusFail
		if err := cache.SetCacheInfo(info); err != nil {
			log.Println("set cache fail", err)
			return err
		}
		log.Printf("create resolver fail:%s; stopResolveSignal:%t", err, services.StopResolveSignal)
		return err
	}
	if err := resolver.PrepareResolve(); err != nil {
		return err
	}
	log.Printf("[end InitImportFile] cost: %s", time.Since(startT))
	if err := resolver.Clear(); err != nil {
		return err
	}

	return nil
}

func (p *Importer) init() error {
	p.needCalculateTag = map[string]bool{
		services.ResourceTypeStringProject:           true,
		services.ResourceTypeStringTask:              true,
		services.ResourceTypeStringUser:              true,
		services.ResourceTypeStringTaskAttachmentTmp: true,
		services.ResourceTypeStringTaskAttachment:    true,
	}
	p.mapCount = make(map[string]int64)
	p.importUUID = utils.UUID()
	file, err := log2.CreateTeamLogFile(p.importTask.ImportTeamUUID, p.importUUID)
	if err != nil {
		return err
	}
	p.logFile = file
	info, err := cache.GetCacheInfo()
	if err != nil {
		return err
	}
	info.ImportUUID = p.importUUID
	if err = cache.SetCacheInfo(info); err != nil {
		return err
	}
	if err = p.initOutputFile(); err != nil {
		return err
	}
	return nil
}

func (p *Importer) Import() (err error) {
	defer func() {
		if err = p.setCache(err); err != nil {
			p.writeLog("set cache fail:%+v", err)
		}
		services.StopImportSignal = true
	}()
	startT := time.Now()
	p.writeLog("[start resolve]")
	if err = p.init(); err != nil {
		return err
	}
	resolver, err := createResolver(p.importTask)
	if err != nil {
		p.writeLog("create resolver fail:%+v", err)
		return err
	}

	p.runBatch(services.ResourceTypeStringUser, resolver.NextUser)
	p.runBatch(services.ResourceTypeStringUserGroup, resolver.NextUserGroup)
	p.runBatch(services.ResourceTypeStringUserGroupMember, resolver.NextUserGroupMember)
	p.runBatch(services.ResourceTypeStringGlobalProjectRole, resolver.NextGlobalProjectRole)
	p.runBatch(services.ResourceTypeStringGlobalProjectField, resolver.NextGlobalProjectField)
	p.runBatch(services.ResourceTypeStringIssueType, resolver.NextIssueType)
	p.runBatch(services.ResourceTypeStringProject, resolver.NextProject)
	p.runBatch(services.ResourceTypeStringProjectIssueType, resolver.NextProjectIssueType)
	p.runBatch(services.ResourceTypeStringProjectRole, resolver.NextProjectRole)
	p.runBatch(services.ResourceTypeStringProjectRoleMember, resolver.NextProjectRoleMember)
	p.runBatch(services.ResourceTypeStringGlobalPermission, resolver.NextGlobalPermission)
	p.runBatch(services.ResourceTypeStringProjectPermission, resolver.NextProjectPermission)
	p.runBatch(services.ResourceTypeStringProjectFieldValue, resolver.NextProjectFieldValue)
	p.runContinueBatch(services.ResourceTypeStringTaskStatus, resolver.NextTaskStatus)
	p.runBatch(services.ResourceTypeStringTaskField, resolver.NextTaskField)
	p.runBatch(services.ResourceTypeStringTaskFieldOption, resolver.NextTaskFieldOption)
	p.runBatch(services.ResourceTypeStringIssueTypeField, resolver.NextIssueTypeField)
	p.runBatch(services.ResourceTypeStringIssueTypeLayout, resolver.NextIssueTypeLayout)
	p.runBatch(services.ResourceTypeStringProjectIssueTypeField, resolver.NextProjectIssueTypeField)
	p.runBatch(services.ResourceTypeStringProjectIssueTypeLayout, resolver.NextProjectIssueTypeLayout)
	p.runBatch(services.ResourceTypeStringPriority, resolver.NextPriority)
	p.runBatch(services.ResourceTypeStringTaskLinkType, resolver.NextTaskLinkType)
	p.runBatch(services.ResourceTypeStringWorkflow, resolver.NextWorkflow)
	p.runBatch(services.ResourceTypeStringSprint, resolver.NextSprint)
	p.runBatch(services.ResourceTypeStringTask, resolver.NextTask)
	p.runContinueBatch(services.ResourceTypeStringTaskFieldValue, resolver.NextTaskFieldValue)
	p.runContinueBatch(services.ResourceTypeStringTaskWatcher, resolver.NextTaskWatcher)
	p.runContinueBatch(services.ResourceTypeStringTaskWorkLog, resolver.NextTaskWorkLog)
	p.runContinueBatch(services.ResourceTypeStringTaskComment, resolver.NextTaskComment)
	p.runBatch(services.ResourceTypeStringTaskRelease, resolver.NextTaskRelease)
	p.runBatch(services.ResourceTypeStringTaskLink, resolver.NextTaskLink)
	p.runBatch(services.ResourceTypeStringNotification, resolver.NextNotification)
	p.runContinueBatch(services.ResourceTypeStringTaskAttachmentTmp, resolver.NextTaskAttachment)
	p.runContinueBatch(services.ResourceTypeStringChangeItem, resolver.NextChangeItem)
	p.runOnce(services.ResourceTypeStringConfig, resolver.Config)

	p.setCountCache(services.ResourceTypeStringTaskAttachment, resolver.TotalAttachmentSize())
	p.writeLog("[end resolve] cost: %s", time.Since(startT))

	if err = p.outputFileWriteToRead(); err != nil {
		return err
	}
	go func() {
		if err := p.uploadAttachments(); err != nil {
			p.writeLog("[upload attachments error]: %+v", err)
		}
		uploadAttachmentDoneChannel <- true
	}()

	if err = p.sendStartImportNotice(); err != nil {
		return err
	}
	if err = p.sendData(); err != nil {
		return err
	}
	go func() {
		if err := p.syncImportLog(); err != nil {
			p.writeLog("[sync import log error]: %+v", err)
		}
		syncImportLogChannel <- true
	}()

	<-uploadAttachmentDoneChannel
	if err := p.sendAttachmentsData(); err != nil {
		return err
	}
	<-syncImportLogChannel

	p.clear()
	return nil
}

func (p *Importer) uploadAttachments() error {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[importer] err: %s\n%s", p, debug.Stack())
		}
	}()
	startTime := time.Now()
	p.writeLog("[start upload attachments]")
	reader := bufio.NewReader(p.MapTagFile[services.ResourceTypeStringTaskAttachmentTmp])
	for {
		if services.CheckIsStop() {
			p.writeLog("importer stop")
			break
		}
		readString, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		r := new(resolve.ThirdTaskAttachment)
		if err = json.Unmarshal([]byte(readString), &r); err != nil {
			return err
		}
		filePath := fmt.Sprintf("%s%s", p.importTask.AttachmentsPath, r.FilePath)
		fi, err := os.Open(filePath)
		if err != nil {
			return err
		}
		resourceUUID, err := file.UploadFile(fi, r.FileName)
		r.ResourceUUID = resourceUUID
		line := utils.OutputJSON(r)
		_, err = p.MapTagFile[services.ResourceTypeStringTaskAttachment].WriteString(string(line) + "\n")
		if err != nil {
			p.writeLog("write error:%+v", err)
			return err
		}
	}
	filePath := p.MapTagFile[services.ResourceTypeStringTaskAttachment].Name()
	if err := p.MapTagFile[services.ResourceTypeStringTaskAttachment].Close(); err != nil {
		fmt.Println("close file err", err)
		return err
	}
	openFile, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("open file err", err)
		return err
	}
	p.MapTagFile[services.ResourceTypeStringTaskAttachment] = openFile
	p.writeLog("[end upload attachments] cost: %s", time.Since(startTime))
	return nil
}

func (p *Importer) clear() {
	p.writeLogWithoutTime(fileEndSign)
	p.logFile.Sync()
	for k, _ := range p.MapTagFile {
		p.MapTagFile[k].Close()
	}
}

func (p *Importer) syncImportLog() error {
	defer func() {
		if err := recover(); err != nil {
			p.writeLog("sync import log panic: %+v", err)
		}
	}()
	startTimeSyncLog := time.Now()
	p.writeLog("[start sync import log]")
	accountInfo := new(account.Account)
	if err := accountInfo.Login(); err != nil {
		return err
	}
	retryFunc := func(accInfo *account.Account, perCount int) (length int) {
		logs, err := accInfo.GetImportLog(perCount)
		length = len(logs)
		if err == nil {
			p.writeLogWithoutTime(logs...)
			return
		}
		for j := 0; j < retryCount; j++ {
			p.writeLog("retry count: %d", j)
			if err := accInfo.Login(); err != nil {
				p.writeLog("login err:%+v", err)
				continue
			}
			logs, err = accInfo.GetImportLog(perCount)
			if err != nil {
				p.writeLogWithoutTime(logs...)
				continue
			}
			p.writeLogWithoutTime(logs...)
			length = len(logs)
			break
		}
		p.logFile.Sync()
		return
	}

	needUpdateStatus := map[string]bool{
		common.ImportStatusLabelDone:        true,
		common.ImportStatusLabelFail:        true,
		common.ImportStatusLabelInterrupted: true,
	}

	cacheInfo, err := cache.GetCacheInfo()
	if err != nil {
		p.writeLog("[get cache fail]: %+v", err)
		return err
	}
	for {
		retryFunc(accountInfo, getLogPerCount)
		importStatus, err := accountInfo.GetImportStatus(p.resourceUUID)
		if err != nil {
			log.Println("get import status err", err)
			break
		}
		switch importStatus {
		case common.ImportStatusLabelDone:
			cacheInfo.ImportResult.Status = common.ImportStatusDone
		case common.ImportStatusLabelFail:
			cacheInfo.ImportResult.Status = common.ImportStatusFail
		case common.ImportStatusLabelInterrupted:
			cacheInfo.ImportResult.Status = common.ImportStatusCancel
			services.StopImportSignal = true
		}
		if needUpdateStatus[importStatus] {
			if err := cache.SetCacheInfo(cacheInfo); err != nil {
				p.writeLog("[set cache fail]: %+v", err)
			}
			retryFunc(accountInfo, getLogAtLastCount)
			break
		}
		time.Sleep(retryIntervalSec * time.Second)
	}

	p.writeLog("[end sync import log] cost: %s", time.Since(startTimeSyncLog))
	return nil
}

func (p *Importer) sendStartImportNotice() error {
	startTime := time.Now()
	p.writeLog("[start sendStartImportNotice]")
	accountInfo := new(account.Account)
	if err := accountInfo.Login(); err != nil {
		return err
	}
	info, err := file.PrepareUploadInfo(accountInfo.Cache.BackupName, file.LabelJiraImport, accountInfo)
	if err != nil {
		p.writeLog("prepare upload info err:%s", err)
		return err
	}
	p.resourceUUID = info.ResourceUUID
	if err = accountInfo.SetPassword(p.resourceUUID, p.importTask.Password); err != nil {
		p.writeLog("set password err:%s", err)
		return err
	}
	confirmBody := new(services.ConfirmImportRequest)
	confirmBody.ImportUUID = p.importUUID
	confirmBody.Version = accountInfo.Cache.ResolveResult.JiraVersion
	confirmBody.ServerID = accountInfo.Cache.ResolveResult.JiraServerID
	confirmBody.ResourceUUID = p.resourceUUID
	confirmBody.ImportType = common.ImportTypeImportTools
	confirmBody.FromImportTools = true
	batchTaskUUID, err := accountInfo.ConfirmImport(confirmBody)
	if err != nil {
		p.writeLog("confirm import err:%s", err)
		return err
	}
	services.ImportBatchTaskUUID = batchTaskUUID
	p.writeLog("[end sendStartImportNotice] cost: %s", time.Since(startTime))
	return nil
}

func (p *Importer) sendData() error {
	startTime := time.Now()
	p.writeLog("[start send data]")
	accountInfo := new(account.Account)
	if err := accountInfo.Login(); err != nil {
		return err
	}

	retryFunc := func(resourceTypeString string, data string) error {
		if err := accountInfo.SendImportData(resourceTypeString, data); err == nil {
			p.writeLog("send import data succ: %s", resourceTypeString)
			return nil
		}
		for j := 0; j < retryCount; j++ {
			p.writeLog("retry count: %d", j)
			if err := accountInfo.Login(); err != nil {
				p.writeLog("login err:%+v", err)
				continue
			}
			if err := accountInfo.SendImportData(resourceTypeString, data); err != nil {
				log.Println("send import data err", err)
				continue
			}
			break
		}
		return nil
	}

	for _, tag := range services.MapOutputFile {
		if tag == services.ResourceTypeStringTaskAttachmentTmp || tag == services.ResourceTypeStringTaskAttachment {
			continue
		}
		if services.CheckIsStop() {
			p.writeLog("importer stop")
			break
		}
		fi := p.MapTagFile[tag]
		fileScanner := bufio.NewScanner(fi)
		i := 0
		curData := strings.Builder{}

		for fileScanner.Scan() {
			line := fileScanner.Bytes()
			if len(line) == 0 {
				break
			}
			curData.Write(line)
			curData.Write([]byte("\n"))
			i++
			if i >= sendPerDataCount {
				if err := retryFunc(tag, curData.String()); err != nil {
					p.writeLog("send import data err: %+v", err)
					return err
				}
				curData = strings.Builder{}
				i = 0
			}
		}
		curData.Write([]byte(fileEndSign))
		if err := retryFunc(tag, curData.String()); err != nil {
			log.Println("send import data err", err)
			return err
		}
		p.MapTagFile[tag].Close()
		p.logFile.Sync()
	}
	p.writeLog("[end send data] cost: %s", time.Since(startTime))
	return nil
}

func (p *Importer) sendAttachmentsData() error {
	startTime := time.Now()
	p.writeLog("[start send attachments]")
	accountInfo := new(account.Account)
	if err := accountInfo.Login(); err != nil {
		return err
	}

	retryFunc := func(resourceTypeString string, data string) error {
		if err := accountInfo.SendImportData(resourceTypeString, data); err == nil {
			p.writeLog("send import data succ: %s", resourceTypeString)
			return nil
		}
		for j := 0; j < retryCount; j++ {
			p.writeLog("retry count: %d", j)
			if err := accountInfo.Login(); err != nil {
				p.writeLog("login err:%+v", err)
				continue
			}
			if err := accountInfo.SendImportData(resourceTypeString, data); err != nil {
				log.Println("send import data err", err)
				continue
			}
			break
		}
		return nil
	}

	tag := services.ResourceTypeStringTaskAttachment

	fi := p.MapTagFile[tag]
	fileScanner := bufio.NewScanner(fi)
	i := 0
	curData := strings.Builder{}

	for fileScanner.Scan() {
		line := fileScanner.Bytes()
		if len(line) == 0 {
			break
		}
		if services.CheckIsStop() {
			p.writeLog("importer stop")
			break
		}
		curData.Write(line)
		curData.Write([]byte("\n"))
		i++
		if i >= sendPerDataCount {
			if err := retryFunc(tag, curData.String()); err != nil {
				p.writeLog("send import data err: %+v", err)
				return err
			}
			curData = strings.Builder{}
			i = 0
		}
	}
	curData.Write([]byte(fileEndSign))
	if err := retryFunc(tag, curData.String()); err != nil {
		log.Println("send import data err", err)
		return err
	}
	p.MapTagFile[tag].Close()
	p.logFile.Sync()
	p.writeLog("[end send attachments] cost: %s", time.Since(startTime))
	return nil
}

func (p *Importer) outputFileWriteToRead() error {
	startTime := time.Now()
	p.writeLog("[start change file mode]")
	var err error
	for tag, fi := range p.MapTagFile {
		filePath := fi.Name()
		if err = p.MapTagFile[tag].Close(); err != nil {
			return err
		}
		p.MapTagFile[tag], err = os.OpenFile(filePath, os.O_RDWR, 0666)
		if err != nil {
			return err
		}
	}
	p.writeLog("[end change file mode] cost: %s", time.Since(startTime))
	return nil
}

func (p *Importer) setCache(importErr error) error {
	cacheInfo, e := cache.GetCacheInfo()
	if e != nil {
		return e
	}
	status := common.ImportStatusDone
	if importErr != nil {
		log.Println("import fail", importErr)
		status = common.ImportStatusFail
	}
	if cacheInfo.ImportResult.Status != common.ImportStatusCancel {
		cacheInfo.ImportResult.Status = status
	}
	cacheInfo.ImportResult.DoneTime = time.Now().Unix()
	return cache.SetCacheInfo(cacheInfo)
}

func (p *Importer) setCountCache(tag string, count int64) error {
	cacheInfo, e := cache.GetCacheInfo()
	if e != nil {
		return e
	}
	switch tag {
	case services.ResourceTypeStringProject:
		cacheInfo.ImportScope.ProjectCount = count
	case services.ResourceTypeStringTask:
		cacheInfo.ImportScope.IssueCount = count
	case services.ResourceTypeStringUser:
		cacheInfo.ImportScope.MemberCount = count
	case services.ResourceTypeStringTaskAttachmentTmp:
		cacheInfo.ImportScope.AttachmentCount = count
	case services.ResourceTypeStringTaskAttachment:
		cacheInfo.ImportScope.AttachmentSize = count
	}
	if err := cache.SetCacheInfo(cacheInfo); err != nil {
		return err
	}
	cache.SetExpectTimeCache()
	return nil
}

func (p *Importer) runBatch(tag string, f func() ([]byte, error)) {
	p.writeLog("start resolve: %s", tag)
	for {
		if services.CheckIsStop() {
			p.writeLog("stop signal: %s", tag)
			log.Println("importer stop")
			break
		}
		line, err := f()
		if err != nil {
			panic(err)
		}
		if len(line) == 0 {
			break
		}
		_, err = p.MapTagFile[tag].WriteString(string(line) + "\n")
		if err != nil {
			log.Panicf("write error:%+v", err)
			break
		}
		p.mapCount[tag]++
	}
	if p.needCalculateTag[tag] {
		p.setCountCache(tag, p.mapCount[tag])
	}
	p.MapTagFile[tag].Sync()
	p.writeLog("end resolve: %s", tag)
}

func (p *Importer) runContinueBatch(tag string, f func() ([]byte, bool, error)) {
	p.writeLog("start resolve: %s", tag)
	for {
		if services.CheckIsStop() {
			p.writeLog("stop signal: %s", tag)
			log.Println("importer stop")
			break
		}
		line, continueSignal, err := f()
		if err != nil {
			panic(err)
		}
		if continueSignal {
			continue
		}
		if len(line) == 0 {
			break
		}
		_, err = p.MapTagFile[tag].WriteString(string(line) + "\n")
		if err != nil {
			log.Panicf("write error:%+v", err)
			break
		}
		p.mapCount[tag]++
	}
	if p.needCalculateTag[tag] {
		p.setCountCache(tag, p.mapCount[tag])
	}
	p.MapTagFile[tag].Sync()
}

func (p *Importer) runOnce(tag string, f func() ([]byte, error)) {
	if services.CheckIsStop() {
		p.writeLog("stop signal: %s", tag)
		return
	}
	p.writeLog("start resolve: %s", tag)
	line, err := f()
	if err != nil {
		panic(err)
	}
	if len(line) == 0 {
		return
	}
	_, err = p.MapTagFile[tag].WriteString(string(line) + "\n")
	if err != nil {
		log.Panicf("write error:%+v", err)
		return
	}
	p.mapCount[tag]++
	p.MapTagFile[tag].Sync()
	p.writeLog("end resolve: %s", tag)
}

func NewImporter(task *types.ImportTask) *Importer {
	return &Importer{importTask: task}
}

func (p *Importer) writeLog(input string, param ...interface{}) {
	prefix := fmt.Sprintf("[%s] %s\n", timestamp.NowTimeString(), input)
	msg := fmt.Sprintf(prefix, param...)
	if len(param) > 0 {
		firstParam, ok := param[0].(error)
		if ok && firstParam != nil {
			msg += fmt.Sprintf("%s:%s", firstParam, debug.Stack())
		}
	}
	fmt.Println(msg)
	p.logFile.WriteString(msg)
}

func (p *Importer) writeLogWithoutTime(inputs ...string) {
	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}
		msg := fmt.Sprintf("%s\n", input)
		fmt.Println(msg)
		p.logFile.WriteString(msg)
	}
}

func (p *Importer) initOutputFile() error {
	outputPath := fmt.Sprintf("%s/%s", common.Path, common.OutputDir)
	if err := p.initPath(outputPath); err != nil {
		return err
	}
	info, err := cache.GetCacheInfo()
	if err != nil {
		return err
	}
	p.MapTagFile = make(map[string]*os.File)
	for _, tag := range services.MapOutputFile {
		path := fmt.Sprintf("%s/%s.json", outputPath, tag)
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		p.MapTagFile[tag] = file
		//services.MapOutputFile[tag] = path
	}
	//info.MapOutputFilePath = services.MapOutputFile
	return cache.SetCacheInfo(info)
}
func (p *Importer) initPath(outputPath string) error {
	if utils.CheckPathExist(outputPath) {
		stat, err := os.Stat(outputPath)
		if err != nil {
			return err
		}
		modTime := stat.ModTime().Format("2006-01-02_15:04")
		backPath := fmt.Sprintf("%s-%s-%s", outputPath, modTime, utils.RandomNumberString(5))
		if err = os.Rename(outputPath, backPath); err != nil {
			return err
		}
	}
	err := os.MkdirAll(outputPath, 0755)
	if err != nil {
		return err
	}
	return nil
}
func (p *Importer) initXmlPath() error {
	xmlPath := fmt.Sprintf("%s/%s", common.Path, common.XmlDir)
	return p.initPath(xmlPath)
}
