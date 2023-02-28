package serve

import (
	"fmt"

	"github.com/bangwork/import-tools/serve/services/cookie"

	"github.com/bangwork/import-tools/serve/services/cache"
	"github.com/bangwork/import-tools/serve/services/importer/sync"
	"github.com/bangwork/import-tools/serve/services/log"
)

func Init() {
	onStart(cache.InitCacheFile)
	onStart(log.InitLogDir)
	onStart(sync.InitResolverFactory)
	onStart(cookie.InitCookieMap)
}

func onStart(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Sprintf("Error at onStart: %s\n", err))
	}
}
