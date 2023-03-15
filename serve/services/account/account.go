package account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bangwork/import-tools/serve/services"
	"github.com/juju/errors"

	"github.com/bangwork/import-tools/serve/services/issue_type"

	"github.com/bangwork/import-tools/serve/services/cache"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

const (
	// team
	interruptImport        = "/team/%s/queue/%s/interrupt"
	confirmImportUri       = "team/%s/jira/import"
	setPasswordUri         = "team/%s/jira/user/set_password"
	importStatusUri        = "team/%s/queues/list"
	importProgressUri      = "team/%s/queue/%s/progress"
	sendImportDataUri      = "team/%s/importer/import/%s"
	issueTypeListUri       = "team/%s/issue_types"
	thirdIssueTypeBindUri  = "team/%s/third_issue_type_bind"
	listPermissionRulesUri = "team/%s/lite_context_permission_rules"

	// organization
	importLogUri     = "organization/%s/importer/log/%s/%d"
	importHistoryUri = "organization/%s/importer/history"
	fileConfigUri    = "organization/%s/file_config"
	stampsUri        = "organization/%s/stamps/data"

	// login
	loginUri = "auth/v2/login"

	dataTypeOrgPermissionRules = "org_permission_rules"
	administerOrganization     = "administer_organization"
	superAdministrator         = "super_administrator"
	singleUser                 = "single_user"
)

