import re
import requests
import json
import time
import uuid
from typing import Dict, Any, Optional

class ResponseMessage:
    def __init__(self, data):
        self.id = data.get("id", str(uuid.uuid4()))
        self.timestamp = data.get("timestamp", int(time.time()))
        self.status = data.get("status", "success")
        self.content = data.get("content", "")
        self.context = data.get("context", {})
    
    def is_success(self):
        return self.status == "success"

class RequestMessage:
    def __init__(self, data):
        self.id = data.get("id", str(uuid.uuid4()))
        self.content = data.get("content", "")
        self.timestamp = data.get("timestamp", int(time.time()))
        self.sender_id = data.get("sender_id", "user")
        self.session_id = data.get("session_id", str(uuid.uuid4()))
        self.type = data.get("type", "question")

class ROBOT_CONFIGS:
    @staticmethod
    def get():
        return {
            "【代码专家】": {
                "name": "代码专家",
                "url": "http://localhost:5001/api/chat",
                "description": "处理编程相关的问题"
            },
            "【理财专家】": {
                "name": "理财专家",
                "url": "http://localhost:5002/api/chat",
                "description": "处理理财相关的问题"
            },
            "【主机器人】": {
                "name": "主机器人",
                "url": "http://localhost:18789/api/chat",
                "description": "处理一般问题"
            },
            "【阿智】": {
                "name": "主机器人",
                "url": "http://localhost:18789/api/chat",
                "description": "处理一般问题"
            }
        }

def parse_command_format(message):
    # 支持的简单格式识别
    patterns = [
        (r"^:(code|代码|编程)\s*", "【代码专家】"),
        (r"^:(finance|理财|股票|投资)\s*", "【理财专家】"),
        (r"^:(main|阿智|主)\s*", "【主机器人】"),
        (r"^\[代码\]\s*", "【代码专家】"),
        (r"^\[理财\]\s*", "【理财专家】"),
        (r"^\[主\]\s*", "【主机器人】"),
        (r"^\[阿智\]\s*", "【主机器人】"),
        (r"^//(code|代码)\s*", "【代码专家】"),
        (r"^//(finance|理财)\s*", "【理财专家】"),
        (r"^//(main|阿智)\s*", "【主机器人】"),
        (r"^c\s*", "【代码专家】"),        # 代码专家简化指令
        (r"^f\s*", "【理财专家】"),        # 理财专家简化指令
        (r"^m\s*", "【主机器人】"),        # 主机器人简化指令
        (r"^阿智\s*", "【主机器人】"),     # 支持直接使用名字
        (r"^主\s*", "【主机器人】"),       # 支持主机器人简化
        (r"^code\s*", "【代码专家】"),
        (r"^finance\s*", "【理财专家】"),
        (r"^main\s*", "【主机器人】")
    ]
    
    for pattern, robot_identifier in patterns:
        match = re.match(pattern, message, re.IGNORECASE)
        if match:
            # 提取剩余部分作为内容
            content = message[match.end():].strip()
            return robot_identifier, content
    
    # 检查是否包含特殊字符
    if message.startswith("【") and "】" in message:
        end_pos = message.find("】")
        robot_identifier = message[0:end_pos+1]
        content = message[end_pos+1:].strip()
        if robot_identifier in ROBOT_CONFIGS.get():
            return robot_identifier, content
        else:
            return None, message
    
    # 默认情况下，检查关键词来自动识别
    if any(keyword in message.lower() for keyword in ["代码", "编程", "写代码", "python", "java", "c++", "算法"]):
        return "【代码专家】", message
    elif any(keyword in message.lower() for keyword in ["理财", "股票", "投资", "市场"]):
        return "【理财专家】", message
    
    return None, message

def send_message(target: str, content: str, context: Dict[str, Any] = None) -> ResponseMessage:
    robot_config = ROBOT_CONFIGS.get().get(target)
    if not robot_config:
        return ResponseMessage({
            "status": "error",
            "content": "未找到目标机器人"
        })
    
    data = {
        "id": str(uuid.uuid4()),
        "content": content,
        "timestamp": int(time.time()),
        "type": "question"
    }
    
    if context:
        data.update(context)
    
    try:
        response = requests.post(
            robot_config["url"],
            json=data,
            timeout=30
        )
        
        if response.status_code == 200:
            return ResponseMessage(response.json())
        else:
            return ResponseMessage({
                "status": "error",
                "content": f"HTTP Error: {response.status_code}"
            })
            
    except Exception as e:
        return ResponseMessage({
            "status": "error",
            "content": f"通信错误: {str(e)}"
        })

def recognize_command(message):
    robot_identifier, content = parse_command_format(message)
    return robot_identifier, content

def process_message(message, context: Dict[str, Any] = None) -> str:
    robot_identifier, content = parse_command_format(message)
    
    if robot_identifier:
        # 对于主机器人，直接处理而不发送到API
        if robot_identifier == "【主机器人】" or robot_identifier == "【阿智】":
            return handle_by_main_bot(content, context)
            
        # 检查机器人是否在线
        check_response = check_robot_status(robot_identifier)
        
        if not check_response.is_success():
            return "抱歉，该机器人暂时不可用。请稍后重试。"
            
        # 发送消息给子机器人
        response = send_message(robot_identifier, content, context)
        
        if response.is_success():
            robot_config = ROBOT_CONFIGS.get().get(robot_identifier)
            robot_name = robot_config["name"] if robot_config else robot_identifier
            
            return response.content
        else:
            return "抱歉，在处理请求时发生了错误。请稍后重试。"
    else:
        # 主机器人直接处理
        return handle_by_main_bot(content, context)

