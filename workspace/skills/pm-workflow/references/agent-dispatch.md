# Agent 调度规则

## 调度方式总览

| agent | 调度方式 | 适用场景 |
|-------|---------|---------|
| BE | `sessions_spawn` | 后端代码开发（一次性交付） |
| FE | `sessions_spawn` | 前端代码开发（一次性交付） |
| QA | `sessions_spawn` | 代码审核/测试执行（一次性交付） |
| UI | `sessions_spawn mode=session` + `sessions_send` | 设计稿需要多轮修订，boss 审核 |
| OP | `sessions_spawn mode=session` + `sessions_send` | 运营物料需要多轮修订，boss 审核 |

---

## sessions_spawn（一次性任务）

```javascript
sessions_spawn({
  task: "任务描述，包含项目路径、输出目录、技术要求",
  runtime: "subagent",
  agentId: "be", // be / fe / qa / op
  runTimeoutSeconds: 600,
  mode: "run",
  cleanup: "keep"
})
```

**适用场景：** 代码开发、测试执行、接口文档输出等一次性交付任务。

---

## sessions_send + 持久 session（需要 boss 审核的任务）

### 第一步：拉起持久会话
```javascript
sessions_spawn({
  task: "你是 {agent}，持续待命，等待任务指令。",
  runtime: "subagent",
  agentId: "ui", // ui / op
  mode: "session",  // 持久会话
  cleanup: "keep"
})
```
→ 返回 `childSessionKey`，如 `agent:ui:subagent:xxxxx`

### 第二步：发送任务
```javascript
sessions_send({
  sessionKey: "agent:ui:subagent:xxxxx",
  message: "【任务描述】..."
})
```

### 第三步：OP/UI 通过 sessions_send 回复结果
boss 审核 → 可能打回修改 → 继续 send 修订

### 第四步：任务完成
boss 确认通过后，该 session 可以结束。

---

## 关键调度规则

### FE 启动前置条件
FE 开始开发的时机：**UI设计稿 + API接口文档 双确认后**，两者缺一不可。
- UI 必须 boss 审美确认
- API 必须 BE 输出、FE 确认后冻结
- 未双确认不能调度 FE

### QA 二次 review 规则
- FE/BE 修复完成后，PM **必须再次调度 QA** 进行二次 review
- PM 不能自行判断修复通过，必须等 QA 确认
- 二次 review 通过后才能进入 A08-2 测试执行

### UI/OP boss 审核规则
- UI 设计稿完成后必须 boss 审美确认
- 运营物料完成后必须 boss 内容确认
- 中途打回修改，用 send 继续发修订指令

---

## 各 Agent 职责边界

| Agent | 职责 | 不做什么 |
|-------|------|---------|
| BE | 后端开发、API文档输出 | 不直接与 boss 沟通代码细节 |
| FE | 前端开发、联调主导 | 不自行 spawn 其他 agent |
| QA | 代码审核、测试执行 | 不自行修复，发现问题上报 PM |
| UI | 设计稿输出 | 不直接与 boss 沟通，通过 PM 中转 |
| OP | 运营物料、部署文档 | 不自行决定内容方向，按需求执行 |

---

## 调度串行/并行规则

| 关系 | 时序 |
|------|------|
| BE 和 UI | PRD 确认后 **并行**，互相独立 |
| FE 和 BE | API 文档冻结后 **并行**，各自开发 |
| FE 和 UI | FE 等 UI 确认，但 BE 不等 UI |
| QA | 等联调完成后才能开始；测试用例可提前写 |
| OP | 测试期间可并行；主要等 QA 通过后启动 |
