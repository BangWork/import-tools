package main

import (
	"flag"

	"github.com/bangwork/import-tools/serve"
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/router"
)

var port = flag.Int("p", common.DefaultHTTPPort, "http server port")
var cachePath = flag.String("c", common.Path, "cache file path")

func main() {
	flag.Parse()
	common.SetCachePath(*cachePath)
	serve.Init()
	router.Run(*port)
}