type Account struct {
	URL         string `json:"url"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	LocalHome   string `json:"local_home"`
	BackupName  string `json:"backup_name"`
	FileStorage string `json:"file_storage"`

	AuthHeader map[string]string
	OrgInfo    Organization
	TeamInfo   Team

	Cache *cache.Cache
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	User  User         `json:"user"`
	Org   Organization `json:"org"`
	Teams []Team       `json:"teams"`
}

type confirmImportResponse struct {
	UUID string `json:"uuid"`
}

type loginErrorResponse struct {
	//ErrCode    string `json:"errcode"`
	RetryCount int `json:"retry_count"`
}

type User struct {
	UUID string `json:"uuid"`
}

type Team struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type Organization struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	MultiTeam bool   `json:"visibility"`
}

type Stamps struct {
	OrgPermissionRules PermissionRules `json:"org_permission_rules"`
}

type FileConfig struct {
	FileStorage string `json:"file_storage"`
}

type PermissionRule struct {
	Permission      string `json:"permission"`
	UserDomainType  string `json:"user_domain_type"`
	UserDomainParam string `json:"user_domain_param"`
}

type PermissionRules struct {
	Rules []PermissionRule `json:"permission_rules"`
}

type ImportHistory struct {
	TeamUUID   string `json:"team_uuid"`
	TeamName   string `json:"team_name"`
	ImportList []struct {
		ImportTime   int    `json:"import_time"`
		JiraVersion  string `json:"jira_version"`
		JiraServerID string `json:"jira_server_id"`
	} `json:"import_list"`
}

type ResolveResultResponse struct {
	*cache.ResolveResult `json:"resolve_result"`
	ImportHistory        []*ImportHistory `json:"import_history"`
}

func (r *Account) SetCurrentCache(c *cache.Cache) {
	r.Cache = c
}

func (r *Account) GetImportHistory() ([]*ImportHistory, error) {
	resp, err := r.postLogin()
	if err != nil {
		return nil, common.Errors(common.NetworkError, nil)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, common.Errors(common.AccountError, nil)
	}
	r.AuthHeader = make(map[string]string)
	r.AuthHeader[common.AuthToken] = resp.Header.Get(common.AuthToken)
	r.AuthHeader[common.UserID] = resp.Header.Get(common.UserID)

	history, err := r.getImportHistory()
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (r *Account) Login() error {
	resp, err := r.postLogin()
	if err != nil {
		return common.Errors(common.NetworkError, nil)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return common.Errors(common.AccountError, nil)
	}
	r.AuthHeader = make(map[string]string)
	r.AuthHeader[common.AuthToken] = resp.Header.Get(common.AuthToken)
	r.AuthHeader[common.UserID] = resp.Header.Get(common.UserID)

	info, err := cache.GetCacheInfo(r.Key())
	if err != nil {
		return err
	}
	r.SetCurrentCache(info)
	return nil
}

func (r *Account) SendImportData(resourceTypeString string, data string) error {
	byteData := []byte(data)
	url := common.GenApiUrl(r.URL, fmt.Sprintf(sendImportDataUri, r.Cache.ImportTeamUUID, resourceTypeString))
	resp, err := utils.PostByteWithHeader(url, byteData, r.AuthHeader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func (r *Account) GetImportLog(count int) ([]string, error) {
	url := common.GenApiUrl(r.URL, fmt.Sprintf(importLogUri, r.Cache.OrgUUID, r.Cache.ImportUUID, count))
	resp, err := utils.GetWithHeader(url, r.AuthHeader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	res := make([]string, 0)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Account) GetImportStatus(resourceUUID string) (status string, err error) {
	url := common.GenApiUrl(r.URL, fmt.Sprintf(importStatusUri, r.Cache.ImportTeamUUID))
	resp, err := utils.GetWithHeader(url, r.AuthHeader)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	res := new(services.QueueListResponse)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	type extraStruct struct {
		ResourceUUID string `json:"resource_uuid"`
	}
	for _, v := range res.BatchTasks {
		e := new(extraStruct)
		if err = json.Unmarshal([]byte(v.Extra), e); err != nil {
			log.Println("json err", err)
			continue
		}
		if e.ResourceUUID == resourceUUID {
			status = v.JobStatus
			break
		}
	}
	return status, nil
}

func (r *Account) GetImportStatusByUUID(batchTaskUUID string) (status string, err error) {
	url := common.GenApiUrl(r.URL, fmt.Sprintf(importProgressUri, r.Cache.ImportTeamUUID, batchTaskUUID))
	resp, err := utils.GetWithHeader(url, r.AuthHeader)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	res := new(services.QueueStruct)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	return res.JobStatus, nil
}

func (r *Account) SetPassword(resourceUUID, password string) error {
	url := common.GenApiUrl(r.URL, fmt.Sprintf(setPasswordUri, r.Cache.ImportTeamUUID))
	reqBody := map[string]string{
		"resource_uuid":      resourceUUID,
		"password":           password,
		"confirmed_password": password,
	}
	resp, err := utils.PostJSONWithHeader(url, reqBody, r.AuthHeader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return common.Errors(common.ServerError, nil)
	}
	return nil
}

func (r *Account) InterruptImport() error {
	if services.ImportBatchTaskUUID == "" {
		return nil
	}
	url := common.GenApiUrl(r.URL, fmt.Sprintf(interruptImport, r.Cache.ImportTeamUUID, services.ImportBatchTaskUUID))
	resp, err := utils.PostJSONWithHeader(url, []string{}, r.AuthHeader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return common.Errors(common.ServerError, nil)
	}
	return nil
}

func (r *Account) ConfirmImport(reqBody *services.ConfirmImportRequest) (string, error) {
	url := common.GenApiUrl(r.URL, fmt.Sprintf(confirmImportUri, r.Cache.ImportTeamUUID))
	resp, err := utils.PostJSONWithHeader(url, reqBody, r.AuthHeader)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", common.Errors(common.ServerError, nil)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBody := new(confirmImportResponse)
	if err = json.Unmarshal(data, &respBody); err != nil {
		return "", err
	}

	return respBody.UUID, nil
}

func (r *Account) GetIssueTypeList() (*services.IssueTypeListResponse, error) {
	resp, err := r.postLogin()
	if err != nil {
		return nil, common.Errors(common.NetworkError, nil)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, common.Errors(common.AccountError, nil)
	}
	r.AuthHeader = make(map[string]string)
	r.AuthHeader[common.AuthToken] = resp.Header.Get(common.AuthToken)
	r.AuthHeader[common.UserID] = resp.Header.Get(common.UserID)

	issueTypes, err := r.getIssueTypeList()
	if err != nil {
		return nil, errors.Trace(err)
	}

	thirdIssueTypesBind, err := r.getThirdIssueTypeBind()
	if err != nil {
		return nil, errors.Trace(err)
	}

	res := new(services.IssueTypeListResponse)
	res.JiraList = thirdIssueTypesBind
	res.ONESList = issueTypes

	return res, nil
}

func (r *Account) CheckONESAccount() error {
	if len(r.URL) == 0 || len(r.Email) == 0 || len(r.Password) == 0 || len(r.LocalHome) == 0 || len(r.BackupName) == 0 {
		return common.Errors(common.ParameterMissingError, nil)
	}
	resp, err := r.postLogin()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respBody := new(loginResponse)
	if err = json.Unmarshal(data, &respBody); err != nil {
		return err
	}

	r.AuthHeader = make(map[string]string)
	r.AuthHeader[common.AuthToken] = resp.Header.Get(common.AuthToken)
	r.AuthHeader[common.UserID] = resp.Header.Get(common.UserID)
	r.OrgInfo = respBody.Org
	r.TeamInfo = respBody.Teams[0]
	config, err := r.postOrgConfig()
	if err != nil {
		return err
	}
	r.FileStorage = config.FileStorage
	if respBody.User.UUID == respBody.Org.Owner {
		return nil
	}

	if respBody.Org.MultiTeam {
		havePermission, err := r.checkOrgPermission()
		if err != nil {
			return err
		}
		if !havePermission {
			return common.Errors(common.NotOrganizationAdministratorError, nil)
		}
		return nil
	}

	havePermission, err := r.checkTeamPermission()
	if err != nil {
		return err
	}
	if !havePermission {
		return common.Errors(common.NotSuperAdministratorError, nil)
	}
	return nil
}

func (r *Account) Key() string {
	return cache.GenCacheKey(r.URL)
}

func (r *Account) SetCache() error {
	info, err := cache.GetCacheInfo(r.Key())
	if err != nil {
		return err
	}
	fileSize, err := utils.GetFileSize(common.GenBackupFilePath(r.LocalHome, r.BackupName))
	if err != nil {
		return common.Errors(common.NotFoundError, nil)
	}
	info.ExpectedResolveTime = getExpectedResolveTime(fileSize)
	info.URL = r.URL
	info.Email = r.Email
	info.Password = r.Password
	info.MultiTeam = r.OrgInfo.MultiTeam
	info.ResolveStartTime = time.Now().Unix()
	info.OrgName = r.OrgInfo.Name
	info.OrgUUID = r.OrgInfo.UUID
	info.TeamUUID = r.TeamInfo.UUID
	info.TeamName = r.TeamInfo.Name
	info.LocalHome = r.LocalHome
	info.BackupName = r.BackupName
	info.ResolveStatus = common.ResolveStatusInProgress
	info.ImportUserUUID = r.AuthHeader[common.UserID]
	info.FileStorage = r.FileStorage
	if info.ImportResult != nil {
		info.ImportResult = nil
	}

	return cache.SetCacheInfo(r.Key(), info)
}

func getExpectedResolveTime(fileSize int64) int64 {
	fileSizeM := fileSize / 1024 / 1024
	for size, t := range services.MapResolveTime {
		if fileSizeM > size {
			return t
		}
	}
	return 100
}

func (r *Account) checkOrgPermission() (bool, error) {
	stamps, err := r.postOrgPermission()
	if err != nil {
		return false, err
	}
	for _, rule := range stamps.OrgPermissionRules.Rules {
		if rule.Permission == administerOrganization &&
			rule.UserDomainType == singleUser &&
			rule.UserDomainParam == r.AuthHeader["ones-user-id"] {
			return true, nil
		}
	}
	return false, nil
}

func (r *Account) checkTeamPermission() (bool, error) {
	rules, err := r.postTeamPermission()
	if err != nil {
		return false, err
	}
	for _, rule := range rules.Rules {
		if rule.Permission == superAdministrator &&
			rule.UserDomainType == singleUser &&
			rule.UserDomainParam == r.AuthHeader["ones-user-id"] {
			return true, nil
		}
	}
	return false, nil
}

func (r *Account) postLogin() (*http.Response, error) {
	if len(r.Email) == 0 || len(r.Password) == 0 || len(r.URL) == 0 {
		info, err := cache.GetCacheInfo(r.Key())
		if err != nil {
			return nil, err
		}
		r.URL = info.URL
		r.Email = info.Email
		r.Password = info.Password
	}
	body := new(LoginRequest)
	body.Email = r.Email
	body.Password = r.Password
	url := common.GenApiUrl(r.URL, loginUri)
	resp, err := utils.PostJSON(url, body)
	if err != nil {
		return nil, common.Errors(common.NetworkError, nil)
	}
	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		respBody := new(loginErrorResponse)
		if err = json.Unmarshal(data, &respBody); err != nil {
			return nil, err
		}
		return nil, common.Errors(common.AccountError, respBody)
	}
	return resp, nil
}

func (r *Account) getThirdIssueTypeBind() ([]*services.JiraIssueType, error) {
	uri := fmt.Sprintf(thirdIssueTypeBindUri, r.Cache.ImportTeamUUID)
	url := common.GenApiUrl(r.URL, uri)
	resp, err := utils.GetWithHeader(url, r.AuthHeader)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	res := make([]*services.JiraIssueType, 0)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}

func (r *Account) getIssueTypeList() ([]*services.ONESIssueType, error) {
	uri := fmt.Sprintf(issueTypeListUri, r.Cache.ImportTeamUUID)
	url := common.GenApiUrl(r.URL, uri)
	resp, err := utils.GetWithHeader(url, r.AuthHeader)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	type issueTypeResp struct {
		IssueTypes []*services.ONESIssueType `json:"issue_types"`
	}
	res := new(issueTypeResp)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, errors.Trace(err)
	}
	filterRes := make([]*services.ONESIssueType, 0)
	for _, v := range res.IssueTypes {
		if issue_type.SupportDetailType[v.DetailType] {
			filterRes = append(filterRes, v)
		}
	}
	return filterRes, nil
}

func (r *Account) getImportHistory() ([]*ImportHistory, error) {
	uri := fmt.Sprintf("%s", fmt.Sprintf(importHistoryUri, r.Cache.OrgUUID))
	url := common.GenApiUrl(r.URL, uri)
	resp, err := utils.GetWithHeader(url, r.AuthHeader)
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

func (r *Account) postOrgPermission() (*Stamps, error) {
	reqBody := map[string]int{
		"org_permission_rules": 0,
	}
	uri := fmt.Sprintf("%s?t=%s", fmt.Sprintf(stampsUri, r.OrgInfo.UUID), dataTypeOrgPermissionRules)
	url := common.GenApiUrl(r.URL, uri)
	resp, err := utils.PostJSONWithHeader(url, reqBody, r.AuthHeader)
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

func (r *Account) postOrgConfig() (*FileConfig, error) {
	uri := fmt.Sprintf(fileConfigUri, r.OrgInfo.UUID)
	url := common.GenApiUrl(r.URL, uri)
	resp, err := utils.GetWithHeader(url, r.AuthHeader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData := new(FileConfig)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &respData); err != nil {
		return nil, err
	}
	return respData, nil
}

func (r *Account) postTeamPermission() (*PermissionRules, error) {
	reqBody := map[string]string{
		"context_type": "team",
	}
	url := common.GenApiUrl(r.URL, fmt.Sprintf(listPermissionRulesUri, r.TeamInfo.UUID))
	resp, err := utils.PostJSONWithHeader(url, reqBody, r.AuthHeader)
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
