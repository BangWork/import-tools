package cookie

import (
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils/expire_map"
)

var ExpireMap *expire_map.TTLMap

func InitCookieMap() error {
	ExpireMap = expire_map.New(common.GetCookieExpireTime())
	return nil
}
