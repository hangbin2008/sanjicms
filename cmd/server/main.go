package main

import (
	"fmt"
	"log"

	"github.com/yourusername/jiceng-sanji-exam/internal/api"
	"github.com/yourusername/jiceng-sanji-exam/internal/db"
	"github.com/yourusername/jiceng-sanji-exam/pkg/config"
)

func main() {
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
