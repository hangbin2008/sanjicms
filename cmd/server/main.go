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

// migrateDatabase 执行数据库迁移
func migrateDatabase() error {
	// 读取迁移脚本
	content, err := os.ReadFile("migrations/001_init_schema.sql")
	if err != nil {
		return fmt.Errorf("读取迁移脚本失败: %w", err)
	}

	// 执行迁移脚本
	_, err = db.DB.Exec(string(content))
	if err != nil {
		return fmt.Errorf("执行迁移脚本失败: %w", err)
	}

	log.Println("数据库迁移成功")
	return nil
}
