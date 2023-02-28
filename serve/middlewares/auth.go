package middlewares

import (
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/controllers"
	"github.com/bangwork/import-tools/serve/services/cookie"
	"github.com/gin-gonic/gin"
)

func CheckLogin(c *gin.Context) {
	s, err := c.Cookie(common.LoginCookieName)
	if err != nil {
		controllers.RenderJSON(c, common.Errors(common.LoginCookieExpireError, nil), nil)
		return
	}
	ck := cookie.ExpireMap.Get(s)
	if len(ck) == 0 {
		controllers.RenderJSON(c, common.Errors(common.LoginCookieExpireError, nil), nil)
		return
	}
	c.Next()
}
