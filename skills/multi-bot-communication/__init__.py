"""
多机器人通信模块
"""

__version__ = "1.0.0"
__author__ = "OpenClaw"
__description__ = "主机器人与子机器人通信模块，支持问题识别、分发给对应的子机器人"

from .communication import (
    send_message,
    recognize_question_type,
    process_message,
    health_check,
    generate_health_check_response,
    handle_by_main_bot
)

__all__ = [
    "send_message",
    "recognize_question_type", 
    "process_message",
    "health_check",
    "generate_health_check_response",
    "handle_by_main_bot"
]
