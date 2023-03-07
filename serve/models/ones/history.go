package ones

import (
	"encoding/json"
	"fmt"
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
	"io/ioutil"
)

func GetImportHistory(orgUUID, url string, header map[string]string) ([]*ImportHistory, error) {
	uri := fmt.Sprintf("%s", fmt.Sprintf(importHistoryUri, orgUUID))
	url = common.GenApiUrl(url, uri)
	resp, err := utils.GetWithHeader(url, header)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res := make([]*ImportHistory, 0)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}
