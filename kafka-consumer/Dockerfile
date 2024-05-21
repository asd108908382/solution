FROM ubuntu:latest
LABEL authors="asd108908372"
# 使用官方Go镜像作为构建环境
FROM golang:1.2.2 AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目源码
COPY . .

# 编译构建应用程序
RUN go build -o /kafka-consumer .

# 创建最终镜像，用于运行应用程序
FROM alpine

# 从构建阶段复制编译好的应用程序
COPY --from=builder /kafka-consumer /kafka-consumer

# 设置容器启动时运行的命令
ENTRYPOINT ["/kafka-consumer"]