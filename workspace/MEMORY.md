# MEMORY.md - 长期记忆索引

> 这是索引层，不是仓库。详细内容和教训去对应文件查找。

---

## 👤 用户基本信息

| 字段 | 内容 |
|------|------|
| 姓名 | miller |
| 称呼 | boss |
| 时区 | Asia/Shanghai (GMT+8) |
| 偏好 | 不客套，直接犀利；有判断有主张；不自作主张越权 |

---

## 📁 文件索引

| 文件 | 用途 |
|------|------|
| `memory/projects.md` | 所有项目当前状态总览 |
| `memory/lessons.md` | 踩坑教训，按严重程度分级 |
| `memory/YYYY-MM-DD.md` | 每日结论日志（不是过程！） |
| `workspace/SOUL.md` | PM 人设（胖达人/panda） |
| `workspace/IDENTITY.md` | 身份定义 |
| `workspace/USER.md` | 用户信息 |

---

## ⚡ 核心规则（不常改动）

### 记忆管理（复盘=双写）
每次复盘时，两边同步更新：
- `memory/projects.md` — 项目节点变化
- `memory/lessons.md` — 教训分级（🔴严重/🟡中等/🟢轻微）
- `.learnings/LEARNINGS.md` — 教训、修正、洞察
- `.learnings/ERRORS.md` — 错误/失败记录
- `.learnings/FEATURE_REQUESTS.md` — 做不到的需求

### Agent 调度
- **FE/BE/QA**：用 `sessions_spawn` 派发，完成后自动返回
- **UI/OP**：用 `sessions_spawn mode=session` + `sessions_send`，支持 boss 多轮审核
- **FE 前置**：UI设计稿 + API接口文档 **双确认后**才能启动
- **QA 二次review**：FE/BE 修复完成后，必须再次调度 QA 确认，PM 不自作判断

### 项目结构
- 路径：`Projects/{项目名-版本}`
- 看板：`docs/kanban/KANBAN.md`
- PRD：`docs/requirements/PRD.md`
- 访谈纪要：`docs/requirements/访谈纪要.md`

### 需求访谈
大项目：PRD 前先做结构化访谈 → 输出 `访谈纪要.md` → 再写 PRD

---

## 🧠 当前项目

| 项目 | 阶段 | 状态 |
|------|------|------|
| weather-website-mvp | A12完成 | 测试标注，流程已理顺 |

---

## 📌 关键决策（按时间）

| 日期 | 决策 |
|------|------|
| 2026-04-08 | 归责原则：谁的活谁干，PM 不代办 |
| 2026-04-09 | UI/OP 用 send+session；FE 前置双确认；QA 二次review 强制 |
| 2026-04-09 | 记忆分层：projects + lessons + 每日结论日志 |
| 2026-04-09 | 各 Agent 独立 SOUL（panda/小U/小F/小B/小Q/小O） |
| 2026-04-13 | pm agent 转型为 planner（内容策划师），职责：行业调研+文案创作，双路线流程（路线A服务文案/路线B独立交付） |
| 2026-04-13 | 小文产出 Main PM 二次核查规则：路线A核查重点数据/路线B核查所有数据结论 |
| 2026-04-13 | 小文交付物统一存放 Projects/{项目名}/docs/，命名 {类型}-{版本}-{日期}.md |
| 2026-04-13 | pm agent 转型为 planner（内容策划师），职责：行业调研+文案创作，双路线流程（路线A服务文案/路线B独立交付） |
| 2026-04-13 | 小文产出 Main PM 二次核查规则：路线A核查重点数据/路线B核查所有数据结论 |
| 2026-04-13 | 小文交付物统一存放 Projects/{项目名}/docs/，命名 {类型}-{版本}-{日期}.md |

---

*详细教训见 `memory/lessons.md`*
