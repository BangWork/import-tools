package main

import (
	"flag"

	"github.com/bangwork/import-tools/serve/common"

	_ "github.com/bangwork/import-tools/serve"
	"github.com/bangwork/import-tools/serve/router"
)

var port = flag.Int("p", common.DefaultHTTPPort, "http server port")

func main() {
	flag.Parse()
	router.Run(*port)
}
