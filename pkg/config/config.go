package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Password PasswordConfig
}

type AppConfig struct {
	Name string
}

type ServerConfig struct {
	Port         int
	Host         string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Charset  string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn int
}

type PasswordConfig struct {
	MinLength      int
	RequireLetter  bool
	RequireDigit   bool
	RequireSpecial bool
}

func Load() (*Config, error) {
	config := &Config{}

	// App config
	config.App.Name = getEnv("APP_NAME", "jiceng-sanji-exam")

	// Server config
	config.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
	config.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	config.Server.ReadTimeout = getEnvAsInt("READ_TIMEOUT", 15)
	config.Server.WriteTimeout = getEnvAsInt("WRITE_TIMEOUT", 15)

	// Database config
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnvAsInt("DB_PORT", 3306)
	config.Database.User = getEnv("DB_USER", "root")
	config.Database.Password = getEnv("DB_PASSWORD", "password")
	// 如果没有设置DB_NAME，使用APP_NAME作为数据库名称
	if dbName := getEnv("DB_NAME", ""); dbName != "" {
		config.Database.DBName = dbName
	} else {
		config.Database.DBName = config.App.Name
	}
	config.Database.Charset = getEnv("DB_CHARSET", "utf8mb4")

	// JWT config
	config.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key")
	config.JWT.ExpiresIn = getEnvAsInt("JWT_EXPIRES_IN", 3600)

	// Password config
	config.Password.MinLength = getEnvAsInt("PASSWORD_MIN_LENGTH", 8)
	config.Password.RequireLetter = getEnvAsBool("PASSWORD_REQUIRE_LETTER", true)
	config.Password.RequireDigit = getEnvAsBool("PASSWORD_REQUIRE_DIGIT", true)
	config.Password.RequireSpecial = getEnvAsBool("PASSWORD_REQUIRE_SPECIAL", true)

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}
