# 多阶段构建
FROM golang:1.25-alpine AS builder

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

# 复制必要的文件
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
# 确保static目录存在
RUN mkdir -p ./static
# 复制static目录内容（如果存在）
COPY --from=builder /app/static/* ./static/
COPY --from=builder /app/migrations ./migrations
# 复制.env文件
COPY --from=builder /app/.env ./

# 设置环境变量
ENV GIN_MODE=release

# 暴露端口
EXPOSE 8080

# 启动应用程序
CMD ["./main"]