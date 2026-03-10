# 金融专家Bot项目

## 项目概述

金融专家Bot是一个基于人工智能的股票分析系统，专门针对中国股票市场（A股、B股、港股）提供专业的金融分析服务。

## 功能特性

### 股票数据获取
- 实时股价查询
- 历史数据获取
- 技术指标计算

### 股票分析
- 股票搜索功能
- 股票推荐
- 买入/卖出信号分析
- 市场趋势预测

### 数据分析
- 技术指标分析（MA20、MACD、RSI）
- 股票基础信息查询
- 公司财务数据获取

### 市场资讯
- 实时市场行情
- 财经新闻
- 政策解读

## 技术架构

### 技术栈
- **编程语言**: Go
- **数据获取**: Tushare API、腾讯财经API
- **数据存储**: 待实现（计划使用MySQL/PostgreSQL）
- **缓存**: Redis（待实现）
- **部署**: Docker（待实现）

### 项目结构
```
workspace/
├── finance_agent/          # 金融专家Bot核心代码
├── study/                  # 学习资料
├── work/                   # 工作文档和方案
├── skills/                 # 技能和功能模块
├── memory/                 # 记忆文件
├── logs/                   # 日志文件
└── README.md              # 项目说明
```

## 安装和运行

### 安装依赖
```bash
# 进入项目目录
cd /root/.openclaw/workspace/finance_agent

# 安装Go依赖（如果有）
go mod tidy
```

### 运行项目
```bash
# 进入金融专家Bot目录
cd /root/.openclaw/workspace/finance_agent

# 运行主程序
go run agent.go
```

### 测试运行
```bash
# 运行测试
go test -v
```

## 使用方法

### 基础功能
```bash
# 获取实时股价
curl "http://localhost:8080/api/quote?code=600000&market=A股"

# 搜索股票
curl "http://localhost:8080/api/search?keyword=茅台"

# 获取股票推荐
curl "http://localhost:8080/api/recommend"
```

### 高级功能
```bash
# 技术指标分析
curl "http://localhost:8080/api/analysis?code=600519&market=A股&period=30"

# 市场趋势预测
curl "http://localhost:8080/api/trend?period=30"
```

## 配置说明

### Tushare API配置
在`/root/.openclaw/workspace/finance_agent/config.yaml`中配置Tushare API密钥：

```yaml
tushare:
  token: "your_tushare_api_token"
```

### 腾讯财经API配置
在`/root/.openclaw/workspace/finance_agent/config.yaml`中配置腾讯财经API：

```yaml
qq_finance:
  base_url: "https://qt.gtimg.cn"
```

## 开发说明

### 开发环境要求
- Go 1.16+
- Git 2.0+
- VS Code或其他Go语言编辑器

### 开发流程
1. 创建功能分支
2. 实现功能
3. 提交代码
4. 创建Pull Request
5. 代码审查
6. 合并到开发分支

## 贡献指南

请参考[CONTRIBUTING.md](CONTRIBUTING.md)文件。

## 许可证

本项目采用MIT许可证，详细信息请查看[LICENSE](LICENSE)文件。

## 联系方式

如有问题或建议，请通过以下方式联系：
- 邮箱：openclaw@example.com
- 项目问题跟踪：https://github.com/openclaw/finance-agent/issues
