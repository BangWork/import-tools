package ones

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	cookie2 "github.com/bangwork/import-tools/serve/services/cookie"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginONES(url, email, password string) (*http.Response, error) {
	if len(url) == 0 || len(email) == 0 || len(password) == 0 {
		log.Println("login param missing", url, email, password)
		return nil, common.Errors(common.ParameterMissingError, nil)
	}
	body := new(LoginRequest)
	body.Email = email
	body.Password = password
	url = common.GenApiUrl(url, loginUri)
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

func LoginONESAndSetAuth(cookie string) error {
	cookieValue := cookie2.ExpireMap.Get(cookie)
	if len(cookieValue) == 0 {
		log.Println("cookie expire")
		return common.Errors(common.LoginCookieExpireError, nil)
	}
	cookieInfo, err := DecryptCookieValue(cookieValue)
	if err != nil {
		log.Println("DecryptCookieValue err", err)
		return err
	}
	resp, err := LoginONES(cookieInfo.URL, cookieInfo.Email, cookieInfo.Password)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	cacheInfo := new(CookieCacheInfo)
	cacheInfo.ONESUserUUID = resp.Header.Get(common.UserID)
	cacheInfo.ONESAuthToken = resp.Header.Get(common.AuthToken)

	cookieValue, err = cacheInfo.GenCookieValue()
	if err != nil {
		log.Println("GenCookieValue err", err)
		return common.Errors(common.ServerError, err)
	}
	cookie2.ExpireMap.Put(cookie, cookieValue)
	return nil
}

type CookieCacheInfo struct {
	URL            string `json:"url"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ONESUserUUID   string `json:"ones_user_uuid"`
	ONESAuthToken  string `json:"ones_auth_token"`
	*LoginResponse `json:"login_response"`
	Language       string `json:"language"`
}

func (c *CookieCacheInfo) GenCookie() string {
	key := fmt.Sprintf("%s|%s|%s", c.URL, c.Email, c.Password)
	return utils.CBCEncrypt(key, common.GetEncryptKey())
}

func (c *CookieCacheInfo) GenAuthHeader() map[string]string {
	h := map[string]string{
		common.AuthToken: c.ONESAuthToken,
		common.UserID:    c.ONESUserUUID,
	}
	return h
}

func (c *CookieCacheInfo) GenCookieValue() (string, error) {
	j, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func DecryptCookieValue(v string) (*CookieCacheInfo, error) {
	cookieCacheInfo := new(CookieCacheInfo)
	err := json.Unmarshal([]byte(v), &cookieCacheInfo)
	if err != nil {
		log.Println("json unmarshal err", err)
		return nil, err
	}
	return cookieCacheInfo, nil
}

func DecryptCookieValueByCookie(c string) (*CookieCacheInfo, error) {
	c = cookie2.ExpireMap.Get(c)
	value, err := DecryptCookieValue(c)
	if err != nil {
		return nil, err
	}
	return value, nil
}
