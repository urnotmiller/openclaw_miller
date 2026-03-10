# Git常用命令 - 版本管理指南

---

## 常用Git提交版本命令

### 1. 检查状态和变更
```bash
# 查看当前仓库状态
git status

# 查看文件变更详情
git diff

# 查看提交历史
git log --oneline  # 简洁格式
git log --stat     # 显示文件变更统计
```

### 2. 添加和提交
```bash
# 添加单个文件
git add filename.txt

# 添加所有变更（包括新增、修改、删除）
git add .

# 提交变更（带消息）
git commit -m "提交说明"

# 提交变更（详细模式）
git commit -v  # 显示diff内容

# 只提交修改过的文件（不添加新文件）
git commit -a -m "提交说明"
```

### 3. 分支管理
```bash
# 创建并切换到新分支
git checkout -b new-branch

# 切换到现有分支
git checkout branch-name

# 查看所有分支
git branch

# 删除分支
git branch -d branch-name  # 合并过的
git branch -D branch-name  # 未合并的（强制删除）
```

### 4. 远程操作
```bash
# 推送当前分支到远程仓库
git push origin branch-name

# 拉取远程仓库最新代码
git pull

# 查看远程仓库信息
git remote -v

# 添加新的远程仓库
git remote add origin git@github.com:yourname/repo.git
```

### 5. 标签管理
```bash
# 创建轻量标签
git tag v1.0

# 创建带注释的标签
git tag -a v1.0 -m "版本1.0"

# 查看标签
git tag

# 查看标签详情
git show v1.0

# 推送标签到远程仓库
git push origin --tags
```

### 6. 版本回滚
```bash
# 回滚到之前的提交（不保留历史）
git reset --hard commit-id

# 回滚到之前的提交（保留历史）
git revert commit-id

# 查看提交历史（获取commit-id）
git log --oneline
```

---

## 工作流程建议

### 标准工作流程
```bash
# 1. 查看当前状态
git status

# 2. 查看变更
git diff

# 3. 添加变更
git add .

# 4. 提交
git commit -m "提交说明"

# 5. 推送到远程
git push origin master
```

### 提交说明建议格式
```
类型: 简短描述（不超过50个字符）

详细描述（可选）

变更：
- 新增功能X
- 修复了Y问题
- 优化了Z性能

关联：#123（关联issue号）
```

## 常用类型：
- **feat**：新功能
- **fix**：修复bug
- **docs**：文档变更
- **style**：格式调整（不影响功能）
- **refactor**：代码重构
- **perf**：性能优化
- **test**：测试相关
- **chore**：构建过程或辅助工具的变动

---

## 其他实用命令

### 1. 配置
```bash
# 全局配置
git config --global user.name "Your Name"
git config --global user.email "you@example.com"

# 本地仓库配置
git config user.name "Your Name"
git config user.email "you@example.com"

# 查看配置
git config --list
```

### 2. 缓存区管理
```bash
# 查看缓存区内容
git ls-files

# 取消缓存
git reset HEAD -- filename.txt

# 取消所有缓存
git reset HEAD -- .
```

### 3. 历史记录
```bash
# 显示指定文件的修改历史
git log -- filename.txt

# 显示提交树
git log --graph --oneline
```

---

## 注意事项

1. **定期提交**：建议每完成一个小功能或修复一个bug就提交一次
2. **清晰的提交说明**：提交说明要简洁明了，能清楚表达变更的目的
3. **分支管理**：合理使用分支，每个功能或修复一个分支
4. **远程备份**：定期推送到远程仓库，防止数据丢失
5. **版本控制**：重要节点使用标签标记，方便回滚和管理

---

**最后更新：** 2026年3月10日  
**文档作者：** OpenClaw  
**使用说明：** 本文件为Git版本管理参考文档，可根据实际需求修改和补充。
