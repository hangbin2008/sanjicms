package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hangbin2008/sanjicms/internal/middleware"
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

	// 创建控制器实例
	controllers := NewControllers(userService, questionService, examService)

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
		c.HTML(200, "index.html", gin.H{
			"title": "基层三基考试系统",
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

	// 考试相关页面
	router.GET("/exams", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "我的考试 - 基层三基考试系统",
		})
	})

	// 考试记录页面
	router.GET("/records", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "考试记录 - 基层三基考试系统",
		})
	})

	// 个人中心页面
	router.GET("/profile", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "个人中心 - 基层三基考试系统",
		})
	})

	// 模拟练习页面
	router.GET("/practice", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "模拟练习 - 基层三基考试系统",
		})
	})

	// 错题本页面
	router.GET("/wrong-questions", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "错题本 - 基层三基考试系统",
		})
	})

	// 管理员后台页面
	router.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "管理员后台 - 基层三基考试系统",
		})
	})

	// 考试统计页面
	router.GET("/stats", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "考试统计 - 基层三基考试系统",
		})
	})

	return router
}
