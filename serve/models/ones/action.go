package ones

import (
	"fmt"
	"net/http"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

func InterruptImport(teamUUID, url, importBatchTaskUUID string, header map[string]string) error {
	if importBatchTaskUUID == "" {
		return nil
	}
	url = common.GenApiUrl(url, fmt.Sprintf(interruptImport, teamUUID, importBatchTaskUUID))
	resp, err := utils.PostJSONWithHeader(url, []string{}, header)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return common.Errors(common.ServerError, nil)
	}
	return nil
}
