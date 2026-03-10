import sys
import os

# 添加当前目录到路径
sys.path.insert(0, os.path.abspath(os.path.dirname(__file__)))

import communication as mbcomm

def test_health_check():
    """测试健康检查功能"""
    print("1. 测试健康检查功能:")
    results = mbcomm.health_check()
    
    assert 'main' in results
    assert 'code-expert' in results
    assert 'finance-expert' in results
    assert results['main']['status'] == 'healthy'
    
    print(f"✅ 主机器人状态: {results['main']['status']}")
    print(f"✅ 代码专家状态: {results['code-expert']['status']}")
    print(f"✅ 理财专家状态: {results['finance-expert']['status']}")
    
    return results

def test_message_recognition():
    """测试消息识别功能"""
    print("\n2. 测试消息识别功能:")
    
    test_cases = [
        ("帮我写一个Python函数", "code-expert"),
        ("分析一下600000.SH的股票走势", "finance-expert"),
        ("你好", None),
        ("我是miller", None)
    ]
    
    for question, expected_type in test_cases:
        actual_type = mbcomm.recognize_question_type(question)
        
        result_msg = f"{'✅' if actual_type == expected_type else '❌'} '{question}'"
        
        if actual_type:
            print(f"{result_msg} -> 识别为 '{actual_type}'")
        else:
            print(f"{result_msg} -> 主机器人处理")
            
    return test_cases

def test_send_message():
    """测试发送消息功能"""
    print("\n3. 测试发送消息功能:")
    
    test_messages = {
        "code-expert": "写一个简单的Python函数计算阶乘",
        "finance-expert": "分析上海证券交易所股票走势"
    }
    
    results = {}
    for robot, message in test_messages.items():
        try:
            response = mbcomm.send_message(robot, message)
            print(f"✅ {robot} 响应: {response.content[:100]}")
            results[robot] = response
        except Exception as e:
            print(f"❌ {robot} 发送失败: {e}")
            print("   可能子机器人尚未启动或配置不正确")
            
    return results

def test_message_processing():
    """测试消息处理流程"""
    print("\n4. 测试消息处理流程:")
    
    test_questions = [
        "你好，我是miller",
        "帮我写一个Python程序计算斐波那契数列",
        "分析一下上证指数最近一周的走势"
    ]
    
    for question in test_questions:
        try:
            result = mbcomm.process_message(question)
            print(f"✅ '{question}' -> '{result[:100]}'")
        except Exception as e:
            print(f"❌ 处理失败: {e}")

def run_all_tests():
    """运行所有测试"""
    print("=== 多机器人通信系统测试 ===")
    
    try:
        health_results = test_health_check()
    except Exception as e:
        print(f"❌ 健康检查失败: {e}")
        return False
        
    test_message_recognition()
    
    try:
        send_results = test_send_message()
    except Exception as e:
        print(f"❌ 发送消息测试失败: {e}")
        
    try:
        test_message_processing()
    except Exception as e:
        print(f"❌ 消息处理测试失败: {e}")
        
    print("\n=== 测试完成 ===")
    
    # 检查是否有健康问题
    if health_results.get('main', {}).get('status') != 'healthy':
        print("❌ 主机器人未正常运行")
        return False
        
    # 检查子机器人
    if health_results.get('code-expert', {}).get('status') != 'healthy':
        print("⚠️  代码专家未正常运行（可能未启动）")
        
    if health_results.get('finance-expert', {}).get('status') != 'healthy':
        print("⚠️  理财专家未正常运行（可能未启动）")
        
    return True

def health_check_to_string():
    """将健康检查结果转换为字符串"""
    try:
        return mbcomm.generate_health_check_response()
    except Exception as e:
        return f"健康检查失败: {e}"

if __name__ == "__main__":
    print("=== 多机器人通信系统测试 ===")
    print(f"测试时间: {mbcomm.RequestMessage('test', 'test', 'test').to_dict()['timestamp']}")
    print()
    
    if run_all_tests():
        print("\n✅ 所有测试通过")
        print(health_check_to_string())
    else:
        print("\n❌ 测试失败")
