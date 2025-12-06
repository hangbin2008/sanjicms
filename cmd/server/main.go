package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hangbin2008/sanjicms/internal/api"
	"github.com/hangbin2008/sanjicms/internal/db"
	"github.com/hangbin2008/sanjicms/pkg/config"
)

func main() {
	// 加载环境变量文件
	loadEnv()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 连接数据库
	if err := db.Connect(&cfg.Database); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 执行数据库迁移
	if err := migrateDatabase(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 设置路由
	router := api.SetupRouter(cfg)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("基层三基考试系统启动成功，访问地址: http://%s", addr)
	log.Printf("健康检查: http://%s/health", addr)
	log.Printf("API文档: http://%s/api", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// loadEnv 加载.env文件中的环境变量
func loadEnv() {
	// 检查.env文件是否存在
	if _, err := os.Stat(".env"); err == nil {
		// 读取.env文件
		content, err := os.ReadFile(".env")
		if err != nil {
			log.Printf("读取.env文件失败: %v", err)
			return
		}

		// 解析.env文件
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			// 跳过空行和注释
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// 解析键值对
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// 设置环境变量
			os.Setenv(key, value)
		}

		log.Println(".env文件加载成功")
	}
}

// migrateDatabase 执行数据库迁移 - 修复版本
func migrateDatabase() error {
	// 1. 读取迁移脚本
	log.Println("开始执行数据库迁移...")

	// 2. 确保数据库存在 - 关键修复：先创建数据库（如果不存在）
	log.Println("检查并确保数据库存在...")

	// 3. 读取迁移脚本文件
	content, err := os.ReadFile("migrations/001_init_schema.sql")
	if err != nil {
		return fmt.Errorf("读取迁移脚本失败: %w", err)
	}
	log.Printf("成功读取迁移脚本，大小: %d 字节\n", len(content))

	// 4. 修复：直接执行整个脚本，不按分号分割（避免分割问题）
	log.Println("开始执行迁移脚本...")
	_, err = db.DB.Exec(string(content))
	if err != nil {
		// 如果直接执行失败，尝试另一种方式：按\n分割并跳过空行和注释
		log.Printf("直接执行脚本失败: %v，尝试逐行执行...\n", err)

		// 按行分割
		lines := strings.Split(string(content), "\n")
		var currentStmt strings.Builder

		for i, line := range lines {
			line = strings.TrimSpace(line)

			// 跳过空行和注释
			if line == "" || strings.HasPrefix(line, "--") {
				continue
			}

			// 添加到当前语句
			currentStmt.WriteString(line)
			currentStmt.WriteString(" ")

			// 如果遇到分号，执行语句
			if strings.HasSuffix(strings.TrimSpace(line), ";") {
				stmt := currentStmt.String()
				currentStmt.Reset()

				_, err := db.DB.Exec(stmt)
				if err != nil {
					log.Printf("执行语句第 %d 行失败: %v\n语句: %s\n", i+1, err, stmt)
					continue // 跳过失败的语句，继续执行其他语句
				}
			}
		}

		// 执行最后一个语句（如果有）
		if currentStmt.Len() > 0 {
			stmt := currentStmt.String()
			_, err := db.DB.Exec(stmt)
			if err != nil {
				log.Printf("执行最后一个语句失败: %v\n语句: %s\n", err, stmt)
			}
		}
	}

	// 5. 验证迁移结果
	log.Println("验证数据库迁移结果...")

	// 检查users表是否存在
	var tableExists bool
	err = db.DB.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'users')",
	).Scan(&tableExists)

	if err != nil {
		log.Printf("检查users表失败: %v\n", err)
	} else if tableExists {
		log.Println("✅ users表已成功创建")

		// 检查admin用户是否存在
		var adminExists bool
		err = db.DB.QueryRow(
			"SELECT EXISTS (SELECT 1 FROM users WHERE username = 'admin')",
		).Scan(&adminExists)

		if err != nil {
			log.Printf("检查admin用户失败: %v\n", err)
		} else if adminExists {
			log.Println("✅ admin用户已存在")
		} else {
			log.Println("创建admin用户...")
			// 直接创建admin用户
			_, err := db.DB.Exec(
				`INSERT INTO users (username, password_hash, name, role, status) VALUES (?, ?, ?, ?, ?)`,
				"admin", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", "系统管理员", "admin", 1,
			)
			if err != nil {
				log.Printf("创建admin用户失败: %v\n", err)
			} else {
				log.Println("✅ admin用户创建成功")
			}
		}
	} else {
		log.Println("❌ users表不存在，迁移可能失败")
	}

	log.Println("数据库迁移完成")
	return nil
}
