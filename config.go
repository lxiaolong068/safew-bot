package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用程序配置结构
type Config struct {
	BotToken          string
	ForwardTargetChat int64
	LogLevel          string
	SuperAdmins       []int64
	PollTimeout       int
}

// LoadConfig 从环境变量和.env文件加载配置
func LoadConfig() (*Config, error) {
	// 尝试加载.env文件（如果存在）
	if err := godotenv.Load(); err != nil {
		log.Printf("Info: .env file not found or cannot be loaded: %v", err)
		log.Println("Info: Will use system environment variables")
	} else {
		log.Println("Info: .env file loaded successfully")
	}

	config := &Config{
		LogLevel:    "INFO",
		PollTimeout: 30, // 默认30秒超时
	}

	// 必需的配置项
	botToken := os.Getenv("SAFEW_BOT_TOKEN")
	if botToken == "" {
		return nil, errors.New("SAFEW_BOT_TOKEN environment variable is required")
	}
	config.BotToken = botToken

	// 可选配置项
	if forwardTarget := os.Getenv("FORWARD_TARGET_CHAT"); forwardTarget != "" {
		if chatID, err := strconv.ParseInt(forwardTarget, 10, 64); err == nil {
			config.ForwardTargetChat = chatID
		} else {
			log.Printf("Warning: Invalid FORWARD_TARGET_CHAT value: %s", forwardTarget)
		}
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	if timeout := os.Getenv("POLL_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil && t > 0 {
			config.PollTimeout = t
		} else {
			log.Printf("Warning: Invalid POLL_TIMEOUT value: %s, using default", timeout)
		}
	}

	// 超级管理员配置
	if adminIDs := os.Getenv("SUPER_ADMINS"); adminIDs != "" {
		// 简单的逗号分隔解析，后续可以改进
		// 格式: "123456789,987654321"
		// 这里暂时不实现复杂解析，后续可以扩展
		log.Printf("SUPER_ADMINS configuration detected, parsing not implemented yet")
	}

	return config, nil
}

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	if c.BotToken == "" {
		return errors.New("bot token cannot be empty")
	}

	if c.PollTimeout <= 0 {
		return errors.New("poll timeout must be positive")
	}

	// 验证日志级别
	validLogLevels := map[string]bool{
		"DEBUG": true,
		"INFO":  true,
		"WARN":  true,
		"ERROR": true,
	}

	if !validLogLevels[c.LogLevel] {
		return errors.New("invalid log level: must be DEBUG, INFO, WARN, or ERROR")
	}

	return nil
}

// IsSuperAdmin 检查用户是否为超级管理员
func (c *Config) IsSuperAdmin(userID int64) bool {
	for _, adminID := range c.SuperAdmins {
		if adminID == userID {
			return true
		}
	}
	return false
} 