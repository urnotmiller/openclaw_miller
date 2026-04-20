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

## 🤖 Agent 矩阵

| Agent | 角色 | 工作区 |
|-------|------|--------|
| panda（我） | 项目经理 PM | `workspace` |
| 小U | UI设计 | `workspace-ui` |
| 小F | 前端开发 | `workspace-fe` |
| 小B | 后端开发 | `workspace-be` |
| 小Q | 测试/审核 | `workspace-qa` |
| 小O | 运营 | `workspace-op` |
| 小文 | 内容策划（planner） | `workspace-planner` |

---

## ⚡ 核心规则（不常改动）

### 记忆管理（复盘=双写）
每次复盘时，两边同步更新：
- `memory/projects.md` — 项目节点变化
- `memory/lessons.md` — 教训分级（🔴严重/🟡中等/🟢轻微）
- `.learnings/LEARNINGS.md` — 教训、修正、洞察
- `.learnings/ERRORS.md` — 错误/失败记录
- `.learnings/FEATURE_REQUESTS.md` — 做不到的需求

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
| 2026-04-07 | 新增 QA agent（测试/审核），飞书 qa 账号绑定，长连接模式 |
| 2026-04-13 | pm agent 转型为 planner（内容策划师），职责：行业调研+文案创作，双路线流程（路线A服务文案/路线B独立交付） |
| 2026-04-13 | 小文产出 Main PM 二次核查规则：路线A核查重点数据/路线B核查所有数据结论 |
| 2026-04-13 | 小文交付物统一存放 Projects/{项目名}/docs/，命名 {类型}-{版本}-{日期}.md |
| 2026-04-15 | 飞书云盘目录结构体系建立：01项目文档/02案例库/03运营素材/04知识沉淀/05归档 |
| 2026-04-15 | .openclaw 目录同步至 GitHub（urnotmiller/openclaw_miller），排除 credentials/、*.sqlite、backups/、.env |
| 2026-04-16 | PM协作规范升级：Karpathy协作四原则（Think/Simplicity/Surgical/Goal-Driven），SOUL.md留原则、AGENTS.md留操作手册 |

---

*详细教训见 `memory/lessons.md`*
