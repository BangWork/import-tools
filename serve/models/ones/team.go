package ones

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

func CheckTeamPermission(url, teamUUID string, header map[string]string) (bool, error) {
	rules, err := getTeamPermission(url, teamUUID, header)
	if err != nil {
		return false, err
	}
	for _, rule := range rules.Rules {
		if rule.Permission == superAdministrator &&
			rule.UserDomainType == singleUser &&
			rule.UserDomainParam == header[common.UserID] {
			return true, nil
		}
	}
	return false, nil
}

func getTeamPermission(url, teamUUID string, header map[string]string) (*PermissionRules, error) {
	reqBody := map[string]string{
		"context_type": "team",
	}
	url = common.GenApiUrl(url, fmt.Sprintf(listPermissionRulesUri, teamUUID))
	resp, err := utils.PostJSONWithHeader(url, reqBody, header)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	stampsData := new(PermissionRules)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &stampsData); err != nil {
		return nil, err
	}
	return stampsData, nil
}
