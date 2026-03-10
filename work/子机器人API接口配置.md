# 子机器人API接口配置

## 概述

子机器人现在是直接启动的OpenClaw实例，而不是我们专门编写的API服务器。因此，它们的API接口可能不符合我们预期的格式。我们需要了解OpenClaw实例的API接口。

## 检查OpenClaw的API接口

### 查看主机器人API接口

我们可以通过查看OpenClaw源代码中的API接口定义：

```javascript
// 这可能位于 OpenClaw 源代码中
// 检查是否有相关的API路由文件
const fs = require('fs');
const path = require('path');

// 查找 OpenClaw 源代码
const npmPrefix = require('child_process').execSync('npm prefix -g', { encoding: 'utf-8' }).trim();
const openclawPath = path.join(npmPrefix, 'node_modules/openclaw');

console.log(`OpenClaw 安装在: ${openclawPath}`);

// 查找API相关文件
console.log('\n查找API相关文件:');
const apiFiles = [];
const searchPatterns = [
    'api',
    'server',
    'router',
    'endpoint'
];

// 递归查找
function searchFiles(dir) {
    const files = fs.readdirSync(dir);
    files.forEach(file => {
        const filePath = path.join(dir, file);
        const stat = fs.statSync(filePath);
        
        if (stat.isDirectory()) {
            searchFiles(filePath);
        } else {
            if (searchPatterns.some(pattern => 
                filePath.toLowerCase().includes(pattern) && 
                (file.endsWith('.js') || file.endsWith('.ts')))) {
                apiFiles.push(filePath);
            }
        }
    });
}

try {
    searchFiles(path.join(openclawPath, 'src'));
} catch (e) {
    console.log('无法访问 src 目录，尝试其他位置');
    try {
        searchFiles(path.join(openclawPath, 'dist'));
    } catch (e2) {
        console.log('无法访问 dist 目录');
    }
}

console.log(apiFiles);
```

## 可能的解决方案

### 方案1：使用OpenClaw的内置通信方式

OpenClaw可能有内置的通信方式，我们可以通过查看它的文档或源代码找到：

```javascript
// 检查是否有相关的OpenClaw命令
const fs = require('fs');
const openclawConfig = '/root/.openclaw/openclaw.json';

if (fs.existsSync(openclawConfig)) {
    const config = JSON.parse(fs.readFileSync(openclawConfig, 'utf-8'));
    console.log('OpenClaw配置:');
    console.log(config);
}
```

### 方案2：为子机器人创建专门的API服务器

我们可以创建专门的API服务器来模拟子机器人的功能，或者调整现有的子机器人：

```javascript
// 示例子机器人API服务器
const express = require('express');
const app = express();

app.use(express.json());

// 健康检查
app.get('/health', (req, res) => {
    res.json({
        status: 'healthy',
        message: '我是子机器人，正在运行'
    });
});

// 聊天接口
app.post('/api/chat', (req, res) => {
    const message = req.body.content || req.body.message;
    
    if (!message) {
        return res.status(400).json({
            status: 'error',
            message: '缺少消息内容'
        });
    }
    
    // 根据消息类型返回不同的响应
    let response;
    if (req.body.target === 'code-expert') {
        response = handleCodeQuestion(message);
    } else if (req.body.target === 'finance-expert') {
        response = handleFinanceQuestion(message);
    } else {
        response = '我是子机器人，但不知道你想要什么';
    }
    
    res.json({
        id: req.body.id || Date.now(),
        requestId: req.body.id,
        timestamp: Date.now(),
        status: 'success',
        content: response,
        context: {
            source: req.body.target,
            target: 'main',
            executionTime: 0.5
        }
    });
});

function handleCodeQuestion(question) {
    const lowerQuestion = question.toLowerCase();
    
    if (lowerQuestion.includes('阶乘')) {
        return `def factorial(n):
    if n == 0 or n == 1:
        return 1
    return n * factorial(n-1)`;
    }
    
    if (lowerQuestion.includes('斐波那契')) {
        return `def fibonacci(n):
    if n <= 0:
        return 0
    elif n == 1:
        return 1
    return fibonacci(n-1) + fibonacci(n-2)`;
    }
    
    return `我可以帮你编写各种代码，但需要更具体的问题。例如，你可以问：
