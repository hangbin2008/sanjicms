package service

import (
	"fmt"
	"sync"

	"github.com/mojocn/base64Captcha"
)

// CaptchaService 验证码服务
type CaptchaService struct {
	store sync.Map // 用于存储验证码，key: captcha ID, value: captcha answer
}

// NewCaptchaService 创建验证码服务
func NewCaptchaService() *CaptchaService {
	return &CaptchaService{}
}

// GenerateCaptcha 生成验证码
func (s *CaptchaService) GenerateCaptcha() (string, string, error) {
	// 配置验证码生成器 - 使用正确的NewDriverDigit参数，修复编译错误
	driver := base64Captcha.NewDriverDigit(
		120, // 宽度
		48,  // 高度
		4,   // 字符数
		0.0, // 干扰系数 - 完全无干扰，确保清晰
		0,   // 最大干扰点数 - 完全无干扰
	)

	// 创建验证码实例
	captcha := base64Captcha.NewCaptcha(driver, s)

	// 生成验证码
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return "", "", fmt.Errorf("生成验证码失败: %w", err)
	}

	// 存储验证码答案
	s.store.Store(id, answer)

	return id, b64s, nil
}

// VerifyCaptcha 验证验证码
func (s *CaptchaService) VerifyCaptcha(id, answer string) bool {
	if id == "" || answer == "" {
		return false
	}

	// 从存储中获取验证码答案
	storedAnswer, exists := s.store.Load(id)
	if !exists {
		return false
	}

	// 比较验证码
	if storedAnswer != answer {
		return false
	}

	// 验证成功后删除验证码
	s.store.Delete(id)

	return true
}

// Get 实现base64Captcha.Store接口
func (s *CaptchaService) Get(id string, clear bool) string {
	value, exists := s.store.Load(id)
	if !exists {
		return ""
	}

	if clear {
		s.store.Delete(id)
	}

	if v, ok := value.(string); ok {
		return v
	}

	return ""
}

// Set 实现base64Captcha.Store接口
func (s *CaptchaService) Set(id string, value string) error {
	s.store.Store(id, value)
	return nil
}

// Verify 实现base64Captcha.Store接口
func (s *CaptchaService) Verify(id, answer string, clear bool) bool {
	return s.VerifyCaptcha(id, answer)
}
