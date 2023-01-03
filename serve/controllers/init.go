package controllers

import (
	"log"
	"net/http"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/gin-gonic/gin"
)

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