- "写一个Python函数计算阶乘"
- "写一个JavaScript数组排序函数"
- "帮我写一个简单的API服务器"`;
}

function handleFinanceQuestion(question) {
    const lowerQuestion = question.toLowerCase();
    
    if (lowerQuestion.includes('股票') || lowerQuestion.includes('上证指数')) {
        return `我可以帮你分析股票市场，但需要更具体的问题。例如，你可以问：
- "分析一下600000.SH的股票走势"
- "上证指数最近一周的走势分析"
- "推荐一些蓝筹股"

注意：我提供的是信息分析，不是投资建议。投资有风险，入市需谨慎。`;
    }
    
    if (lowerQuestion.includes('理财')) {
        return `我可以帮你提供理财建议，但需要更具体的问题。例如，你可以问：
- "我有10万，应该如何投资理财"
- "短期投资和长期投资的区别"
- "如何进行资产配置"

注意：我提供的是信息分析，不是投资建议。理财有风险，投资需谨慎。`;
    }
    
    return `我可以帮你分析各种金融问题，但需要更具体的问题。例如，你可以问：
- "股票走势分析"
- "投资理财建议"
- "市场风险评估"`;
}

// 启动服务器
const port = process.argv[2] || 5000;
app.listen(port, () => {
    console.log(`子机器人API服务器运行在端口 ${port}`);
});

// 健康检查
console.log(`子机器人API服务器启动在 http://localhost:${port}`);
```

### 方案3：修改主机器人通信代码以匹配OpenClaw的API

我们需要检查OpenClaw实例的API接口：

```javascript
// 使用HTTP请求尝试获取OpenClaw的API文档或帮助信息
const http = require('http');

function fetchApiInfo(port) {
    return new Promise((resolve, reject) => {
        const req = http.request({
            hostname: 'localhost',
            port: port,
            path: '/api',
            method: 'GET'
        }, (res) => {
            let data = '';
            res.on('data', chunk => { data += chunk; });
            res.on('end', () => resolve(data));
        });
        
        req.on('error', reject);
        req.end();
    });
}

async function checkOpenClawApi() {
    console.log('检查 OpenClaw 主机器人 API (18789):');
    try {
        const data = await fetchApiInfo(18789);
        console.log(data);
    } catch (e) {
        console.log('失败:', e.message);
    }
    
    console.log('\n检查 OpenClaw 代码专家 API (5001):');
    try {
        const data = await fetchApiInfo(5001);
        console.log(data);
    } catch (e) {
        console.log('失败:', e.message);
    }
    
    console.log('\n检查 OpenClaw 理财专家 API (5002):');
    try {
        const data = await fetchApiInfo(5002);
        console.log(data);
    } catch (e) {
        console.log('失败:', e.message);
    }
}

checkOpenClawApi();
```

## 实际方案

最实际的方案是，我们需要修改主机器人的通信代码，以适应OpenClaw实例的API接口，或者直接为子机器人创建简单的API服务器。

### 修改通信代码

我们需要修改`communication.py`文件中的`send_message`函数，使用OpenClaw的实际API接口：

```python
import requests

