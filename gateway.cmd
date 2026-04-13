@echo off
rem OpenClaw Gateway (v2026.4.1)
set "TAVILY_API_KEY=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"
set "TMPDIR=C:\Users\mille\AppData\Local\Temp"
set "OPENCLAW_GATEWAY_PORT=18789"
set "OPENCLAW_SYSTEMD_UNIT=openclaw-gateway.service"
set "OPENCLAW_WINDOWS_TASK_NAME=OpenClaw Gateway"
set "OPENCLAW_SERVICE_MARKER=openclaw"
set "OPENCLAW_SERVICE_KIND=gateway"
set "OPENCLAW_SERVICE_VERSION=2026.4.1"
"C:\Program Files\nodejs\node.exe" C:\Users\mille\AppData\Roaming\npm\node_modules\openclaw\dist\index.js gateway --port 18789
