package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"

	"github.com/bangwork/import-tools/serve"
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/router"
	"github.com/bangwork/import-tools/serve/services/cache"
)

func main() {
	initFlags()
	common.SetCachePath(cachePath)
	common.SetMaxScanTokenSize(maxScanTokenSize)
	common.SetEncryptKey(encryptKey)
	serve.Init()
	router.Run(port)
}

var port int
var cachePath string
var encryptKey string
var maxScanTokenSize int

func initFlags() {
	flag.IntVar(&port, "port", -1, "http server port")
	flag.StringVar(&cache.SharedDiskPath, "shared-disk-path", "", "the path of the shared disk")
	flag.StringVar(&cachePath, "cache-path", common.Path, "cache file path")
	flag.StringVar(&encryptKey, "aes-key", common.DefaultAesKey, "cookie encrypt aes key, 16 characters")
	flag.IntVar(&maxScanTokenSize, "max_scan_token_size", bufio.MaxScanTokenSize*1000, "max scan token size")
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

	if cache.SharedDiskPath == "" {
		cache.SharedDiskPath = os.Getenv("IMPORT_TOOLS_SHARED_DISK_PATH")
	}
}
