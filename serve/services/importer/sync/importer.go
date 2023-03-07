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

	"github.com/bangwork/import-tools/serve/models/ones"

	"github.com/bangwork/import-tools/serve/services/file"

	"github.com/bangwork/import-tools/serve/services/account"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services"
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
	importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
	resolver, err := InitImportFile(p.importTask)
	if err != nil {
		log.Println("init import file err", err)
		importCache.ResolveStatus = common.ResolveStatusFail
		common.ImportCacheMap.Set(p.importTask.Cookie, importCache)
		log.Printf("create resolver fail:%s; stopResolveSignal:%t", err, importCache.StopResolveSignal)
		return err
	}
	if err := resolver.PrepareResolve(); err != nil {
		return err
	}
	log.Printf("[end InitImportFile] cost: %s", time.Since(startT))
	resolver.Clear()

	return nil
}

func (p *Importer) init() error {
	p.needCalculateTag = map[string]bool{
		common.ResourceTypeStringProject:           true,
		common.ResourceTypeStringTask:              true,
		common.ResourceTypeStringUser:              true,
		common.ResourceTypeStringTaskAttachmentTmp: true,
		common.ResourceTypeStringTaskAttachment:    true,
	}
	p.mapCount = make(map[string]int64)
	p.importUUID = utils.UUID()
	file, err := log2.CreateTeamLogFile(p.importTask.ImportTeamUUID, p.importUUID)
	if err != nil {
		return err
	}
	p.logFile = file
	info, err := common.GetCacheInfo(p.importTask.Key)
	if err != nil {
		return err
	}
	info.ImportUUID = p.importUUID
	if err = common.SetCacheInfo(p.importTask.Key, info); err != nil {
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
		importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
		importCache.StopImportSignal = true
		common.ImportCacheMap.Set(p.importTask.Cookie, importCache)
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

	p.runBatch(common.ResourceTypeStringUser, resolver.NextUser)
	p.runBatch(common.ResourceTypeStringUserGroup, resolver.NextUserGroup)
	p.runBatch(common.ResourceTypeStringUserGroupMember, resolver.NextUserGroupMember)
	p.runBatch(common.ResourceTypeStringGlobalProjectRole, resolver.NextGlobalProjectRole)
	p.runBatch(common.ResourceTypeStringGlobalProjectField, resolver.NextGlobalProjectField)
	p.runBatch(common.ResourceTypeStringIssueType, resolver.NextIssueType)
	p.runBatch(common.ResourceTypeStringProject, resolver.NextProject)
	p.runBatch(common.ResourceTypeStringProjectIssueType, resolver.NextProjectIssueType)
	p.runBatch(common.ResourceTypeStringProjectRole, resolver.NextProjectRole)
	p.runBatch(common.ResourceTypeStringProjectRoleMember, resolver.NextProjectRoleMember)
	p.runBatch(common.ResourceTypeStringGlobalPermission, resolver.NextGlobalPermission)
	p.runBatch(common.ResourceTypeStringProjectPermission, resolver.NextProjectPermission)
	p.runBatch(common.ResourceTypeStringProjectFieldValue, resolver.NextProjectFieldValue)
	p.runContinueBatch(common.ResourceTypeStringTaskStatus, resolver.NextTaskStatus)
	p.runBatch(common.ResourceTypeStringTaskField, resolver.NextTaskField)
	p.runBatch(common.ResourceTypeStringTaskFieldOption, resolver.NextTaskFieldOption)
	p.runBatch(common.ResourceTypeStringIssueTypeField, resolver.NextIssueTypeField)
	p.runBatch(common.ResourceTypeStringIssueTypeLayout, resolver.NextIssueTypeLayout)
	p.runBatch(common.ResourceTypeStringProjectIssueTypeField, resolver.NextProjectIssueTypeField)
	p.runBatch(common.ResourceTypeStringProjectIssueTypeLayout, resolver.NextProjectIssueTypeLayout)
	p.runBatch(common.ResourceTypeStringPriority, resolver.NextPriority)
	p.runBatch(common.ResourceTypeStringTaskLinkType, resolver.NextTaskLinkType)
	p.runBatch(common.ResourceTypeStringWorkflow, resolver.NextWorkflow)
	p.runBatch(common.ResourceTypeStringSprint, resolver.NextSprint)
	p.runBatch(common.ResourceTypeStringTask, resolver.NextTask)
	p.runContinueBatch(common.ResourceTypeStringTaskFieldValue, resolver.NextTaskFieldValue)
	p.runContinueBatch(common.ResourceTypeStringTaskWatcher, resolver.NextTaskWatcher)
	p.runContinueBatch(common.ResourceTypeStringTaskWorkLog, resolver.NextTaskWorkLog)
	p.runContinueBatch(common.ResourceTypeStringTaskComment, resolver.NextTaskComment)
	p.runBatch(common.ResourceTypeStringTaskRelease, resolver.NextTaskRelease)
	p.runBatch(common.ResourceTypeStringTaskLink, resolver.NextTaskLink)
	p.runBatch(common.ResourceTypeStringNotification, resolver.NextNotification)
	p.runContinueBatch(common.ResourceTypeStringTaskAttachmentTmp, resolver.NextTaskAttachment)
	p.runContinueBatch(common.ResourceTypeStringChangeItem, resolver.NextChangeItem)
	p.runOnce(common.ResourceTypeStringConfig, resolver.Config)

	p.setCountCache(common.ResourceTypeStringTaskAttachment, resolver.TotalAttachmentSize())
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
	reader := bufio.NewReader(p.MapTagFile[common.ResourceTypeStringTaskAttachmentTmp])

	if err := ones.LoginONESAndSetAuth(p.importTask.Cookie); err != nil {
		p.writeLog("LoginONESAndSetAuth err:%+v", err)
		return err
	}
	for {
		importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
		if importCache.CheckIsStop() {
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
		resourceUUID, err := file.UploadFile(p.importTask.Cookie, p.importTask.ImportTeamUUID, fi, r.FileName)
		if err != nil || resourceUUID == "" {
			p.writeLog("upload file err: %+v, %s", err, resourceUUID)
			for j := 0; j < retryCount; j++ {
				p.writeLog("upload file retry count: %d", j)
				if err := ones.LoginONESAndSetAuth(p.importTask.Cookie); err != nil {
					p.writeLog("upload file login err:%+v", err)
					continue
				}
				resourceUUID, err = file.UploadFile(p.importTask.Cookie, p.importTask.ImportTeamUUID, fi, r.FileName)
				if err != nil || resourceUUID == "" {
					log.Println("upload file err", err, resourceUUID)
					continue
				}
				break
			}
		}
		fi.Close()
		if err != nil {
			p.writeLog("upload file error: %+v", err)
			continue
		}
		if resourceUUID == "" {
			p.writeLog("resource uuid empty", r.FileName)
			continue
		}
		r.ResourceUUID = resourceUUID
		line := utils.OutputJSON(r)
		_, err = p.MapTagFile[common.ResourceTypeStringTaskAttachment].WriteString(string(line) + "\n")
		if err != nil {
			p.writeLog("write error:%+v", err)
			return err
		}
	}
	filePath := p.MapTagFile[common.ResourceTypeStringTaskAttachment].Name()
	if err := p.MapTagFile[common.ResourceTypeStringTaskAttachment].Close(); err != nil {
		fmt.Println("close file err", err)
		return err
	}
	openFile, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("open file err", err)
		return err
	}
	p.MapTagFile[common.ResourceTypeStringTaskAttachment] = openFile
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
	if err := ones.LoginONESAndSetAuth(p.importTask.Cookie); err != nil {
		p.writeLog("LoginONESAndSetAuth err%+v", err)
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
			if err := ones.LoginONESAndSetAuth(p.importTask.Cookie); err != nil {
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

	for {
		importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
		retryFunc(accountInfo, getLogPerCount)
		importStatus, err := accountInfo.GetImportStatus(p.resourceUUID)
		if err != nil {
			log.Println("get import status err", err)
			break
		}
		switch importStatus {
		case common.ImportStatusLabelDone:
			importCache.ImportResult.Status = common.ImportStatusDone
		case common.ImportStatusLabelFail:
			importCache.ImportResult.Status = common.ImportStatusFail
		case common.ImportStatusLabelInterrupted:
			importCache.ImportResult.Status = common.ImportStatusCancel
			importCache.StopImportSignal = true
		}
		if needUpdateStatus[importStatus] {
			common.ImportCacheMap.Set(p.importTask.Cookie, importCache)
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
	if err := ones.LoginONESAndSetAuth(p.importTask.Cookie); err != nil {
		p.writeLog("LoginONESAndSetAuth err%+v", err)
		return err
	}
	cookieValue, err := ones.DecryptCookieValueByCookie(p.importTask.Cookie)
	if err != nil {
		return err
	}

	info, err := file.PrepareUploadInfo(p.importTask.BackupName, file.LabelJiraImport, file.EntityTypeJiraLabel, p.importTask.URL, p.importTask.ImportTeamUUID, cookieValue.GenAuthHeader())
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
	importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
	importCache.ImportBatchTaskUUID = batchTaskUUID
	common.ImportCacheMap.Set(p.importTask.Cookie, importCache)
	p.writeLog("[end sendStartImportNotice] cost: %s", time.Since(startTime))
	return nil
}

func (p *Importer) sendData() error {
	startTime := time.Now()
	p.writeLog("[start send data]")
	accountInfo := new(account.Account)
	if err := ones.LoginONESAndSetAuth(p.importTask.Cookie); err != nil {
		p.writeLog("LoginONESAndSetAuth err%+v", err)
		return err
	}

	retryFunc := func(resourceTypeString string, data string) error {
		if err := accountInfo.SendImportData(resourceTypeString, data); err == nil {
			p.writeLog("send import data succ: %s", resourceTypeString)
			return nil
		}
		for j := 0; j < retryCount; j++ {
			p.writeLog("retry count: %d", j)
			if err := ones.LoginONESAndSetAuth(p.importTask.Cookie); err != nil {
				p.writeLog("LoginONESAndSetAuth err%+v", err)
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

	for _, tag := range common.MapOutputFile {
		if tag == common.ResourceTypeStringTaskAttachmentTmp || tag == common.ResourceTypeStringTaskAttachment {
			continue
		}
		importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
		if importCache.CheckIsStop() {
			p.writeLog("importer stop")
			break
		}
		fi := p.MapTagFile[tag]
		fileScanner := bufio.NewScanner(fi)
		fileScanner.Buffer([]byte{}, common.GetMaxScanTokenSize())
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
		if fileScanner.Err() != nil {
			p.writeLog("scan err: %v\n", fileScanner.Err())
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
	if err := ones.LoginONESAndSetAuth(p.importTask.Cookie); err != nil {
		p.writeLog("login err:%+v", err)
		return err
	}

	retryFunc := func(resourceTypeString string, data string) error {
		if err := accountInfo.SendImportData(resourceTypeString, data); err == nil {
			p.writeLog("send import data succ: %s", resourceTypeString)
			return nil
		}
		for j := 0; j < retryCount; j++ {
			p.writeLog("retry count: %d", j)
			if err := ones.LoginONESAndSetAuth(p.importTask.Cookie); err != nil {
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

	tag := common.ResourceTypeStringTaskAttachment

	fi := p.MapTagFile[tag]
	fileScanner := bufio.NewScanner(fi)
	fileScanner.Buffer([]byte{}, common.GetMaxScanTokenSize())
	i := 0
	curData := strings.Builder{}

	for fileScanner.Scan() {
		line := fileScanner.Bytes()
		if len(line) == 0 {
			break
		}
		importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
		if importCache.CheckIsStop() {
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
	if fileScanner.Err() != nil {
		p.writeLog("scan file err: %+v", fileScanner.Err())
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
	cacheInfo, e := common.GetCacheInfo(p.importTask.Key)
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
	return common.SetCacheInfo(p.importTask.Key, cacheInfo)
}

func (p *Importer) setCountCache(tag string, count int64) error {
	cacheInfo, e := common.GetCacheInfo(p.importTask.Key)
	if e != nil {
		return e
	}
	switch tag {
	case common.ResourceTypeStringProject:
		cacheInfo.ImportScope.ProjectCount = count
	case common.ResourceTypeStringTask:
		cacheInfo.ImportScope.IssueCount = count
	case common.ResourceTypeStringUser:
		cacheInfo.ImportScope.MemberCount = count
	case common.ResourceTypeStringTaskAttachmentTmp:
		cacheInfo.ImportScope.AttachmentCount = count
	case common.ResourceTypeStringTaskAttachment:
		cacheInfo.ImportScope.AttachmentSize = count
	}
	if err := common.SetCacheInfo(p.importTask.Key, cacheInfo); err != nil {
		return err
	}
	common.SetExpectTimeCache(p.importTask.Key)
	return nil
}

func (p *Importer) runBatch(tag string, f func() ([]byte, error)) {
	p.writeLog("start resolve: %s", tag)
	for {
		importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
		if importCache.CheckIsStop() {
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
		importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
		if importCache.CheckIsStop() {
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
	importCache := common.ImportCacheMap.Get(p.importTask.Cookie)
	if importCache.CheckIsStop() {
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
	outputPath := fmt.Sprintf("%s/%s", common.GetCachePath(), common.OutputDir)
	if err := p.initPath(outputPath); err != nil {
		return err
	}
	info, err := common.GetCacheInfo(p.importTask.Key)
	if err != nil {
		return err
	}
	p.MapTagFile = make(map[string]*os.File)
	for _, tag := range common.MapOutputFile {
		path := fmt.Sprintf("%s/%s.json", outputPath, tag)
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		p.MapTagFile[tag] = file
		//services.MapOutputFile[tag] = path
	}
	//info.MapOutputFilePath = services.MapOutputFile
	return common.SetCacheInfo(p.importTask.Key, info)
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
	xmlPath := fmt.Sprintf("%s/%s", common.GetCachePath(), common.XmlDir)
	return p.initPath(xmlPath)
}
