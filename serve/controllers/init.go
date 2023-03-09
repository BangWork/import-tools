package controllers

import (
	"log"
	"net/http"

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
