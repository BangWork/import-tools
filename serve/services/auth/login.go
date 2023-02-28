package auth

import (
	"fmt"

	"github.com/bangwork/import-tools/serve/models/ones"
	"github.com/bangwork/import-tools/serve/services"

	"github.com/bangwork/import-tools/serve/common"
	cookie2 "github.com/bangwork/import-tools/serve/services/cookie"
	"github.com/bangwork/import-tools/serve/utils"
)

func Login(req *services.LoginRequest) (string, error) {
	if len(req.URL) == 0 || len(req.Email) == 0 || len(req.Password) == 0 {
		return "", common.Errors(common.ParameterMissingError, req)
	}
	_, err := ones.LoginONES(req.URL, req.Email, req.Password)
	if err != nil {
		return "", err
	}
	cookie := GenCookie(req.URL, req.Email, req.Password)
	cookie2.ExpireMap.Put(cookie, req.Email)
	return cookie, nil
}

func Logout(cookie string) {
	cookie2.ExpireMap.Del(cookie)
}

func GenCookie(url, username, password string) string {
	str := fmt.Sprintf("%s|%s|%s", utils.Base64Encode(url), utils.Base64Encode(username), utils.Base64Encode(password))
	return utils.CBCEncrypt(str, common.GetEncryptKey())
}
