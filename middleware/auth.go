package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(302, "/auth/login")
			c.Abort()
			return
		}
		// Set isAuthenticated for templates
		c.Set("isAuthenticated", true)
		c.Next()
	}
}

// CheckAuth Middleware to check authentication without redirection
func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Set("isAuthenticated", false)
		} else {
			c.Set("isAuthenticated", true)
		}
		c.Next()
	}
}
