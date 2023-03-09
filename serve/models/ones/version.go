package ones

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

func CheckONESVersion(url string, header map[string]string) (bool, error) {
	uri := fmt.Sprintf(checkVersionUri, common.CurrentVersion)
	url = common.GenApiUrl(url, uri)
	resp, err := utils.GetWithHeader(url, header)
	if err != nil {
		log.Println("check ones version err", err)
		return false, err
	}
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, common.Errors(common.ONESVersionError, nil)
}
