# 金融专家Bot项目贡献指南

## 欢迎参与

欢迎您参与金融专家Bot项目的开发！我们非常高兴您能加入我们的团队，共同开发这个优秀的金融分析系统。

## 开发流程

### 1. 代码获取
```bash
# 克隆项目到本地
git clone [项目地址]
cd [项目目录]

# 确保使用开发分支
git checkout develop
```

### 2. 创建功能分支
```bash
# 创建新的功能分支
git checkout -b feature/your-feature-name develop

# 或者创建错误修复分支
git checkout -b fix/bug-description master
```

### 3. 开发和测试
```bash
# 安装依赖
cd finance_agent
go mod tidy

# 运行项目
go run agent.go

# 运行测试
go test -v
```

### 4. 提交代码
```bash
# 查看修改的文件
git status

# 暂存更改
git add .

# 提交代码
git commit -m "feat: 添加股票搜索功能"
```

### 5. 代码审查
```bash
# 推送到远程仓库
git push origin feature/your-feature-name

# 在GitHub上创建Pull Request
# 等待至少一个团队成员审查
# 审查通过后合并到develop分支
```

### 6. 完成开发
```bash
# 合并到开发分支后，删除本地功能分支
git checkout develop
git pull origin develop
git branch -d feature/your-feature-name
```

## 提交规范

### 提交消息格式
```
<类型>: <主题>

<详细说明>

[可选的引用信息]
```

### 提交类型说明
- **feat**：新功能开发
- **fix**：错误修复
- **refactor**：代码重构
- **docs**：文档更新
- **style**：代码格式调整（不改变功能）
- **test**：测试相关
- **chore**：构建过程或辅助工具的变更

### 提交示例
```
feat: 添加股票搜索功能

- 实现股票代码和名称搜索
- 添加市场分类筛选
- 优化搜索结果展示

refs: #123
```

## 代码风格

### Go语言规范
- 遵循Go语言官方规范
- 使用`go fmt`格式化代码
- 变量和函数名使用驼峰命名法
- 文件名使用小写字母，单词间用下划线分隔

### 代码结构
- 每个文件只包含一个主题的代码
- 使用适当的包结构
- 避免过长的函数（建议不超过50行）
- 提供适当的代码注释

### 错误处理
- 对所有可能的错误进行适当处理
- 使用有意义的错误信息
- 不要忽略错误
- 使用defer和recover处理panic

## 测试要求

### 单元测试
- 为每个函数编写单元测试
- 测试覆盖率不低于70%
- 使用表驱动测试
- 保持测试简单和可读性

### 集成测试
- 测试模块间的交互
- 测试API接口
- 测试数据获取和处理流程

### 测试命令
```bash
# 运行所有测试
go test -v

# 只运行特定包的测试
cd finance_agent
go test -v ./...

# 生成测试覆盖率报告
go test -coverprofile coverage.out ./...
go tool cover -html coverage.out
```

## 文档要求

### 代码注释
- 为所有导出的函数、类型和变量添加注释
- 注释应该清晰、简洁
- 使用标准的Go文档格式

### README文档
- 包含项目概述、功能特性、安装方法等
- 使用Markdown格式
- 保持文档更新

### CHANGELOG
- 记录版本变更
- 使用语义化版本号
- 包含新增功能、修复内容和破坏性变更

## 问题反馈

### 发现Bug
1. 首先在GitHub Issues中搜索是否有类似问题
2. 打开一个新的Issue，标题清晰明确
3. 详细描述问题，包括
   - 预期行为
   - 实际行为
   - 复现步骤
   - 系统信息
   - 日志输出（如果有）

### 提出功能建议
1. 在GitHub Issues中搜索是否有类似建议
2. 打开一个新的Issue，标题清晰明确
3. 详细描述建议，包括
   - 功能概述
   - 预期用途
   - 实现思路（如果有）

### 参与讨论
- 积极参与问题和功能建议的讨论
- 提供有价值的反馈
- 尊重其他贡献者的意见

## 社区规范

### 行为准则
- 尊重所有贡献者
- 保持专业和礼貌的沟通方式
- 不攻击个人，只讨论代码和想法
- 帮助新成员
- 接受建设性的批评

### 代码审查标准
- 代码质量
- 功能完整性
- 测试覆盖率
- 代码风格
- 错误处理

### 奖励和认可
- 定期识别和奖励优秀贡献者
- 为重要功能添加贡献者信息
- 在项目文档中感谢贡献者

## 许可证

本项目采用MIT许可证，详细信息请查看[LICENSE](LICENSE)文件。

## 联系方式

如有问题或建议，请通过以下方式联系：
- 邮箱：openclaw@example.com
- 项目问题跟踪：https://github.com/openclaw/finance-agent/issues
