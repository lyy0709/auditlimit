#!/bin/bash

set -e

# 构建 linux/amd64 架构的二进制文件
gf build main.go -a amd64 -s linux -p ./temp

# 使用 docker buildx 构建 amd64 架构的镜像并直接推送（在 macOS ARM 上交叉编译）
# --provenance=false 和 --sbom=false 禁用 attestation，确保生成纯 amd64 镜像
now=$(date +"%Y%m%d%H%M%S")

docker buildx build --platform linux/amd64 \
  --provenance=false \
  --sbom=false \
  -f manifest/docker/Dockerfile \
  -t lyy0709/auditlimit:latest \
  -t lyy0709/auditlimit:$now \
  --push .
echo "release success" $now
