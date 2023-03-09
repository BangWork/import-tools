package middlewares

import (
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/controllers"
	"github.com/bangwork/import-tools/serve/models/ones"
	"github.com/bangwork/import-tools/serve/services/cookie"
	"github.com/gin-gonic/gin"
)

func CheckLogin(c *gin.Context) {
	s, err := c.Cookie(common.LoginCookieName)
	//s = "aaaaaaaaa"
	if err != nil {
		controllers.RenderJSONAndStop(c, common.Errors(common.LoginCookieExpireError, err), nil)
		return
	}
	cookieValue := cookie.ExpireMap.Get(s)
	if len(cookieValue) == 0 {
		controllers.RenderJSONAndStop(c, common.Errors(common.LoginCookieExpireError, nil), nil)
		return
	}
	cookieInfo, err := ones.DecryptCookieValue(cookieValue)
	if err != nil {
		controllers.RenderJSONAndStop(c, common.Errors(common.ServerError, err), nil)
		return
	}
	c.Set("cookie", s)
	c.Set("userUUID", cookieInfo.ONESUserUUID)
	c.Set("orgUUID", cookieInfo.LoginResponse.Org.UUID)
	c.Set("url", cookieInfo.URL)
	c.Set("onesHeader", map[string]string{
		common.AuthToken: cookieInfo.ONESAuthToken,
		common.UserID:    cookieInfo.ONESUserUUID,
	})
	c.Next()
}
