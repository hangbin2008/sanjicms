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
	// 配置验证码生成器 - 使用字符串驱动，避免数字驱动可能的畸变问题
	driver := base64Captcha.NewDriverString(
		48,  // 高度
		120, // 宽度
		0,   // 噪声
		base64Captcha.Option{
			UseNoise:        true,         // 使用噪声
			NoiseLevel:      0.2,          // 噪声水平 - 低噪声
			UseSineLine:     false,        // 不使用正弦线干扰
			UseFont:         true,         // 使用字体
			Fonts:           []string{},   // 使用默认字体
			Width:           120,          // 宽度
			Height:          48,           // 高度
			FontSize:        32,           // 字体大小
			CharSet:         "0123456789", // 字符集 - 只使用数字
			Length:          4,            // 字符数
			ShowLineOptions: 0,            // 不显示干扰线
			BkgColor:        nil,          // 默认背景色
			TextColor:       nil,          // 默认文字颜色
		},
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
