package jira

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/juju/errors"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services"
	"github.com/bangwork/import-tools/serve/services/cache"
	jira "github.com/bangwork/import-tools/serve/services/importer/resolve"
	"github.com/bangwork/import-tools/serve/services/importer/types"
	"github.com/bangwork/import-tools/serve/utils/timestamp"
)

var (
	handleTags = map[string]struct{}{
		"User":                        {},
		"ApplicationUser":             {},
		"Group":                       {},
		"Membership":                  {},
		"Status":                      {},
		"ProjectRole":                 {},
		"Project":                     {},
		"ProjectRoleActor":            {},
		"ProjectCategory":             {},
		"CustomField":                 {},
		"Component":                   {},
		"CustomFieldOption":           {},
		"Resolution":                  {},
		"Priority":                    {},
		"IssueLinkType":               {},
		"IssueType":                   {},
		"ConfigurationContext":        {},
		"OptionConfiguration":         {},
		"Workflow":                    {},
		"WorkflowSchemeEntity":        {},
		"NodeAssociation":             {},
		"Version":                     {},
		"IssueLink":                   {},
		"Issue":                       {},
		"CustomFieldValue":            {},
		"Label":                       {},
		"EventType":                   {},
		"SchemePermissions":           {},
		"Action":                      {},
		"UserAssociation":             {},
		"Worklog":                     {},
		"GlobalPermissionEntry":       {},
		"FileAttachment":              {},
		"FieldLayout":                 {},
		"FieldLayoutSchemeEntity":     {},
		"FieldScreenScheme":           {},
		"FieldScreenSchemeItem":       {},
		"FieldScreenTab":              {},
		"IssueTypeScreenSchemeEntity": {},
		"FieldLayoutItem":             {},
		"FieldScreenLayoutItem":       {},
		"FieldConfigSchemeIssueType":  {},
		"FieldConfigScheme":           {},
		"ChangeGroup":                 {},
		"ChangeItem":                  {},
		"Notification":                {},
	}
)

type handlerEntityFile struct {
	ImportMessage        *types.ImportTask
	Tags                 map[string]struct{}
	Reader               io.ReadCloser
	hoursPerDay          string
	daysPerWeek          string
	began                bool
	beganMap             map[string]bool
	entityFile           *os.File
	entityFilePathMap    map[string]string
	entityFileMap        map[string]*os.File
	entityFileScannerMap map[string]*jira.XmlScanner
	entityTmpFile        *os.File
	nowTimeString        string
	ResolveResult        cache.ResolveResult
}

func processedEntityFile(msg *types.ImportTask, reader io.ReadCloser) (map[string]*jira.XmlScanner, map[string]string, *handlerEntityFile, error) {
	handler := newHandlerEntityFile(msg, handleTags, reader)
	err := handler.scan()
	if err != nil {
		return nil, nil, handler, err
	}
	fileList, mapFilePath := handler.fileLists()
	return fileList, mapFilePath, handler, nil
}

func newHandlerEntityFile(msg *types.ImportTask, tags map[string]struct{}, reader io.ReadCloser) *handlerEntityFile {
	o := &handlerEntityFile{
		ImportMessage:        msg,
		Tags:                 tags,
		Reader:               reader,
		beganMap:             map[string]bool{},
		entityFileScannerMap: make(map[string]*jira.XmlScanner, 0),
		nowTimeString:        timestamp.NowTimeString(),
	}
	return o
}

func (o *handlerEntityFile) fileLists() (map[string]*jira.XmlScanner, map[string]string) {
	o.entityFileMap = make(map[string]*os.File)
	o.entityFileScannerMap = make(map[string]*jira.XmlScanner)
	var mapTagFilePath = make(map[string]string)
	for tag, file := range o.entityFilePathMap {
		fi, err := os.Open(file)
		if err != nil {
			log.Printf("open file fail: %s", err)
		}
		scanner := jira.NewXmlScanner(fi, entityRootTag)
		o.entityFileScannerMap[tag] = scanner
		mapTagFilePath[tag] = fi.Name()
	}
	return o.entityFileScannerMap, mapTagFilePath
}

func importerFileDir() string {
	root := common.Path
	path := fmt.Sprintf("%s/%s", root, common.XmlDir)
	return path
}

func (o *handlerEntityFile) createFile() error {
	root := importerFileDir()
	err := os.MkdirAll(root, 0755)
	if err != nil {
		return err
	}
	o.entityFilePathMap = make(map[string]string)
	o.entityFileMap = make(map[string]*os.File)
	for tag, _ := range o.Tags {
		path := fmt.Sprintf("%s_*.xml", tag)
		f, err := ioutil.TempFile(root, path)
		if err != nil {
			return err
		}
		o.entityFileMap[tag] = f
		o.entityFilePathMap[tag] = f.Name()
	}

	tmpPath := fmt.Sprintf("entities_tmp_*.xml")
	o.entityTmpFile, err = ioutil.TempFile(root, tmpPath)

	return nil
}

func (o *handlerEntityFile) closeFile() error {
	for tag, _ := range handleTags {
		err := o.entityFileMap[tag].Close()
		if err != nil {
			return err
		}
	}
	tmpPath := o.entityTmpFile.Name()
	err := os.Remove(tmpPath)
	if err != nil {
		return err
	}
	return nil
}

func printOnly(r rune) rune {
	// uint32(r) == 10 is line breaks \n
	if unicode.IsPrint(r) || uint32(r) == 10 {
		return r
	}
	return -1
}

func (o *handlerEntityFile) scan() error {
	if err := o.createFile(); err != nil {
		return err
	}
	defer o.closeFile()

	log.Println("[jira import] start tmp xml scanner")
	tmpScanner := bufio.NewScanner(o.Reader)
	tmpScanner.Buffer([]byte{}, bufio.MaxScanTokenSize*10)

	for tmpScanner.Scan() {
		if services.StopResolveSignal {
			return common.Errors(common.StopResolve, nil)
		}
		line := tmpScanner.Text()

		data := strings.Map(printOnly, line)
		if _, e := o.entityTmpFile.WriteString(data + "\n"); e != nil {
			log.Printf("write line %s to tmpFile fail.%+v", line, e)
			return errors.Trace(e)
		}
	}

	log.Println("[jira import] end tmp xml scanner")
	tmpFilePath := o.entityTmpFile.Name()
	o.entityTmpFile.Close()
	var err error
	o.Reader, err = os.Open(tmpFilePath)
	if err != nil {
		return err
	}

	logMap := make(map[string]string)

	log.Println("[jira import] start xml scanner")
	scanner := jira.NewXmlScanner(o.Reader, entityRootTag)
	for {
		if services.StopResolveSignal {
			return common.Errors(common.StopResolve, nil)
		}
		e := scanner.NextElement()
		if e == nil {
			break
		}
		if _, found := o.Tags[e.Tag]; !found {
			continue
		}

		_, ok := logMap[e.Tag]
		if !ok {
			logMap[e.Tag] = ""
			log.Println(e.Tag)
		}

		line := e.Encode()
		if o.entityFileMap[e.Tag] == nil {
			log.Printf("file os closed: %s", e.Tag)
			return errors.Errorf("file os closed:%s", e.Tag)
		}
		o.entityFileMap[e.Tag].WriteString(line + "\n")
	}
	log.Println("[jira import] end xml scanner")
	o.ResolveResult = scanner.ResolveResult
	o.daysPerWeek = scanner.DaysPerWeek
	o.hoursPerDay = scanner.HoursPerDay
	return nil
}
