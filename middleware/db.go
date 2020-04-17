package middleware

import (
	"ginserver/db"
	"github.com/gin-gonic/gin"
)

// ConnectMgo ...
func ConnectMgo() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := db.Session()
		c.Set("db", session)
		c.Next()
	}
}
