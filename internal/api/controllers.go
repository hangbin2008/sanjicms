package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hangbin2008/sanjicms/internal/models"
	"github.com/hangbin2008/sanjicms/internal/service"
)

// Controllers API控制器
type Controllers struct {
	userService     *service.UserService
	questionService *service.QuestionService
	examService     *service.ExamService
}

// NewControllers 创建API控制器
func NewControllers(
	userService *service.UserService,
	questionService *service.QuestionService,
	examService *service.ExamService,
) *Controllers {
	return &Controllers{
		userService:     userService,
		questionService: questionService,
		examService:     examService,
	}
}

// Register 用户注册
func (c *Controllers) Register(ctx *gin.Context) {
	var req models.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.RegisterUser(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "用户注册成功",
		"user": models.UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			Name:       user.Name,
			Role:       user.Role,
			Phone:      user.Phone,
			IDCard:     user.IDCard,
			Department: user.Department,
			JobTitle:   user.JobTitle,
			Status:     user.Status,
			CreatedAt:  user.CreatedAt,
		},
	})
}

// Login 用户登录
func (c *Controllers) Login(ctx *gin.Context) {
	var req models.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.userService.LoginUser(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data":    response,
	})
}

// GetCurrentUser 获取当前用户信息
func (c *Controllers) GetCurrentUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	user, err := c.userService.GetUserByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取用户信息成功",
		"user": models.UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			Name:       user.Name,
			Role:       user.Role,
			Phone:      user.Phone,
			IDCard:     user.IDCard,
			Department: user.Department,
			JobTitle:   user.JobTitle,
			Avatar:     user.Avatar,
			Status:     user.Status,
			CreatedAt:  user.CreatedAt,
		},
	})
}

// UpdateUser 更新用户信息
func (c *Controllers) UpdateUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req models.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.UpdateUser(userID.(int), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "更新用户信息成功",
		"user": models.UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			Name:       user.Name,
			Role:       user.Role,
			Phone:      user.Phone,
			IDCard:     user.IDCard,
			Department: user.Department,
			JobTitle:   user.JobTitle,
			Avatar:     user.Avatar,
			Status:     user.Status,
			CreatedAt:  user.CreatedAt,
		},
	})
}

// CreateQuestionBank 创建题库
func (c *Controllers) CreateQuestionBank(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req models.QuestionBankCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bank, err := c.questionService.CreateQuestionBank(&req, userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "题库创建成功",
		"bank":    bank,
	})
}

// GetQuestionBankByID 根据ID获取题库
func (c *Controllers) GetQuestionBankByID(ctx *gin.Context) {
	bankID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的题库ID"})
		return
	}

	bank, err := c.questionService.GetQuestionBankByID(bankID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "题库不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取题库成功",
		"bank":    bank,
	})
}

// ListQuestionBanks 获取题库列表
func (c *Controllers) ListQuestionBanks(ctx *gin.Context) {
	subject := ctx.Query("subject")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	banks, total, err := c.questionService.ListQuestionBanks(subject, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取题库列表成功",
		"data": gin.H{
			"banks":     banks,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateQuestion 创建题目
func (c *Controllers) CreateQuestion(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req models.QuestionCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := c.questionService.CreateQuestion(&req, userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "题目创建成功",
		"question": question,
	})
}

// GetQuestionByID 根据ID获取题目
func (c *Controllers) GetQuestionByID(ctx *gin.Context) {
	questionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	question, err := c.questionService.GetQuestionByID(questionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "获取题目成功",
		"question": question,
	})
}

// ListQuestionsByBank 获取题库下的题目列表
func (c *Controllers) ListQuestionsByBank(ctx *gin.Context) {
	bankID, err := strconv.Atoi(ctx.Param("bank_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的题库ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	questions, total, err := c.questionService.ListQuestionsByBank(bankID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取题目列表成功",
		"data": gin.H{
			"questions": questions,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GenerateExam 生成试卷
func (c *Controllers) GenerateExam(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req models.ExamGenerateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exam, err := c.examService.GenerateExam(&req, userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "试卷生成成功",
		"exam":    exam,
	})
}

// GetExamByID 根据ID获取试卷
func (c *Controllers) GetExamByID(ctx *gin.Context) {
	examID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的试卷ID"})
		return
	}

	exam, err := c.examService.GetExamByID(examID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "试卷不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取试卷成功",
		"exam":    exam,
	})
}

// StartExam 开始考试
func (c *Controllers) StartExam(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	examID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的试卷ID"})
		return
	}

	record, err := c.examService.StartExam(examID, userID.(int))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "开始考试成功",
		"record":  record,
	})
}

// SubmitExam 提交试卷
func (c *Controllers) SubmitExam(ctx *gin.Context) {
	var req models.ExamSubmitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := c.examService.SubmitExam(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "提交试卷成功",
		"record":  record,
	})
}

// GetExamRecord 获取考试记录
func (c *Controllers) GetExamRecord(ctx *gin.Context) {
	recordID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	record, err := c.examService.GetExamRecord(recordID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取考试记录成功",
		"record":  record,
	})
}

// ListExamRecords 获取用户考试记录
func (c *Controllers) ListExamRecords(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	records, total, err := c.examService.ListExamRecords(userID.(int), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取考试记录成功",
		"data": gin.H{
			"records":   records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetExamStats 获取考试统计数据
func (c *Controllers) GetExamStats(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	stats, err := c.examService.GetExamStats(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取考试统计成功",
		"stats":   stats,
	})
}

// ListExams 获取试卷列表
func (c *Controllers) ListExams(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	exams, total, err := c.examService.ListExams(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取试卷列表成功",
		"data": gin.H{
			"exams":     exams,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetPracticeQuestions 获取练习题目
func (c *Controllers) GetPracticeQuestions(ctx *gin.Context) {
	// 这里需要实现GetPracticeQuestions方法
	ctx.JSON(http.StatusOK, gin.H{
		"message":   "获取练习题目成功",
		"questions": []models.Question{},
	})
}

// SubmitPractice 提交练习答案
func (c *Controllers) SubmitPractice(ctx *gin.Context) {
	// 这里需要实现SubmitPractice方法
	ctx.JSON(http.StatusOK, gin.H{
		"message": "提交练习答案成功",
		"score":   0,
		"correct": 0,
		"total":   0,
	})
}

// ListWrongQuestions 获取错题列表
func (c *Controllers) ListWrongQuestions(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	// 这里需要实现ListWrongQuestions方法
	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取错题列表成功",
		"data": gin.H{
			"questions": []models.Question{},
			"total":     0,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// RemoveWrongQuestion 移除错题
func (c *Controllers) RemoveWrongQuestion(ctx *gin.Context) {
	// 这里需要实现RemoveWrongQuestion方法
	ctx.JSON(http.StatusOK, gin.H{
		"message": "移除错题成功",
	})
}
