# agent-browser-clawdbot — 高性能浏览器自动化技能

基于 Vercel Labs `agent-browser` CLI，提供确定性元素选择的浏览器自动化，适合多步骤工作流、复杂SPA操作。

## 安装说明
- 需要先全局安装 CLI: `npm install -g agent-browser`
- 然后安装 Chrome: `agent-browser install`

## 能力
提供比内置canvas更适合AI多步骤交互的浏览器自动化：
- 基于可访问树快照，带ref引用，确定性选择元素
- 支持会话隔离，多用户/多账户隔离
- 支持状态持久化（保存cookie/localStorage）
- 网络控制（拦截/mock请求）
- 支持截图、PDF导出

## 核心命令
```
# 导航
agent-browser open <url>
agent-browser back|forward|reload|close

# 快照（获取元素ref）
agent-browser snapshot -i --json

# 交互（使用@ref）
agent-browser click @e2
agent-browser fill @e3 "text"
agent-browser type @e3 "text"

# 获取信息
agent-browser get text @e1 --json
agent-browser get url --json

# 等待
agent-browser wait --load networkidle
agent-browser wait --text "Welcome"

# 会话
agent-browser --session admin open site.com
agent-browser session list

# 状态
agent-browser state save auth.json
agent-browser state load auth.json

# 截图
agent-browser screenshot --full page.png
```

## 作者
- 技能：Yossi Elkrief (@MaTriXy)
- CLI：Vercel Labs

## 链接
- https://clawhub.ai/matrixy/agent-browser-clawdbot
- https://github.com/vercel/agent-browser
