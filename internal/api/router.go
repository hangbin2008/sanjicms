package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hangbin2008/sanjicms/internal/middleware"
	"github.com/hangbin2008/sanjicms/internal/models"
	"github.com/hangbin2008/sanjicms/internal/service"
	"github.com/hangbin2008/sanjicms/pkg/config"
)

// SetupRouter 配置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	// 创建Gin引擎
	router := gin.Default()

	// 设置静态文件服务
	router.Static("/static", "./static")

	// 设置模板引擎
	router.LoadHTMLGlob("./templates/*")

	// 调试：检查模板文件是否存在
	// 注意：这行代码会在生产环境中产生日志，可以根据需要保留或删除
	router.GET("/debug/templates", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "模板文件加载配置",
			"pattern": "./templates/*",
		})
	})

	// 添加认证检查中间件，用于前端页面访问控制
	router.Use(middleware.AuthCheck())

	// 创建JWT中间件
	jwtConfig := middleware.NewJWTConfig(cfg)

	// 创建服务实例
	userService := service.NewUserService(cfg, jwtConfig)
	questionService := service.NewQuestionService()
	examService := service.NewExamService(questionService)
	captchaService := service.NewCaptchaService()

	// 创建控制器实例
	controllers := NewControllers(userService, questionService, examService, captchaService)

	// 健康检查路由 - 只有站长可以访问
	router.GET("/health", middleware.RoleAuth("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "基层三基考试系统 API 运行正常",
		})
	})

	// 公共路由组
	public := router.Group("/api")
	{
		// 验证码路由
		public.GET("/captcha", controllers.GenerateCaptcha)
		// 用户认证路由
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
	}

	// 受保护路由组
	protected := router.Group("/api")
	protected.Use(middleware.JWTAuth(jwtConfig))
	{
		// 用户相关路由
		protected.GET("/user/me", controllers.GetCurrentUser)
		protected.PUT("/user/me", controllers.UpdateUser)

		// 调试路由组 - 只有站长可以访问
		debug := protected.Group("/debug")
		debug.Use(middleware.RoleAuth("admin"))
		{
			debug.GET("/templates", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "模板文件加载配置",
					"pattern": "./templates/*",
				})
			})
		}

		// 题库相关路由（需要管理员权限）
		bank := protected.Group("/banks")
		bank.Use(middleware.RoleAuth("admin", "manager"))
		{
			bank.POST("/", controllers.CreateQuestionBank)
			bank.GET("/", controllers.ListQuestionBanks)
			bank.GET("/:id", controllers.GetQuestionBankByID)
		}

		// 题目相关路由（需要管理员权限）
		question := protected.Group("/questions")
		question.Use(middleware.RoleAuth("admin", "manager"))
		{
			question.POST("/", controllers.CreateQuestion)
			question.GET("/bank/:bank_id", controllers.ListQuestionsByBank)
			question.GET("/:id", controllers.GetQuestionByID)
		}

		// 试卷相关路由
		exam := protected.Group("/exams")
		{
			// 生成试卷（需要管理员权限）
			exam.POST("/generate", middleware.RoleAuth("admin", "manager"), controllers.GenerateExam)
			// 获取试卷列表
			exam.GET("/", controllers.ListExams)
			// 获取试卷详情
			exam.GET("/:id", controllers.GetExamByID)
			// 开始考试
			exam.POST("/:id/start", controllers.StartExam)
			// 提交试卷
			exam.POST("/submit", controllers.SubmitExam)
		}

		// 考试记录相关路由
		record := protected.Group("/records")
		{
			record.GET("/", controllers.ListExamRecords)
			record.GET("/:id", controllers.GetExamRecord)
			record.GET("/stats", controllers.GetExamStats)
		}

		// 模拟练习相关路由
		practice := protected.Group("/practice")
		{
			// 获取练习题目
			practice.GET("/questions", controllers.GetPracticeQuestions)
			// 提交练习答案
			practice.POST("/submit", controllers.SubmitPractice)
		}

		// 错题本相关路由
		wrong := protected.Group("/wrong-questions")
		{
			// 获取错题列表
			wrong.GET("/", controllers.ListWrongQuestions)
			// 移除错题
			wrong.DELETE("/:id", controllers.RemoveWrongQuestion)
		}

		// 管理员功能路由组 - 只有管理员可以访问
		admin := protected.Group("/admin")
		admin.Use(middleware.RoleAuth("admin", "manager"))
		{
			// 这里可以添加更多管理员功能
		}
	}

	// 前端页面路由
	// 首页 - 登录成功后显示
	router.GET("/", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		userName := ""
		if exists {
			// 获取当前用户信息
			user, err := userService.GetUserByID(userID.(int))
			if err == nil && user.Name != "" {
				userName = user.Name
			}
		}
		c.HTML(200, "index.html", gin.H{
			"title":    "基层三基考试系统",
			"userName": userName,
			"stats": gin.H{
				"totalExams":    0,
				"avgScore":      0,
				"passedExams":   0,
				"upcomingExams": 0,
			},
			"recentExams": []interface{}{},
		})
	})

	// 登录页面
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{
			"title": "登录 - 基层三基考试系统",
		})
	})

	// 注册页面
	router.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{
			"title": "注册 - 基层三基考试系统",
		})
	})

	// 登出路由
	router.GET("/logout", func(c *gin.Context) {
		// 清除Cookie
		c.SetCookie("token", "", -1, "/", "", false, true)
		// 重定向到登录页面
		c.Redirect(http.StatusFound, "/login")
	})

	// 考试记录页面
	router.GET("/records", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		userName := ""
		if exists {
			// 获取当前用户信息
			user, err := userService.GetUserByID(userID.(int))
			if err == nil {
				if user.Name != "" {
					userName = user.Name
				} else {
					userName = user.Username
				}
			}
		}
		c.HTML(200, "index.html", gin.H{
			"title":    "考试记录 - 基层三基考试系统",
			"userName": userName,
		})
	})

	// 个人中心页面
	router.GET("/profile", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		// 获取当前用户信息
		user, err := userService.GetUserByID(userID.(int))
		if err != nil {
			c.HTML(200, "profile.html", gin.H{
				"title":    "个人中心 - 基层三基考试系统",
				"error":    "获取用户信息失败",
				"user":     models.User{},
				"userName": "",
			})
			return
		}

		userName := ""
		if user.Name != "" {
			userName = user.Name
		} else {
			userName = user.Username
		}

		c.HTML(200, "profile.html", gin.H{
			"title":    "个人中心 - 基层三基考试系统",
			"user":     user,
			"userName": userName,
		})
	})

	// 模拟练习页面
	router.GET("/practice", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		userName := ""
		if exists {
			// 获取当前用户信息
			user, err := userService.GetUserByID(userID.(int))
			if err == nil {
				if user.Name != "" {
					userName = user.Name
				} else {
					userName = user.Username
				}
			}
		}
		c.HTML(200, "index.html", gin.H{
			"title":    "模拟练习 - 基层三基考试系统",
			"userName": userName,
		})
	})

	// 错题本页面
	router.GET("/wrong-questions", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		userName := ""
		if exists {
			// 获取当前用户信息
			user, err := userService.GetUserByID(userID.(int))
			if err == nil {
				if user.Name != "" {
					userName = user.Name
				} else {
					userName = user.Username
				}
			}
		}
		c.HTML(200, "index.html", gin.H{
			"title":    "错题本 - 基层三基考试系统",
			"userName": userName,
		})
	})

	// 管理员后台页面
	router.GET("/admin", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		userName := ""
		if exists {
			// 获取当前用户信息
			user, err := userService.GetUserByID(userID.(int))
			if err == nil {
				if user.Name != "" {
					userName = user.Name
				} else {
					userName = user.Username
				}
			}
		}
		c.HTML(200, "index.html", gin.H{
			"title":    "管理员后台 - 基层三基考试系统",
			"userName": userName,
		})
	})

	// 管理员用户列表页面
	router.GET("/admin/users", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		userName := ""
		if exists {
			// 获取当前用户信息
			user, err := userService.GetUserByID(userID.(int))
			if err == nil {
				if user.Name != "" {
					userName = user.Name
				} else {
					userName = user.Username
				}
			}
		}
		// 获取用户列表数据
		users := []models.User{}
		// 模拟分页数据
		currentPage := 1
		totalPages := 1
		pages := []int{1}
		c.HTML(200, "admin_users.html", gin.H{
			"title":       "用户管理 - 管理员后台",
			"userName":    userName,
			"users":       users,
			"currentPage": currentPage,
			"totalPages":  totalPages,
			"pages":       pages,
		})
	})

	// 考试统计页面
	router.GET("/stats", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		userName := ""
		if exists {
			// 获取当前用户信息
			user, err := userService.GetUserByID(userID.(int))
			if err == nil {
				if user.Name != "" {
					userName = user.Name
				} else {
					userName = user.Username
				}
			}
		}
		// 获取考试统计数据
		stats := gin.H{
			"totalParticipants": 100,
			"avgScore":          75.5,
			"maxScore":          98,
			"minScore":          45,
			"passRate":          85,
			"excellentRate":     25,
		}
		c.HTML(200, "stats.html", gin.H{
			"title":    "考试统计 - 基层三基考试系统",
			"userName": userName,
			"stats":    stats,
		})
	})

	// 我的考试页面
	router.GET("/exams", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		userName := ""
		if exists {
			// 获取当前用户信息
			user, err := userService.GetUserByID(userID.(int))
			if err == nil {
				if user.Name != "" {
					userName = user.Name
				} else {
					userName = user.Username
				}
			}
		}
		// 获取考试数据
		upcomingExams := []models.Exam{}
		pastExams := []models.ExamRecord{}

		c.HTML(200, "exams.html", gin.H{
			"title":         "我的考试 - 基层三基考试系统",
			"userName":      userName,
			"upcomingExams": upcomingExams,
			"pastExams":     pastExams,
		})
	})

	// 考试详情页面
	router.GET("/exam/:id/details", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "考试详情 - 基层三基考试系统",
		})
	})

	// 开始考试页面
	router.GET("/exam/:id/start", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "开始考试 - 基层三基考试系统",
		})
	})

	// 考试结果页面
	router.GET("/record/:id", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "考试结果 - 基层三基考试系统",
		})
	})

	// 试卷查看页面
	router.GET("/exam/:id/paper", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "查看试卷 - 基层三基考试系统",
		})
	})

	return router
}
