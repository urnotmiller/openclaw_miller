#!/bin/bash

# 自动生成当天日志文件的脚本

WORKSPACE="/root/.openclaw/workspace"
DATE=$(date +%Y-%m-%d)
LOG_FILE="$WORKSPACE/memory/$DATE.md"

# 检查记忆文件是否存在
if [ ! -f "$LOG_FILE" ]; then
    echo "记忆文件 $LOG_FILE 不存在，不生成日志文件"
    exit 0
fi

# 生成当天的日志文件
LOG_CONTENT=$(cat "$LOG_FILE")

# 创建日志文件
TIMESTAMP=$(date +"%Y-%m-%d %H:%M:%S")
echo "生成时间: $TIMESTAMP" > "$WORKSPACE/daily_log_$DATE.txt"
echo "----------------------------------------" >> "$WORKSPACE/daily_log_$DATE.txt"
echo "$LOG_CONTENT" >> "$WORKSPACE/daily_log_$DATE.txt"

echo "当天日志文件已生成: $WORKSPACE/daily_log_$DATE.txt"

# 将日志文件同步到GitHub仓库
cd "$WORKSPACE" || exit
git add "daily_log_$DATE.txt"
git commit -m "Add daily log for $DATE"
git push origin master

echo "日志文件已同步到GitHub仓库"
