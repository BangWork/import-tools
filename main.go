package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"

	"github.com/bangwork/import-tools/serve"
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/router"
)

var port int
var cachePath string
var encryptKey string
var maxScanTokenSize int
var installArea string
var jiraLocalHome string

func initFlags() {
	flag.IntVar(&port, "port", -1, "http server port")
	flag.StringVar(&common.SharedDiskPath, "shared-disk-path", "", "the path of the shared disk")
	flag.StringVar(&cachePath, "cache-path", common.Path, "cache file path")
	flag.StringVar(&encryptKey, "aes-key", common.DefaultAesKey, "cookie encrypt aes key, 16 characters")
	flag.IntVar(&maxScanTokenSize, "max-scan-token-size", bufio.MaxScanTokenSize*1000, "max scan token size")
	flag.StringVar(&installArea, "area", common.InstallAreaAsia, "install area")
	flag.StringVar(&jiraLocalHome, "jira-local-home", "", "jira local home")
	flag.Parse()

	if port == -1 {
		p := os.Getenv("IMPORT_TOOLS_PORT")
		if p != "" {
			var err error
			port, err = strconv.Atoi(p)
			if err != nil {
				port = common.DefaultHTTPPort
			}
		}
	}

	if common.SharedDiskPath == "" {
		common.SharedDiskPath = os.Getenv("IMPORT_TOOLS_SHARED_DISK_PATH")
	}
}

func main() {
	initFlags()
	common.SetCachePath(cachePath)
	common.SetMaxScanTokenSize(maxScanTokenSize)
	common.SetEncryptKey(encryptKey)
	common.SetInstallArea(installArea)
	common.SetJiraLocalHome(jiraLocalHome)
	serve.Init()
	router.Run(port)
}
