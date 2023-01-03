package middlewares

import "github.com/gin-gonic/gin"

func Logger() gin.HandlerFunc {
	conf := gin.LoggerConfig{}
	return gin.LoggerWithConfig(conf)
}
