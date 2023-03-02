package common

import "github.com/bangwork/import-tools/serve/common"

var (
	TeamWorkdays = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
)

type Config struct {
	Area string `json:"area"`
}

func GetConfig() *Config {
	c := new(Config)
	c.Area = common.GetInstallArea()
	return c
}
