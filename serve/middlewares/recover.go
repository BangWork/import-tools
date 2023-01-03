package middlewares

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	write := io.MultiWriter(os.Stderr, os.Stdout, gin.DefaultWriter)
	return gin.RecoveryWithWriter(write)
}
