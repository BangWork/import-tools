package ones

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/juju/errors"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
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

func GetJiraConfigInfo(orgUUID, url string, header map[string]string) (*JiraInfoStruct, error) {
	uri := fmt.Sprintf("%s", fmt.Sprintf(getOrgConfigUri, orgUUID, orgConfigTypeEnum))
	url = common.GenApiUrl(url, uri)
	resp, err := utils.GetWithHeader(url, header)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res := new(OrgConfigResponse)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	jiraInfo := new(JiraInfoStruct)
	err = json.Unmarshal([]byte(res.JiraInfo), jiraInfo)
	if err != nil {
		return nil, err
	}
	return jiraInfo, nil
}

func SetJiraConfigInfo(orgUUID, url string, header map[string]string, body interface{}) error {
	uri := fmt.Sprintf("%s", fmt.Sprintf(updateOrgConfigUri, orgUUID))
	url = common.GenApiUrl(url, uri)
	marshal, err := json.Marshal(body)
	if err != nil {
		return errors.Trace(err)
	}
	req := new(UpdateOrgConfigRequest)
	req.ConfigType = orgConfigTypeEnum
	req.ConfigData = OrgConfigResponse{JiraInfo: string(marshal)}
	resp, err := utils.PostJSONWithHeader(url, req, header)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Trace(common.Errors(common.ServerError, nil))
	}
	return nil
}