def check_robot_status(robot_identifier: str) -> ResponseMessage:
    config = ROBOT_CONFIGS.get().get(robot_identifier)
    if not config:
        return ResponseMessage({
            "status": "error",
            "content": "未找到目标机器人"
        })
    
    health_check_url = config["url"].replace("/api/chat", "/health")
    
    try:
        response = requests.get(health_check_url, timeout=5)
        
        if response.status_code == 200:
            data = response.json()
            return ResponseMessage({
                "status": "success",
                "content": f"{config['name']} 状态良好"
            })
        else:
            return ResponseMessage({
                "status": "error",
                "content": f"{config['name']} 状态异常"
            })
            
    except Exception as e:
        # 对于主机器人，如果健康检查失败，我们仍然认为它在线
        # 因为主机器人可能没有独立的API接口
        if robot_identifier == "【主机器人】" or robot_identifier == "【阿智】":
            return ResponseMessage({
                "status": "success",
                "content": "主机器人正在运行"
            })
            
        return ResponseMessage({
            "status": "error",
            "content": f"{config['name']} 无法访问"
        })

def handle_by_main_bot(content: str, context: Dict[str, Any] = None) -> str:
    responses = {
        "你好": "你好！我是主机器人。我可以帮你解答各种问题。你需要什么帮助？",
        "健康检查": get_health_check_summary(),
        "状态": get_health_check_summary(),
        "介绍": "我是一个智能多机器人系统，由三个机器人组成：\n\n- 代码专家：负责编程相关问题\n- 理财专家：负责理财相关问题\n- 主机器人：负责一般问题和协调\n\n你可以使用指令前缀（如:code或:finance）来呼叫相应的机器人。",
        "使用帮助": "我是一个智能多机器人系统，由三个机器人组成：\n\n- 代码专家：负责编程相关问题\n- 理财专家：负责理财相关问题\n- 主机器人：负责一般问题和协调\n\n你可以使用指令前缀（如:code或:finance）来呼叫相应的机器人。"
    }
    
    lower_content = content.lower().strip()
    
    for keyword, response in responses.items():
        if keyword in lower_content:
            return response
            
    return "抱歉，我暂时无法回答这个问题。你可以尝试使用其他格式或呼叫其他机器人。"

def get_health_check_summary():
    result = []
    
    configs = ROBOT_CONFIGS.get()
    
    for robot_identifier, config in configs.items():
        # 跳过重复的配置（主机器人和阿智指向同一配置）
        if config["name"] in [item["机器人"] for item in result]:
            continue
            
        status_response = check_robot_status(robot_identifier)
        
        result.append({
            "机器人": config["name"],
            "状态": "✅ 在线" if status_response.is_success() else "❌ 离线",
            "信息": status_response.content
        })
        
    summary = []
    for item in result:
        summary.append(f"• {item['机器人']}: {item['状态']}")
        
    return "当前系统状态：\n" + "\n".join(summary)

def analyze_stock(stock_code, start_date, end_date):
    try:
        data = {
            "id": str(uuid.uuid4()),
            "content": f"分析股票: {stock_code} 从 {start_date} 到 {end_date}",
            "timestamp": int(time.time()),
            "type": "analysis",
            "params": {
                "stock_code": stock_code,
                "start_date": start_date,
                "end_date": end_date
            }
        }
        
        response = send_message("【理财专家】", data["content"], data)
        
        if response.is_success():
            return response.content
        else:
            return f"分析失败: {response.content}"
            
    except Exception as e:
        return f"分析过程中出错: {str(e)}"

def get_robot_info(robot_identifier):
    config = ROBOT_CONFIGS.get().get(robot_identifier)
    if config:
        return {
            "name": config["name"],
            "url": config["url"],
            "description": config["description"]
        }
    return None

def send_health_check():
    return get_health_check_summary()

if __name__ == "__main__":
    # 测试函数
    print("=== 多机器人通信系统测试 ===\n")
    
    # 测试指令识别
    test_cases = [
        ":code 写一个阶乘函数",
        "【代码专家】写一个阶乘函数",
        "[代码] 写一个阶乘函数",
        "//code 写一个阶乘函数",
        "c 写一个阶乘函数",
        ":finance 分析股票市场",
        "【理财专家】分析股票市场",
        "[理财] 分析股票市场",
        "//finance 分析股票市场",
        "f 分析股票市场",
        ":main 你好",
        "【主机器人】你好",
        "[主] 你好",
        "//main 你好",
        "m 你好",
        "阿智 你好"
    ]
    
    print("测试指令识别:")
    print("-" * 40)
    for case in test_cases:
        try:
            robot, content = parse_command_format(case)
            if robot:
                print(f"✓ {case:<40} → {robot:<8} | {content}")
            else:
                print(f"✗ {case:<40} → 未识别")
        except Exception as e:
            print(f"✗ {case:<40} → 错误: {e}")
    
    print("\n" + "-" * 40)
    
    # 健康检查
    print("\n测试健康检查:")
    health_response = send_health_check()
    print(health_response)
