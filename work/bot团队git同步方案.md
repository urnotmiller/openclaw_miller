# Bot团队git同步方案

## 概述

为了确保我们的bot团队能够高效、规范地进行代码管理和版本控制，我们制定了以下git同步方案。

## 项目结构优化

### 1. 目录结构调整
```
workspace/
├── finance_agent/          # 金融专家bot代码
├── study/                  # 学习资料
├── work/                   # 工作文档
├── skills/                 # 技能目录
├── memory/                 # 记忆文件
├── logs/                   # 日志
├── .gitignore              # 忽略文件配置
├── README.md               # 项目说明
└── CONTRIBUTING.md         # 贡献指南
```

### 2. .gitignore配置
```gitignore
# 通用忽略模式
*.swp
*.swo
*~
.DS_Store
Thumbs.db

# 日志和临时文件
*.log
*.tmp
temp/
cache/

# 依赖和编译文件
node_modules/
vendor/
dist/
build/
*.exe
*.dll
*.so
*.dylib

# 测试文件
*.test
*.spec
coverage/
.nyc_output/

# 配置文件
.env
.env.local
.env.*.local

# 数据库文件
*.db
*.sqlite
*.sqlite3

# 操作系统临时文件
*.tmp
*.temp

# 记忆和会话文件
memory/*.md
logs/*.txt
```

## 分支策略

### 1. 主分支
```
master          # 生产环境分支（稳定版本）
develop         # 开发分支（最新功能）
```

### 2. 功能分支
```
feature/功能名称    # 新功能开发分支
fix/问题描述       # 错误修复分支
refactor/重构内容  # 代码重构分支
```

### 3. 分支创建规范
```bash
# 创建功能分支
git checkout -b feature/stock-analysis develop

# 创建错误修复分支
git checkout -b fix/data-processing-bug master

# 创建重构分支
git checkout -b refactor/optimize-api develop
```

## 代码提交规范

### 1. 提交消息格式
```
<类型>: <主题>

<详细说明>

[可选的引用信息]
```

### 2. 提交类型说明
- **feat**：新功能开发
- **fix**：错误修复
- **refactor**：代码重构
- **docs**：文档更新
- **style**：代码格式调整（不改变功能）
- **test**：测试相关
- **chore**：构建过程或辅助工具的变更

### 3. 提交示例
```
feat: 添加股票搜索功能

- 实现股票代码和名称搜索
- 添加市场分类筛选
- 优化搜索结果展示

refs: #123
```

### 4. 提交频率
- 每完成一个小功能或修复一个bug就提交一次
- 保持提交内容简洁、明确
- 确保每次提交只包含一个逻辑更改

## 同步流程

### 1. 每日同步流程
```bash
# 早晨工作前
git checkout develop
git pull origin develop

# 开始工作前
git checkout -b feature/your-task develop

# 下午工作结束
git add .
git commit -m "feat: 完成XX功能"
git push origin feature/your-task

# 每日站立会议前
git checkout develop
git pull origin develop
```

### 2. 代码审查和合并
```bash
# 创建Pull Request（PR）
# 等待至少一个团队成员审查
# 审查通过后合并到develop分支

# 每周一次主分支合并
git checkout master
git pull origin master
git merge --no-ff develop
git push origin master
```

### 3. 紧急修复流程
```bash
# 创建紧急修复分支
git checkout -b hotfix/critical-bug master

# 修复问题
git commit -m "fix: 修复XX问题"

# 合并到master和develop
git checkout master
git merge hotfix/critical-bug
git push origin master

git checkout develop
git merge hotfix/critical-bug
git push origin develop

# 删除修复分支
git branch -d hotfix/critical-bug
```

## 团队协作规范

### 1. 代码审查要求
- 所有PR必须经过至少一个团队成员的审查
- 审查时要关注代码质量、功能完整性、性能优化
- 审查通过后才能合并到主分支

### 2. 冲突处理
- 遇到代码冲突时，先与相关人员沟通
- 使用合理的冲突解决方案
- 确保冲突解决后代码能够正常运行

### 3. 版本控制最佳实践
- 定期更新分支
- 使用标签管理重要版本
- 定期备份代码
- 使用.gitattributes文件管理二进制文件

## 自动化工具

### 1. 代码检查和格式化
```json
// .eslintrc.json（示例）
{
  "extends": ["eslint:recommended"],
  "env": {
    "node": true,
    "es2021": true
  },
  "parserOptions": {
    "ecmaVersion": "latest",
    "sourceType": "module"
  },
  "rules": {
    "semi": ["error", "always"],
    "quotes": ["error", "single"]
  }
}
```

### 2. Git钩子配置
```javascript
// .husky/pre-commit
#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

npm run lint
npm run test
```

### 3. CI/CD配置
```yaml
# .github/workflows/main.yml（GitHub Actions）
name: CI/CD

on:
  push:
    branches: [develop, master]
  pull_request:
    branches: [develop]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Use Node.js
        uses: actions/setup-node@v2
        with:
          node-version: '14.x'
      - name: Install dependencies
        run: npm ci
      - name: Run linting
        run: npm run lint
      - name: Run tests
        run: npm run test:coverage
```

## 文档管理

### 1. 项目说明文档
```markdown
# README.md

## 项目概述
金融专家bot是一个基于人工智能的股票分析系统。

## 功能特性
- 股票数据获取
- 技术指标分析
- 股票推荐
- 市场趋势预测

## 技术栈
- Go语言
- Tushare API
- 腾讯财经API

## 安装和运行
```

### 2. 贡献指南
```markdown
# CONTRIBUTING.md

## 开发流程
1. 创建功能分支
2. 提交代码
3. 创建PR
4. 代码审查
5. 合并到开发分支

## 提交规范
参考代码提交规范章节
```

### 3. 变更日志
```markdown
# CHANGELOG.md

## [1.0.0] - 2026-03-10
### 新增功能
- 完成Tushare数据获取对接
- 添加股票推荐功能
- 优化腾讯财经API
```

## 风险控制

### 1. 代码质量保证
- 使用自动化测试
- 定期进行代码审查
- 使用代码质量工具

### 2. 数据安全
- 敏感数据加密存储
- API密钥安全管理
- 定期备份数据

### 3. 灾难恢复
- 定期备份代码和数据
- 使用版本控制系统进行恢复
- 制定灾难恢复计划

## 总结

本git同步方案旨在确保我们的bot团队能够高效、规范地进行代码管理和版本控制。通过合理的分支策略、提交规范和协作流程，我们将能够：

1. 提高开发效率
2. 保证代码质量
3. 降低风险
4. 促进团队协作
5. 确保项目能够按时完成并达到高质量标准

所有团队成员都应遵守这些规范，以确保项目的成功实施。
