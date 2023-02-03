package main

import (
	"flag"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/router"
	"github.com/bangwork/import-tools/serve/services/cache"

	_ "github.com/bangwork/import-tools/serve"
)

var port int

func initFlags() {
	flag.IntVar(&port, "port", common.DefaultHTTPPort, "http server port")
	flag.StringVar(&cache.SharedDiskPath, "shared-disk-path", "", "the path of the shared disk")
	flag.Parse()
}

func main() {
	initFlags()
	router.Run(port)
}
