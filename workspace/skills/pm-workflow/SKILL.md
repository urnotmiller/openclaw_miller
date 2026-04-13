---
name: pm-workflow
description: 项目管理流程 skill。当你说'新建项目'、'启动项目'、'启动新项目'、'推进节点'、'推进项目'、'项目流程'时触发。本 skill 仅胖达人(panda)使用，包含作业顺序、调度规则、目录结构规范和模板。
---

# PM Workflow - 项目管理流程

## 触发条件
当 boss 说以下关键词时触发：
- 新建项目
- 启动项目 / 启动新项目
- 推进节点 / 推进项目
- 项目流程

## 核心规则

### 目录结构（新项目必须遵循）

```
Projects/{项目名-版本}/
├── code/
│   ├── backend/
│   ├── frontend/
│   └── common/
├── docs/
│   ├── apis/
│   ├── design/
│   ├── kanban/        ← KANBAN.md 放这里
│   ├── operation/
│   ├── requirements/  ← PRD.md + 访谈纪要.md 放这里
│   └── test/
├── config/
├── scripts/
└── logs/
```

- 看板文件固定命名：`KANBAN.md`
- PRD 文件固定命名：`docs/requirements/PRD.md`
- 访谈纪要固定命名：`docs/requirements/访谈纪要.md`
- 项目路径：`C:\Users\mille\.openclaw\Projects\{项目名-版本}`

### 作业顺序（依赖链）

```
① 需求访谈（结构化提问 → 输出 docs/requirements/访谈纪要.md）
   ↓
② PRD 确认（基于访谈结果 → PRD写入docs/requirements/PRD.md）
   ↓
③ API 接口设计（BE输出，FE确认） ─┐
   ↓                                 ├─ parallel
④ UI 设计（boss审美确认） ─────────┘
   ↓
⑤ BE 开发（等API） ───────────┐
   ↓                           ├─ parallel
   FE 开发（等API+UI双确认） ──┘
   ↓
⑥ 前后端联调（FE主导）
   ↓
⑦ A08-1 QA 代码审核
   ↓ 发现问题 → FE/BE 修复
   ↓
⑧ QA 二次review ← PM必须再次调度QA确认
   ↓
⑨ A08-2 QA 测试执行
   ↓
⑩ 运营物料（OP执行，boss审核）← 可在测试期间并行
   ↓
⑪ A10 部署 → A12 走查 → A13 引流
```

### 调度规则

| agent | 调度方式 | 场景 |
|-------|---------|------|
| BE | spawn | 后端代码任务（一次性交付） |
| FE | spawn | 前端代码任务（一次性交付） |
| QA | spawn | 测试/审核任务（一次性交付） |
| UI | **send + 持久session** | 设计稿需要boss审核，可能多轮修订 |
| OP | **send + 持久session** | 运营物料需要boss审核，可能多轮修订 |

**send + session 操作流程：**
1. `sessions_spawn mode=session` 拉起持久会话
2. `sessions_send` 发任务
3. OP/UI 完成后通过 `sessions_send` 回复
4. boss 审核 → 可能打回修改 → 继续 send 修订

### 关键前置条件

- **FE 启动前置条件**：UI设计稿 + API接口文档 **双确认后**才能调度 FE
- **QA 二次review**：FE/BE 修复完成后，PM 必须再次调度 QA 确认，PM 不能自行判断修复通过
- **UI/OP boss审核**：UI设计稿和运营物料必须 boss 审美/内容确认后才能推进

## 模板文件

- 看板模板：查看 `references/KANBAN-template.md`
- PRD 模板：查看 `references/PRD-template.md`
- 调度规则详情：查看 `references/agent-dispatch.md`
