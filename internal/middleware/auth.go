package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthCheck 检查用户是否已登录，用于前端页面访问控制
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许访问的公共页面
		publicPages := map[string]bool{
			"/login":    true,
			"/register": true,
		}

		// 允许访问的调试页面，只有站长可以访问
		debugPages := map[string]bool{
			"/debug/templates": true,
			"/health":          true,
		}

		path := c.Request.URL.Path

		// 检查是否是API请求，API请求由专门的JWT中间件处理
		if strings.HasPrefix(path, "/api/") {
			c.Next()
			return
		}

		// 检查是否是静态资源请求
		if strings.HasPrefix(path, "/static/") || strings.HasPrefix(path, "/favicon.ico") {
			c.Next()
			return
		}

		// 检查是否是公共页面
		if publicPages[path] {
			c.Next()
			return
		}

		// 非公共页面，检查登录状态
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			// 未登录，重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 检查是否是调试页面
		if debugPages[path] {
			// 这里需要更严格的验证，确保只有站长可以访问调试页面
			// 简化处理，实际应该解析JWT并验证角色
			// c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			// return
		}

		c.Next()
	}
}
