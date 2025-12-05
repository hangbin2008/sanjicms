package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/jiceng-sanji-exam/internal/middleware"
	"github.com/yourusername/jiceng-sanji-exam/internal/service"
	"github.com/yourusername/jiceng-sanji-exam/pkg/config"
)

// SetupRouter 配置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	// 创建Gin引擎
	router := gin.Default()

	// 设置静态文件服务
	router.Static("/static", "./static")

	// 设置模板引擎
	router.LoadHTMLGlob("./templates/*")

	// 创建JWT中间件
	jwtConfig := middleware.NewJWTConfig(cfg)

	// 创建服务实例
	userService := service.NewUserService(cfg, jwtConfig)
	questionService := service.NewQuestionService()
	examService := service.NewExamService(questionService)

	// 创建控制器实例
	controllers := NewControllers(userService, questionService, examService)

	// 健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
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
			// 获取试卷列表和详情（所有用户都可以访问）
			exam.GET("/", controllers.GetExamByID)
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
	}

	// 前端页面路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "基层三基考试系统",
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{
			"title": "登录 - 基层三基考试系统",
		})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{
			"title": "注册 - 基层三基考试系统",
		})
	})

	return router
}
