package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")                                // 允许所有来源
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // 允许的请求方法
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")     // 允许的请求头

		// 如果是OPTIONS请求，提前返回
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
