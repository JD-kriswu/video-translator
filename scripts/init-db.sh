#!/bin/bash

# Video Translator 数据库初始化脚本
# 用法: ./scripts/init-db.sh [mysql_user] [mysql_password]

MYSQL_USER="${1:-root}"
MYSQL_PASS="${2:-}"
MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
DB_NAME="video_translator"

echo "🔧 Video Translator 数据库初始化"
echo "================================"
echo "MySQL Host: $MYSQL_HOST:$MYSQL_PORT"
echo "MySQL User: $MYSQL_USER"
echo "Database:   $DB_NAME"
echo ""

# 构建 MySQL 命令
if [ -n "$MYSQL_PASS" ]; then
    MYSQL_CMD="mysql -h$MYSQL_HOST -P$MYSQL_PORT -u$MYSQL_USER -p$MYSQL_PASS"
else
    MYSQL_CMD="mysql -h$MYSQL_HOST -P$MYSQL_PORT -u$MYSQL_USER"
fi

# 创建数据库
echo "📦 创建数据库..."
$MYSQL_CMD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

if [ $? -ne 0 ]; then
    echo "❌ 创建数据库失败"
    exit 1
fi

# 创建用户表
echo "📋 创建用户表..."
$MYSQL_CMD $DB_NAME << 'EOF'
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3),
    UNIQUE INDEX idx_username (username),
    UNIQUE INDEX idx_email (email),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
EOF

if [ $? -ne 0 ]; then
    echo "❌ 创建表失败"
    exit 1
fi

echo ""
echo "✅ 数据库初始化完成！"
echo ""
echo "下一步："
echo "1. 确保 Redis 已启动: redis-server"
echo "2. 更新 config.json 中的数据库配置"
echo "3. 启动服务: ./video-translator-server"
