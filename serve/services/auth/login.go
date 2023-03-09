package auth

import (
	"encoding/json"
	"io/ioutil"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/models/ones"
	"github.com/bangwork/import-tools/serve/services"
	cookie2 "github.com/bangwork/import-tools/serve/services/cookie"
)

func Login(req *services.LoginRequest) (string, error) {
	if len(req.URL) == 0 || len(req.Email) == 0 || len(req.Password) == 0 {
		return "", common.Errors(common.ParameterMissingError, req)
	}
	resp, err := ones.LoginONES(req.URL, req.Email, req.Password)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	cacheInfo := new(ones.CookieCacheInfo)
	cacheInfo.URL = req.URL
	cacheInfo.Email = req.Email
	cacheInfo.Password = req.Password
	cacheInfo.ONESUserUUID = resp.Header.Get(common.UserID)
	cacheInfo.ONESAuthToken = resp.Header.Get(common.AuthToken)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBody := new(ones.LoginResponse)
	if err = json.Unmarshal(data, &respBody); err != nil {
		return "", err
	}

	header := map[string]string{
		common.UserID:    cacheInfo.ONESUserUUID,
		common.AuthToken: cacheInfo.ONESAuthToken,
	}
	if e := checkLoginPermission(respBody, req.URL, header); e != nil {
		return "", e
	}

	if e := checkONESVersion(req.URL, header); e != nil {
		return "", e
	}

	cacheInfo.LoginResponse = respBody
	cookie := cacheInfo.GenCookie()
	cookieValue, err := cacheInfo.GenCookieValue()
	if err != nil {
		return "", common.Errors(common.ServerError, err)
	}
	//cookie = "aaaaaaaaa"
	cookie2.ExpireMap.Put(cookie, cookieValue)
	return cookie, nil
}

func checkLoginPermission(r *ones.LoginResponse, url string, header map[string]string) error {
	if r.User.UUID == r.Org.Owner {
		return nil
	}
	if r.Org.MultiTeam {
		havePermission, err := ones.CheckOrgPermission(url, r.Org.UUID, header)
		if err != nil {
			return err
		}
		if !havePermission {
			return common.Errors(common.NotOrganizationAdministratorError, nil)
		}
		return nil
	}

	havePermission, err := ones.CheckTeamPermission(url, r.Teams[0].UUID, header)
	if err != nil {
		return err
	}
	if !havePermission {
		return common.Errors(common.NotSuperAdministratorError, nil)
	}
	return nil
}

func checkONESVersion(url string, header map[string]string) error {
	_, err := ones.CheckONESVersion(url, header)
	return err
}

func Logout(cookie string) {
	cookie2.ExpireMap.Del(cookie)
}
