# Learnings

Corrections, insights, and knowledge gaps captured during development.

**Categories**: correction | insight | knowledge_gap | best_practice

---

## [LRN-20260409-001] best_practice

**Logged**: 2026-04-09T11:45:00+08:00
**Priority**: high
**Status**: active
**Area**: agent-dispatch

### Summary
UI 和 OP 的任务需要多轮修订，应使用 sessions_send + 持久 session 调度，不应用 sessions_spawn。spawn 是一次性任务，无法支持反馈-修改循环。

### Details
- 适用场景：UI设计稿（审美见仁见智）、运营物料（内容需审核调整）
- spawn 适合：BE/FE/QA 的代码任务（一次性交付）
- send+session 操作：先用 spawn mode=session 拉起，再 send 发任务，OP/UI 完成通过 send 回复

### Metadata
- Source: weather-website-mvp project review
- Pattern-Key: agent.dispatch.send-vs-spawn

---

## [LRN-20260409-002] best_practice

**Logged**: 2026-04-09T11:45:00+08:00
**Priority**: high
**Status**: active
**Area**: workflow

### Summary
FE 启动的前置条件必须是 UI设计稿 + API接口文档 双确认，缺一不可。避免 FE 基于未经确认的设计稿开发。

### Details
- UI 必须 boss 审美确认
- API 必须 BE 输出、FE 确认后冻结
- 未双确认不能调度 FE

### Metadata
- Source: weather-website-mvp project review
- Pattern-Key: workflow.fe-prerequisite

---

## [LRN-20260409-003] best_practice

**Logged**: 2026-04-09T11:45:00+08:00
**Priority**: high
**Status**: active
**Area**: qa-process

### Summary
QA 代码审核发现的问题，修复后必须再次调度 QA 进行二次 review 确认，PM 不能自行判断修复通过。

### Details
- 流程：QA审核 → 发现问题 → FE/BE修复 → QA二次review → 通过后 → A08-2测试执行
- PM 不得自行判断修复有效性
- 二次 review 是 QA 一票否决权的体现

### Metadata
- Source: weather-website-mvp project review
- Pattern-Key: qa.second-review

---

## [LRN-20260409-004] best_practice

**Logged**: 2026-04-09T11:45:00+08:00
**Priority**: medium
**Status**: active
**Area**: requirements

### Summary
大项目需求沟通应使用结构化访谈，提升效率。访谈输出访谈纪要，再基于纪要输出 PRD。

### Details
- 访谈框架：商业目标/用户画像/核心场景/约束条件/优先级
- 输出：docs/requirements/访谈纪要.md
- PRD 基于访谈纪要产出

### Metadata
- Source: weather-website-mvp project review
- Pattern-Key: requirements.interview-flow

---

## [LRN-20260409-005] insight

**Logged**: 2026-04-09T11:45:00+08:00
**Priority**: medium
**Status**: active
**Area**: workflow

### Summary
BE 和 UI 在 PRD 确认后可以并行开发，互相不依赖。FE 和 BE 都只依赖 API 文档冻结。

### Details
- 并行关系：BE开发 和 UI设计（都只等 PRD）
- FE 等 UI + API 双确认
- QA 测试用例可提前写，实际测试等联调完成

### Metadata
- Source: weather-website-mvp project review
- Pattern-Key: workflow.parallel-execution

---

## [LRN-20260404-001] best_practice

**Logged**: 2026-04-04T14:26:00+08:00
**Priority**: medium
**Status**: pending
**Area**: config

### Summary
OpenClaw执行本地命令需要先完成设备配对，有pending配对请求时需要用户批准才能正常使用

### Details
- 初始状态下，新会话会有一个待处理的配对请求
- 执行命令会报错 "pairing required"
- 需要用户执行 `openclaw devices approve --latest` 或者批准具体请求ID才能继续
- 配对批准后才能正常使用git、clawhub等本地命令

### Suggested Action
- 遇到 "pairing required" 错误时，直接引导用户批准pending配对
- 可以用 `openclaw devices list` 查看待处理请求

### Metadata
- Source: conversation
- Related Files: 
- Tags: openclaw, pairing, auth
- Pattern-Key: openclaw.pairing.required

