# 使用官方Go镜像作为构建环境
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖（使用国内代理加速）
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

# 复制源代码
COPY . .

# 构建应用程序
RUN go build -o main ./cmd/server

# 使用轻量级Alpine镜像作为运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制构建好的二进制文件
COPY --from=builder /app/main .

# 从builder阶段复制静态文件和模板
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/migrations ./migrations

# 设置环境变量
ENV GIN_MODE=release

# 暴露端口
EXPOSE 8080

# 启动应用程序
CMD ["./main"]