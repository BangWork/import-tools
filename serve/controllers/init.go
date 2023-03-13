package controllers

import (
	"log"
	"net/http"

	"github.com/juju/errors"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/gin-gonic/gin"
)

func getONESHeader(c *gin.Context) (r map[string]string) {
	h, exist := c.Get("onesHeader")
	if !exist {
		log.Println("header not exist")
		return
	}
	r = h.(map[string]string)
	return
}

func getOrgUUID(c *gin.Context) (r string) {
	return c.GetString("orgUUID")
}

func getTeamUUID(c *gin.Context) (r string) {
	return c.Param("teamUUID")
}

func getCookie(c *gin.Context) (r string) {
	return c.GetString("cookie")
}

func getONESUrl(c *gin.Context) (r string) {
	return c.GetString("url")
}

func getUserUUID(c *gin.Context) (r string) {
	return c.GetString("userUUID")
}

func RenderJSON(c *gin.Context, err error, obj interface{}) {
	if err == nil {
		r := gin.H{
			"code": http.StatusOK,
		}
		if obj != nil {
			r["body"] = obj
		}
		c.JSON(http.StatusOK, r)
		c.Next()
		return
	}
	coErr, ok := err.(*common.Err)
	if !ok {
		log.Printf("%+v\n", errors.Trace(err))
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"err": err.Error(),
		})
	} else {
		r := gin.H{
			"code":     coErr.Code,
			"err_code": coErr.ErrCode,
		}
		if coErr.Body != nil {
			r["body"] = coErr.Body
		}
		log.Printf("ERROR: %+v", r)
		c.JSON(http.StatusOK, r)
	}
	c.Next()
}

func RenderJSONAndStop(c *gin.Context, err error, obj interface{}) {
	if err == nil {
		r := gin.H{
			"code": http.StatusOK,
		}
		if obj != nil {
			r["body"] = obj
		}
		c.JSON(http.StatusOK, r)
		c.Abort()
		return
	}
	res := err.(*common.Err)
	r := gin.H{
		"code":     res.Code,
		"err_code": res.ErrCode,
	}
	if res.Body != nil {
		r["body"] = res.Body
	}
	log.Printf("ERROR: %+v", r)
	c.JSON(http.StatusOK, r)
	c.Abort()
}
