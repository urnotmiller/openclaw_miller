---
name: multi-bot-communication
description: 主机器人与子机器人通信，支持问题识别、分发给对应的子机器人。
---

# 多机器人通信技能

## 适用场景

当需要根据问题类型识别并分发给对应的子机器人时，使用该技能。

## 使用步骤

1. 确保子机器人已启动
2. 使用 `process_message` 函数处理用户消息
3. 系统会自动识别问题类型并分发给对应的子机器人
4. 处理结果会返回给主机器人

## 机器人配置

```python
ROBOT_CONFIGS = {
    "code-expert": {
        "name": "代码专家",
        "url": "http://localhost:5001/api/chat",
        "description": "处理代码相关问题",
        "keywords": ["代码", "编程", "开发", "算法", "调试", "Python", "JavaScript", "Go"],
    },
    "finance-expert": {
        "name": "理财专家",
        "url": "http://localhost:5002/api/chat",
        "description": "处理理财和股票相关问题",
        "keywords": ["股票", "理财", "投资", "基金", "分析", "市场", "上证指数", "创业板"],
    }
}
```

## 主要功能

### 1. 消息处理

```python
from multi_bot_communication.communication import process_message

response = process_message("帮我写一个Python函数计算阶乘")
print(response)
```

### 2. 健康检查

```python
from multi_bot_communication.communication import generate_health_check_response

health_status = generate_health_check_response()
print(health_status)
```

### 3. 问题识别

```python
from multi_bot_communication.communication import recognize_question_type

question_type = recognize_question_type("帮我写代码")
print(question_type)  # 输出: code-expert
```

## 支持的关键词识别

### 代码专家

- 关键词: "代码", "编程", "开发", "算法", "调试", "Python", "JavaScript", "Go"
- 处理问题: 编程相关问题

### 理财专家

- 关键词: "股票", "理财", "投资", "基金", "分析", "市场", "上证指数", "创业板"
- 处理问题: 理财和股票相关问题

## 输出格式

### 成功响应

```
我已经让代码专家来帮你回答这个问题:

def factorial(n):
    if n == 0 or n == 1:
        return 1
    return n * factorial(n-1)
```

### 失败响应

```
很抱歉，我在处理请求时遇到了一些问题:

连接失败
```

## 测试

使用 `test_communication.py` 文件运行测试:

```bash
cd /root/.openclaw/workspace/skills/multi-bot-communication
python test_communication.py
```

## 常见问题解决

### 子机器人未响应

1. 检查子机器人是否已启动
2. 检查网络连接是否正常
3. 查看子机器人日志

### 问题类型未识别

1. 确保问题包含关键词
2. 可以考虑添加更多关键词
3. 检查关键词是否已正确配置
