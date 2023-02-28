package controllers

import (
	"log"
	"net/http"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
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
