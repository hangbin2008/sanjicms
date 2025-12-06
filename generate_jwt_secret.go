package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	// 定义密钥长度（字节），建议至少32字节（256位）
	keyLength := 32
	
	// 生成随机字节
	randomBytes := make([]byte, keyLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成随机密钥失败: %v\n", err)
		os.Exit(1)
	}
	
	// 转换为base64编码字符串
	jwtSecret := base64.StdEncoding.EncodeToString(randomBytes)
	
	fmt.Println("生成的JWT强密钥:")
	fmt.Println(jwtSecret)
	fmt.Printf("\n密钥长度: %d 字节 (%d 位)\n", keyLength, keyLength*8)
	fmt.Println("\n使用方法:")
	fmt.Println("1. 将此密钥复制到.env文件中的JWT_SECRET字段")
	fmt.Printf("2. 或作为环境变量设置：export JWT_SECRET='%s'\n", jwtSecret)
}
