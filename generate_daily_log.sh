#!/bin/bash

# 自动生成当天日志文件的脚本
WORKSPACE="/root/.openclaw/workspace"
DATE=$(date +%Y-%m-%d)
MEMORY_DIR="$WORKSPACE/memory"
LOG_DIR="$WORKSPACE/logs/daily_log"

# 确保目录存在
mkdir -p $MEMORY_DIR
mkdir -p $LOG_DIR

# 记忆文件路径
MEMORY_FILE="$MEMORY_DIR/$DATE.md"

# 检查并创建记忆文件（如果不存在）
if [ ! -f "$MEMORY_FILE" ]; then
    echo "# $DATE.md - Daily Log" > "$MEMORY_FILE"
    echo "" >> "$MEMORY_FILE"
    echo "## 当天活动记录" >> "$MEMORY_FILE"
    echo "" >> "$MEMORY_FILE"
    echo "- **日期:** $DATE" >> "$MEMORY_FILE"
    echo "- **主要活动:** 系统自动创建记忆文件" >> "$MEMORY_FILE"
    echo "- **系统状态:** 正常运行" >> "$MEMORY_FILE"
    echo "" >> "$MEMORY_FILE"
    echo "## 备注" >> "$MEMORY_FILE"
    echo "- 记忆文件由系统自动创建" >> "$MEMORY_FILE"
    echo "- 请在当天的对话和工作中更新此文件" >> "$MEMORY_FILE"
    
    # 此时记忆文件内容较少，不需要立即生成日志文件
    echo "记忆文件已创建，但内容较少，将在晚上9点统一生成日志文件"
    exit 0
fi

# 检查记忆文件内容是否足够
MEM_SIZE=$(wc -c < "$MEMORY_FILE")
if [ $MEM_SIZE -lt 200 ]; then
    echo "记忆文件内容较少，将在晚上9点统一生成日志文件"
    exit 0
fi

# 生成当天的日志文件
LOG_CONTENT=$(cat "$MEMORY_FILE")

# 创建日志文件
LOG_FILE="$LOG_DIR/daily_log_$DATE.txt"
TIMESTAMP=$(date +"%Y-%m-%d %H:%M:%S")
echo "生成时间: $TIMESTAMP" > "$LOG_FILE"
echo "----------------------------------------" >> "$LOG_FILE"
echo "$LOG_CONTENT" >> "$LOG_FILE"

echo "当天日志文件已生成: $LOG_FILE"

# 将日志文件同步到GitHub仓库
cd "$WORKSPACE" || exit
git add "$LOG_FILE" "$MEMORY_FILE"
git commit -m "Add daily log for $DATE" 2>/dev/null || echo "没有需要提交的更改"
git push origin master 2>/dev/null || echo "无法同步到GitHub（可能需要配置SSH密钥）"

echo "日志文件处理完成"
