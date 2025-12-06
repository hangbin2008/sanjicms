package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hangbin2008/sanjicms/pkg/config"
)

var DB *sql.DB

func Connect(config *config.DatabaseConfig) error {
	var err error

	// 修复注册失败问题：确保数据库存在
	// 1. 首先创建一个不指定数据库的DSN，用于创建数据库
	tempDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Charset)

	// 2. 打开临时连接（不指定数据库）
	tempDB, err := sql.Open("mysql", tempDSN)
	if err != nil {
		return fmt.Errorf("failed to open temp database connection: %w", err)
	}
	defer tempDB.Close()

	// 3. 创建数据库（如果不存在）
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET %s COLLATE %s_unicode_ci",
		config.DBName, config.Charset, config.Charset)
	if _, err := tempDB.Exec(createDBQuery); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// 4. 使用指定数据库创建正式连接
	dsn := config.GetDSN()
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// 5. 设置连接池参数
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 6. 测试连接
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connection established successfully to database: %s", config.DBName)
	return nil
}

func Close() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}
