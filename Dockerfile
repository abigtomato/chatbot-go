# 使用 golang 官方镜像提供 Go 运行环境，并且命名为 buidler 以便后续引用
FROM golang:1.19-alpine AS builder

# 启用 Go Modules 并设置 GOPROXY
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

# 更新安装源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置工作目录
RUN mkdir /app
ADD . /app/
WORKDIR /app

# 编译 Go 源码
RUN go build -o chatbot-go .

FROM alpine:latest

# 更新安装源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装相关软件
RUN apk update && apk add --no-cache bash supervisor ca-certificates

# 设置时区
ARG TZ="Asia/Shanghai"
ENV TZ ${TZ}

# 设置工作目录
RUN mkdir /app && apk upgrade \
    && apk add bash tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

WORKDIR /app
COPY --from=builder /app/ .
RUN chmod +x chatbot-go && cp config.dev.yaml config.yaml

EXPOSE 8090

CMD ./chatbot-go