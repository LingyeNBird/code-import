# 使用官方的Golang镜像作为构建器
FROM golang:latest AS builder

# 设置工作目录
WORKDIR /app

# 复制所有文件到工作目录
COPY . .

# 设置时区为Asia/Shanghai
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "-X 'ccrctl/cmd.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')' -X 'ccrctl/cmd.Version=$(git rev-parse HEAD)' " -o cnb-code-import ./

# 使用Alpine镜像作为最终镜像
FROM alpine:latest

# 安装Git和Git LFS所需的依赖项
RUN apk update && apk add --no-cache \
    git \
    git-lfs \
    openssh-client \
    openssl \
    ca-certificates \
    tzdata && update-ca-certificates && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 设置工作目录
WORKDIR /app

# 从构建器复制编译好的应用程序和二进制文件
COPY --from=builder /app/cnb-code-import .
COPY --from=builder /app/entrypoint.sh .

# 设置entrypoint脚本可执行权限
RUN chmod +x /app/entrypoint.sh

# 设置容器启动时的入口点脚本
ENTRYPOINT ["/app/entrypoint.sh"]