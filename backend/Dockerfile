# 使用官方Go镜像作为构建环境
FROM golang:1.22.2-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件，并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译可执行文件
RUN go build -o main .

# 使用一个更小的镜像来运行应用
FROM alpine:latest

WORKDIR /root/

# 复制编译好的可执行文件
COPY --from=builder /app/main .

# 复制.env文件（如果有）
COPY --from=builder /app/.env .

# 暴露端口
EXPOSE 8080

# 运行可执行文件
CMD ["./main"]