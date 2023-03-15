package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"

	"github.com/bangwork/import-tools/serve"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services/cache"

	"github.com/bangwork/import-tools/serve/router"
)

var (
	port             int
	cachePath        string
	maxScanTokenSize int
)

func initFlags() {
	flag.IntVar(&port, "port", -1, "http server port")
	flag.StringVar(&cachePath, "cache-path", common.Path, "cache path")
	flag.StringVar(&cache.SharedDiskPath, "shared-disk-path", "", "the path of the shared disk")
	flag.IntVar(&maxScanTokenSize, "max-scan-token-size", bufio.MaxScanTokenSize*1000, "max scan token size")
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

	common.SetCachePath(cachePath)
	common.SetMaxScanTokenSize(maxScanTokenSize)
}

func main() {
	initFlags()
	serve.Init()
	router.Run(port)
}
