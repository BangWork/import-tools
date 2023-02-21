package main

import (
	"bufio"
	"flag"

	"github.com/bangwork/import-tools/serve"
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/router"
)

var port = flag.Int("p", common.DefaultHTTPPort, "http server port")
var cachePath = flag.String("c", common.Path, "cache file path")
var maxScanTokenSize = flag.Int("max_scan_token_size", bufio.MaxScanTokenSize*1000, "max scan token size")

func main() {
	flag.Parse()
	common.SetCachePath(*cachePath)
	common.SetMaxScanTokenSize(*maxScanTokenSize)
	serve.Init()
	router.Run(*port)
}
