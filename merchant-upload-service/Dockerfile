# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go模块文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o merchant-upload .

# 运行阶段
FROM alpine:3.18

WORKDIR /app

# 复制构建产物
COPY --from=builder /app/merchant-upload .
COPY --from=builder /app/www ./www

# 创建www目录并设置权限
RUN mkdir -p /app/www && chmod 755 /app/www

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./merchant-upload"]
    