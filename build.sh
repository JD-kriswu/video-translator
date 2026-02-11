#!/bin/bash

# Video Translator 编译脚本

set -e

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$PROJECT_DIR"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🔨 Video Translator 编译脚本${NC}"
echo "================================"

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go 未安装，请先安装 Go 1.21+${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "Go 版本: ${YELLOW}$GO_VERSION${NC}"

# 安装依赖
echo -e "\n${GREEN}📦 安装依赖...${NC}"
go mod tidy

# 编译 CLI 版本
echo -e "\n${GREEN}🔧 编译 CLI 版本...${NC}"
go build -ldflags="-s -w" -o video-translator ./cmd/main.go
echo -e "✅ video-translator 编译完成"

# 编译 Web 服务版本
echo -e "\n${GREEN}🔧 编译 Web 服务版本...${NC}"
go build -ldflags="-s -w" -o video-translator-server ./cmd/server/main.go
echo -e "✅ video-translator-server 编译完成"

# 显示结果
echo -e "\n${GREEN}🎉 编译完成！${NC}"
echo "================================"
ls -lh video-translator video-translator-server

echo -e "\n${YELLOW}使用方式:${NC}"
echo "  CLI:  ./video-translator -url <video-url>"
echo "  Web:  ./manage.sh start"
