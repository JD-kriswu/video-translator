#!/bin/bash

# Video Translator 服务管理脚本

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$PROJECT_DIR"

APP_NAME="video-translator-server"
PID_FILE="$PROJECT_DIR/.server.pid"
LOG_FILE="$PROJECT_DIR/server.log"
PORT=8080

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 获取进程 PID
get_pid() {
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            echo "$pid"
            return 0
        fi
    fi
    # 尝试通过端口查找
    local pid=$(lsof -ti:$PORT 2>/dev/null)
    echo "$pid"
}

# 启动服务
start() {
    local pid=$(get_pid)
    if [ -n "$pid" ]; then
        echo -e "${YELLOW}⚠️  服务已在运行 (PID: $pid)${NC}"
        return 1
    fi

    if [ ! -f "$PROJECT_DIR/$APP_NAME" ]; then
        echo -e "${RED}❌ 可执行文件不存在，请先运行 ./build.sh${NC}"
        return 1
    fi

    echo -e "${GREEN}🚀 启动服务...${NC}"
    
    # 设置生产模式
    export GIN_MODE=release
    
    nohup "$PROJECT_DIR/$APP_NAME" > "$LOG_FILE" 2>&1 &
    local new_pid=$!
    echo "$new_pid" > "$PID_FILE"
    
    sleep 1
    
    if ps -p "$new_pid" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 服务启动成功${NC}"
        echo -e "   PID:  $new_pid"
        echo -e "   地址: ${YELLOW}http://localhost:$PORT${NC}"
        echo -e "   日志: $LOG_FILE"
    else
        echo -e "${RED}❌ 服务启动失败，请查看日志: $LOG_FILE${NC}"
        rm -f "$PID_FILE"
        return 1
    fi
}

# 停止服务
stop() {
    local pid=$(get_pid)
    if [ -z "$pid" ]; then
        echo -e "${YELLOW}⚠️  服务未运行${NC}"
        rm -f "$PID_FILE"
        return 0
    fi

    echo -e "${GREEN}🛑 停止服务 (PID: $pid)...${NC}"
    kill "$pid" 2>/dev/null
    
    # 等待进程结束
    local count=0
    while ps -p "$pid" > /dev/null 2>&1 && [ $count -lt 10 ]; do
        sleep 1
        count=$((count + 1))
    done
    
    if ps -p "$pid" > /dev/null 2>&1; then
        echo -e "${YELLOW}强制终止进程...${NC}"
        kill -9 "$pid" 2>/dev/null
    fi
    
    rm -f "$PID_FILE"
    echo -e "${GREEN}✅ 服务已停止${NC}"
}

# 重启服务
restart() {
    stop
    sleep 1
    start
}

# 查看状态
status() {
    local pid=$(get_pid)
    if [ -n "$pid" ]; then
        echo -e "${GREEN}✅ 服务运行中${NC}"
        echo -e "   PID:  $pid"
        echo -e "   地址: http://localhost:$PORT"
        echo -e "   日志: $LOG_FILE"
    else
        echo -e "${RED}❌ 服务未运行${NC}"
    fi
}

# 查看日志
logs() {
    if [ -f "$LOG_FILE" ]; then
        tail -f "$LOG_FILE"
    else
        echo -e "${YELLOW}⚠️  日志文件不存在${NC}"
    fi
}

# 使用帮助
usage() {
    echo "Video Translator 服务管理脚本"
    echo ""
    echo "用法: $0 {start|stop|restart|status|logs}"
    echo ""
    echo "命令:"
    echo "  start    启动服务"
    echo "  stop     停止服务"
    echo "  restart  重启服务"
    echo "  status   查看状态"
    echo "  logs     查看日志"
}

# 主逻辑
case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    logs)
        logs
        ;;
    *)
        usage
        exit 1
        ;;
esac
