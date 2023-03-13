package ones

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
	"github.com/juju/errors"
)

func GetUnBoundIssueTypes(url, teamUUID string, header map[string]string) ([]*UnBoundIssueTypes, error) {
	url = common.GenApiUrl(url, fmt.Sprintf(unboundIssueTypesUri, teamUUID))
	resp, err := utils.GetWithHeader(url, header)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	res := make([]*UnBoundIssueTypes, 0)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}

func GetThirdIssueTypeBind(url, teamUUID string, h map[string]string) ([]*BindIssueType, error) {
	uri := fmt.Sprintf(thirdIssueTypeBindUri, teamUUID)
	url = common.GenApiUrl(url, uri)
	resp, err := utils.GetWithHeader(url, h)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	res := make([]*BindIssueType, 0)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}

func MapThirdIssueTypeBind(url, teamUUID string, h map[string]string) (map[string]*BindIssueType, error) {
	bind, err := GetThirdIssueTypeBind(url, teamUUID, h)
	if err != nil {
		return nil, errors.Trace(err)
	}
	mapRes := make(map[string]*BindIssueType)
	for _, v := range bind {
		mapRes[v.IssueTypeID] = v
	}
	return mapRes, nil
}
