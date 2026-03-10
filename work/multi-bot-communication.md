# 多机器人通信方案

## 架构设计

### 1. 通信协议
采用HTTP API通信协议，每个机器人提供以下接口：

```
GET /api/health           // 健康检查
POST /api/chat            // 聊天接口
POST /api/execute         // 执行接口
GET /api/skills           // 获取技能列表
```

### 2. 数据格式
```json
{
  "id": "request-id",
  "source": "main",
  "target": "code-expert",
  "timestamp": 1710000000000,
  "type": "question",
  "content": "请帮我编写一个Python函数计算阶乘",
  "context": {
    "sessionId": "12345",
    "userId": "miller",
    "userName": "Miller"
  }
}
```

### 3. 响应格式
```json
{
  "id": "response-id",
  "requestId": "request-id",
  "timestamp": 1710000001000,
  "status": "success",
  "type": "answer",
  "content": "def factorial(n):\n    if n == 0 or n == 1:\n        return 1\n    return n * factorial(n-1)",
  "context": {
    "source": "code-expert",
    "target": "main",
    "executionTime": 1.23
  }
}
```

## 实现方式

### 1. 主机器人通信代码

#### 发送消息函数
```python
import requests
import json
import uuid
import time

def send_message(target, content, message_type="question", context={}):
    """
    发送消息给子机器人
    :param target: 目标机器人名称
    :param content: 消息内容
    :param message_type: 消息类型
    :param context: 上下文信息
    """
    endpoints = {
        "code-expert": "http://localhost:5001/api/chat",
        "finance-expert": "http://localhost:5002/api/chat"
    }
    
    if target not in endpoints:
        return {"status": "error", "content": f"未知的目标机器人: {target}"}
    
    url = endpoints[target]
    request_id = str(uuid.uuid4())
    
    payload = {
        "id": request_id,
        "source": "main",
        "target": target,
        "timestamp": int(time.time() * 1000),
        "type": message_type,
        "content": content,
        "context": context
    }
    
    try:
        response = requests.post(
            url, 
            data=json.dumps(payload), 
            headers={"Content-Type": "application/json"},
            timeout=30
        )
        
        if response.status_code == 200:
            result = response.json()
            result["requestId"] = request_id
            return result
        else:
            return {"status": "error", "content": f"HTTP错误: {response.status_code}"}
            
    except Exception as e:
        return {"status": "error", "content": f"请求失败: {str(e)}"}
```

#### 接收和处理消息函数
```python
from flask import Flask, request, jsonify
import json

app = Flask(__name__)

@app.route("/api/chat", methods=["POST"])
def receive_message():
    """
    接收来自子机器人或其他系统的消息
    """
    try:
        data = request.get_json()
        
        # 验证必要字段
        required_fields = ["id", "source", "target", "timestamp", "type", "content"]
        for field in required_fields:
            if field not in data:
                return jsonify({
                    "id": data.get("id", str(uuid.uuid4())),
                    "requestId": data.get("id"),
                    "timestamp": int(time.time() * 1000),
                    "status": "error",
                    "content": f"缺少必要字段: {field}"
                })
                
        # 根据消息类型处理
        if data["type"] == "question":
            result = handle_question(data["source"], data["content"], data.get("context", {}))
        elif data["type"] == "answer":
            result = handle_answer(data["source"], data["content"], data.get("context", {}))
        elif data["type"] == "error":
            result = handle_error(data["source"], data["content"], data.get("context", {}))
        else:
            result = {"status": "error", "content": f"未知的消息类型: {data['type']}"}
            
        return jsonify({
            "id": str(uuid.uuid4()),
            "requestId": data["id"],
            "timestamp": int(time.time() * 1000),
            "status": result.get("status", "success"),
            "content": result.get("content", ""),
            "context": {
                "source": "main",
                "target": data["source"],
                "executionTime": result.get("executionTime")
            }
        })
        
    except Exception as e:
        return jsonify({
            "id": str(uuid.uuid4()),
            "status": "error",
            "content": f"服务器内部错误: {str(e)}"
        }), 500
```

## 部署方案

### 1. 主机器人
- **端口**: 18789
- **配置文件**: /root/.openclaw/openclaw.json
- **服务**: openclaw-gateway

### 2. 子机器人
#### 代码专家
- **端口**: 5001
- **配置文件**: /root/.openclaw-code-bot/openclaw.json
- **服务**: openclaw-code-bot

#### 理财专家
- **端口**: 5002
- **配置文件**: /root/.openclaw-finance-bot/openclaw.json
- **服务**: openclaw-finance-bot

## 启动和管理

### 启动服务
```bash
# 启动主机器人
systemctl --user start openclaw-gateway

# 启动代码专家
systemctl --user start openclaw-code-bot

# 启动理财专家
systemctl --user start openclaw-finance-bot

# 启动所有服务
./start-all-bots.sh
```

### 检查运行状态
```bash
# 检查所有服务
systemctl --user status openclaw-gateway openclaw-code-bot openclaw-finance-bot

# 查看日志
journalctl --user -f -u openclaw-gateway
journalctl --user -f -u openclaw-code-bot
journalctl --user -f -u openclaw-finance-bot
```

### 停止服务
```bash
# 停止所有服务
./stop-all-bots.sh
```

## 安全考虑

### 1. 访问控制
- 所有API接口只允许本地访问
- 使用HTTP基本认证或令牌认证
- 限制请求频率

### 2. 数据加密
- 使用HTTPS协议（生产环境）
- 数据传输加密
- 敏感信息加密存储

### 3. 防火墙配置
```bash
# 只允许本地访问
iptables -A INPUT -p tcp --dport 18789 -s 127.0.0.1 -j ACCEPT
iptables -A INPUT -p tcp --dport 18789 -j DROP

iptables -A INPUT -p tcp --dport 5001 -s 127.0.0.1 -j ACCEPT
iptables -A INPUT -p tcp --dport 5001 -j DROP

iptables -A INPUT -p tcp --dport 5002 -s 127.0.0.1 -j ACCEPT
iptables -A INPUT -p tcp --dport 5002 -j DROP
```

## 测试方案

### 1. 健康检查
```bash
# 主机器人健康检查
curl -s http://localhost:18789/api/health

# 代码专家健康检查
curl -s http://localhost:5001/api/health

# 理财专家健康检查
curl -s http://localhost:5002/api/health
```

### 2. 通信测试
```bash
# 发送消息给代码专家
curl -X POST http://localhost:18789/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-1",
    "source": "test",
    "target": "main",
    "timestamp": 1710000000000,
    "type": "question",
    "content": "请帮我计算阶乘",
    "context": {
      "sessionId": "test-session"
    }
  }'
```

### 3. 集成测试
```bash
# 发送复杂问题给理财专家
curl -X POST http://localhost:18789/api/execute \
  -H "Content-Type: application/json" \
  -d '{
    "id": "execute-1",
    "source": "user",
    "target": "main",
    "timestamp": 1710000000000,
    "command": "analyze_stock",
    "params": {
      "stockCode": "600000.SH",
      "startDate": "2026-01-01",
      "endDate": "2026-03-10"
    }
  }'
```
