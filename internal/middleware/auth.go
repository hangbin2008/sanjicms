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
			"/":         true, // 首页需要检查登录状态
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
			// 如果是首页，检查用户是否已登录
			if path == "/" {
				// 从Cookie获取token，这里简化处理，实际应该使用JWT验证
				token, err := c.Cookie("token")
				if err != nil || token == "" {
					// 未登录，重定向到登录页面
					c.Redirect(http.StatusFound, "/login")
					c.Abort()
					return
				}
			}
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

		c.Next()
	}
}
