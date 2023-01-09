package jira

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bangwork/import-tools/serve/common"

	"github.com/bangwork/import-tools/serve/services"

	"github.com/bangwork/import-tools/serve/utils/timestamp"

	jira "github.com/bangwork/import-tools/serve/services/importer/resolve"
	"github.com/bangwork/import-tools/serve/services/importer/types"
)

type handlerObjectFile struct {
	ImportTask    *types.ImportTask
	Reader        io.ReadCloser
	file          *os.File
	began         bool
	nowTimeString string
}

func processedObjectFile(importTask *types.ImportTask, reader io.ReadCloser) (string, error) {
	handler := &handlerObjectFile{
		ImportTask:    importTask,
		Reader:        reader,
		nowTimeString: timestamp.NowTimeString(),
	}
	err := handler.scan()
	if err != nil {
		return "", err
	}
	return handler.fileName(), nil
}

func (o *handlerObjectFile) createFile() error {
	root := common.Path
	root = fmt.Sprintf("%s/%s", root, common.XmlDir)
	err := os.MkdirAll(root, 0755)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("objects_*.xml")
	f, err := ioutil.TempFile(root, path)
	if err != nil {
		return err
	}
	o.file = f
	return nil
}

func (o *handlerObjectFile) scan() error {
	if err := o.createFile(); err != nil {
		return err
	}
	defer o.closeFile()

	log.Println("[jira import] start object xml scanner")
	scanner := jira.NewXmlScanner(o.Reader, activeObjectRootTag)
	for {
		e := scanner.NextElement()
		if e == nil {
			break
		}
		if services.StopResolveSignal {
			return common.Errors(common.StopResolve, nil)
		}
		if e.Tag == "data" {
			tableName, found := e.AttrMap["tableName"]
			if !found {
				continue
			}
			if strings.Contains(tableName, "_SPRINT") && !strings.Contains(tableName, "EX_SPRINT") {
				line := e.Encode()
				o.file.WriteString(line + "\n")
				break
			}
		}
	}
	log.Println("[jira import] end object xml scanner")
	return nil
}

func (o *handlerObjectFile) fileName() string {
	return o.file.Name()
}

func (o *handlerObjectFile) closeFile() error {
	return o.file.Close()
}