def send_message(target: str, content: str, context: Dict[str, Any] = None) -> ResponseMessage:
    """
    发送消息给子机器人
    :param target: 目标机器人名称
    :param content: 消息内容
    :param context: 上下文信息
    :return: 响应消息
    """
    if target not in ROBOT_CONFIGS:
        return ResponseMessage({
            "status": "error",
            "content": f"未知的目标机器人: {target}"
        })
    
    robot_config = ROBOT_CONFIGS[target]
    
    # 尝试使用OpenClaw的实际API接口
    # 这需要我们了解OpenClaw的API接口格式
    # 以下是一个示例，可能不准确
    
    try:
        # 尝试获取OpenClaw的API信息
        # 首先获取令牌
        auth_response = requests.post(
            f"{robot_config['url']}/auth/login",
            data={"username": "admin", "password": "password"}
        )
        
        if auth_response.status_code != 200:
            return ResponseMessage({
                "status": "error",
                "content": f"认证失败: {auth_response.status_code}"
            })
            
        auth_data = auth_response.json()
        
        # 然后发送消息
        response = requests.post(
            f"{robot_config['url']}/api/chat",
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {auth_data['token']}"
            },
            json={
                "content": content,
                "context": context,
                "sender": context.get('user', {}).get('name', '用户')
            },
            timeout=30
        )
        
        if response.status_code == 200:
            data = response.json()
            return ResponseMessage({
                "id": data.get("id"),
                "status": "success",
                "content": data.get("content"),
                "context": data.get("context", {})
            })
        
        return ResponseMessage({
            "status": "error",
            "content": f"API请求失败: {response.status_code}"
        })
        
    except Exception as e:
        print(f"Error communicating with {target}: {e}")
        return ResponseMessage({
            "status": "error",
            "content": f"通信失败: {e}"
        })
```

### 使用专门的子机器人API服务器

另一个方案是，我们为代码专家和理财专家创建简单的API服务器：

```javascript
// 代码专家API服务器 (port 5001)
const express = require('express');
const app = express();
app.use(express.json());

app.get('/health', (req, res) => {
    res.json({
        status: 'healthy',
        message: '代码专家API服务器正在运行'
    });
});

app.post('/api/chat', (req, res) => {
    res.json({
        status: 'success',
        content: handleCodeQuestion(req.body.content),
        context: {
            source: 'code-expert',
            target: 'main'
        }
    });
});

function handleCodeQuestion(question) {
    // 这里是简单的代码问题回答逻辑
    if (question.includes('阶乘')) {
        return `def factorial(n):
    if n == 0 or n == 1:
        return 1
    return n * factorial(n-1)`;
    }
    
    return "我是代码专家，但这个问题我暂时无法回答";
}

app.listen(5001, () => {
    console.log('代码专家API服务器运行在端口 5001');
});

// 理财专家API服务器 (port 5002)
const express2 = require('express');
const app2 = express2();
app2.use(express2.json());

app2.get('/health', (req, res) => {
    res.json({
        status: 'healthy',
        message: '理财专家API服务器正在运行'
    });
});

app2.post('/api/chat', (req, res) => {
    res.json({
        status: 'success',
        content: handleFinanceQuestion(req.body.content),
        context: {
            source: 'finance-expert',
            target: 'main'
        }
    });
});

function handleFinanceQuestion(question) {
    // 这里是简单的理财问题回答逻辑
    if (question.includes('股票')) {
        return "我可以帮你分析股票，但我现在无法提供实时数据";
    }
    
    return "我是理财专家，但这个问题我暂时无法回答";
}

app2.listen(5002, () => {
    console.log('理财专家API服务器运行在端口 5002');
});
```

### 运行这些服务器

```bash
# 保存上述代码为 separate-api-server.js
# 安装依赖
npm install express

# 启动服务器
node separate-api-server.js &
```

## 总结

我们有几个选择：

1. **了解OpenClaw的API接口并使用它们**：我们需要查看OpenClaw的源代码或文档，以找出它们的API接口格式。

2. **使用简单的专门API服务器**：我们为子机器人创建简单的API服务器，这些服务器可以处理我们发送的消息。

3. **修改通信代码以处理OpenClaw实例**：我们修改主机器人的通信代码，以更好地与OpenClaw子实例通信。

选择哪种方案取决于我们的需求和OpenClaw的功能。
