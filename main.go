package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/bangwork/import-tools/serve"
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/router"
	"github.com/bangwork/import-tools/serve/services/cache"
)

var port int
var cachePath string
var encryptKey string

func initFlags() {
	flag.IntVar(&port, "port", -1, "http server port")
	flag.StringVar(&cache.SharedDiskPath, "shared-disk-path", "", "the path of the shared disk")
	flag.StringVar(&cachePath, "cache-path", common.Path, "cache file path")
	flag.StringVar(&encryptKey, "aes-key", common.DefaultAesKey, "cookie encrypt aes key, 16 characters")
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

func main() {
	initFlags()
	common.SetCachePath(cachePath)
	common.SetEncryptKey(encryptKey)
	serve.Init()
	router.Run(port)
}
