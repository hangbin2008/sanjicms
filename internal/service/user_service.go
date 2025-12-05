package service

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/yourusername/jiceng-sanji-exam/internal/db"
	"github.com/yourusername/jiceng-sanji-exam/internal/middleware"
	"github.com/yourusername/jiceng-sanji-exam/internal/models"
	"github.com/yourusername/jiceng-sanji-exam/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	config    *config.Config
	jwtConfig *middleware.JWTConfig
}

// NewUserService 创建用户服务
func NewUserService(cfg *config.Config, jwt *middleware.JWTConfig) *UserService {
	return &UserService{
		config:    cfg,
		jwtConfig: jwt,
	}
}

// ValidatePassword 验证密码是否符合要求
func (s *UserService) ValidatePassword(password string) error {
	// 检查密码长度
	if len(password) < s.config.Password.MinLength {
		return errors.New("密码长度不足")
	}

	// 检查是否包含字母
	if s.config.Password.RequireLetter {
		hasLetter := false
		for _, r := range password {
			if unicode.IsLetter(r) {
				hasLetter = true
				break
			}
		}
		if !hasLetter {
			return errors.New("密码必须包含字母")
		}
	}

	// 检查是否包含数字
	if s.config.Password.RequireDigit {
		hasDigit := false
		for _, r := range password {
			if unicode.IsDigit(r) {
				hasDigit = true
				break
			}
		}
		if !hasDigit {
			return errors.New("密码必须包含数字")
		}
	}

	// 检查是否包含特殊字符
	if s.config.Password.RequireSpecial {
		hasSpecial := false
		specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
		for _, r := range password {
			if strings.ContainsRune(specialChars, r) {
				hasSpecial = true
				break
			}
		}
		if !hasSpecial {
			return errors.New("密码必须包含特殊字符")
		}
	}

	return nil
}

// ValidateName 验证姓名是否只包含中文
func (s *UserService) ValidateName(name string) error {
	// 检查姓名是否只包含中文
	match, _ := regexp.MatchString("^[\p{Han}]+$", name)
	if !match {
		return errors.New("姓名只能包含中文")
	}
	return nil
}

// RegisterUser 用户注册
func (s *UserService) RegisterUser(req *models.UserRegisterRequest) (*models.User, error) {
	// 验证姓名
	if err := s.ValidateName(req.Name); err != nil {
		return nil, err
	}

	// 验证密码
	if err := s.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE phone = ?", req.Phone).Scan(&count)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("手机号已存在")
		}
	}

	// 检查身份证号是否已存在
	if req.IDCard != "" {
		err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id_card = ?", req.IDCard).Scan(&count)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("身份证号已存在")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 设置默认角色
	role := req.Role
	if role == "" {
		role = "employee"
	}

	// 插入用户记录
	result, err := db.DB.Exec(`
		INSERT INTO users (username, password_hash, name, role, phone, id_card, department, job_title)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Username, string(hashedPassword), req.Name, role, req.Phone, req.IDCard, req.Department, req.JobTitle)
	if err != nil {
		return nil, err
	}

	// 获取插入的用户ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 查询插入的用户信息
	var user models.User
	err = db.DB.QueryRow(`
		SELECT id, username, name, role, phone, id_card, department, job_title, status, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(
		&user.ID, &user.Username, &user.Name, &user.Role, &user.Phone, &user.IDCard,
		&user.Department, &user.JobTitle, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// LoginUser 用户登录
func (s *UserService) LoginUser(req *models.UserLoginRequest) (*models.LoginResponse, error) {
	// 查询用户
	var user models.User
	err := db.DB.QueryRow(`
		SELECT id, username, password_hash, name, role, phone, id_card, department, job_title, status, created_at, updated_at
		FROM users WHERE username = ?
	`, req.Username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Name, &user.Role, &user.Phone, &user.IDCard,
		&user.Department, &user.JobTitle, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := s.jwtConfig.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &models.LoginResponse{
		User: models.UserResponse{
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
		Token: token,
	}

	return response, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := db.DB.QueryRow(`
		SELECT id, username, name, role, phone, id_card, department, job_title, status, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(
		&user.ID, &user.Username, &user.Name, &user.Role, &user.Phone, &user.IDCard,
		&user.Department, &user.JobTitle, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
