package ones

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

func CheckOrgPermission(url, orgUUID string, header map[string]string) (bool, error) {
	stamps, err := getOrgPermission(url, orgUUID, header)
	if err != nil {
		return false, err
	}
	for _, rule := range stamps.OrgPermissionRules.Rules {
		if rule.Permission == administerOrganization &&
			rule.UserDomainType == singleUser &&
			rule.UserDomainParam == header[common.UserID] {
			return true, nil
		}
	}
	return false, nil
}

func getOrgPermission(url, orgUUID string, header map[string]string) (*Stamps, error) {
	reqBody := map[string]int{
		"org_permission_rules": 0,
	}
	uri := fmt.Sprintf("%s?t=%s", fmt.Sprintf(stampsUri, orgUUID), dataTypeOrgPermissionRules)
	url = common.GenApiUrl(url, uri)
	resp, err := utils.PostJSONWithHeader(url, reqBody, header)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	stampsData := new(Stamps)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &stampsData); err != nil {
		return nil, err
	}
	return stampsData, nil
}
