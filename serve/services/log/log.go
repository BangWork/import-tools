package log

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

const Log = "logs"
const PerCount = 50

func InitLogDir() error {
	filePath := logRootPath()
	if utils.CheckPathExist(filePath) {
		return nil
	}
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return err
	}
	return nil
}

func InitTeamLogDir(teamUUID string) error {
	path := logTeamPath(teamUUID)
	if utils.CheckPathExist(path) {
		return nil
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

func CreateTeamLogFile(teamUUID, importUUID string) (*os.File, error) {
	filePath := logFilePath(teamUUID, importUUID)
	if utils.CheckPathExist(filePath) {
		return os.Open(filePath)
	}
	return os.Create(filePath)
}

func GetLogText(teamUUID, importUUID string, startLine int) []string {
	filePath := logFilePath(teamUUID, importUUID)
	if !utils.CheckPathExist(filePath) {
		return []string{}
	}
	file, err := os.Open(filePath)
	if err != nil {
		return []string{fmt.Sprintf("log file open fail:%+v", err)}
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	lines := make([]string, 0)
	endLine := startLine + PerCount
	for fileScanner.Scan() {
		if lineCount >= startLine && lineCount < endLine {
			test := strings.TrimSpace(fileScanner.Text())
			if len(test) != 0 {
				lines = append(lines, fileScanner.Text())
			}
		}
		lineCount++
	}
	return lines
}

func GetLogTextByImportUUID(teamUUID, importUUID string) []string {
	filePath := logFilePath(teamUUID, importUUID)
	if !utils.CheckPathExist(filePath) {
		return []string{}
	}
	file, err := os.Open(filePath)
	if err != nil {
		return []string{fmt.Sprintf("log file open fail:%+v", err)}
	}
	fileScanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for fileScanner.Scan() {
		test := strings.TrimSpace(fileScanner.Text())
		if len(test) != 0 {
			lines = append(lines, fileScanner.Text())
		}
	}
	return lines
}

func DownloadAllLog(teamUUID string) ([]byte, error) {
	rootPath := logTeamPath(teamUUID)
	fileInfoList, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return nil, nil
	}
	sortFiles := make(utils.ByModTimeAsc, 0)
	for _, v := range fileInfoList {
		sortFiles = append(sortFiles, v)
	}
	sort.Sort(sortFiles)
	fileBytes := make([]byte, 0)
	for _, v := range sortFiles {
		fi, err := os.Open(fmt.Sprintf("%s/%s", rootPath, v.Name()))
		if err != nil {
			continue
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, fi); err != nil {
			return nil, err
		}
		fileBytes = append(fileBytes, buf.Bytes()...)
	}
	return fileBytes, nil
}

func DownloadCurrentLog(teamUUID, importUUID string) ([]byte, error) {
	rootPath := logTeamPath(teamUUID)
	fileInfoList, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return nil, nil
	}
	fiName := ""
	for _, v := range fileInfoList {
		if v.Name() == importUUID {
			fiName = v.Name()
			break
		}
	}
	if fiName == "" {
		return nil, nil
	}
	fi, err := os.Open(fmt.Sprintf("%s/%s", rootPath, fiName))
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, fi); err != nil {
		return nil, err
	}
	fileBytes := buf.Bytes()
	return fileBytes, nil
}

func logRootPath() string {
	filePath := fmt.Sprintf("%s/%s", common.GetCachePath(), Log)
	return filePath
}

func logFilePath(teamUUID, importUUID string) string {
	filePath := fmt.Sprintf("%s/%s", logTeamPath(teamUUID), importUUID)
	return filePath
}

func logTeamPath(teamUUID string) string {
	filePath := fmt.Sprintf("%s/%s/%s", common.GetCachePath(), Log, teamUUID)
	return filePath
}
