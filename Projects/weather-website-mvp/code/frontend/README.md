# 森系天气 - 前端

> 天气预报网站 MVP 前端项目

## 运行方式

直接用浏览器打开 `index.html` 即可运行。

```bash
# 或者使用本地服务器（推荐）
npx serve .
# 然后访问 http://localhost:3000
```

## 技术栈

- 纯 HTML + CSS + JavaScript（无框架依赖）
- 响应式设计，移动优先
- 模拟天气数据（无需后端）

## 功能

- ✅ 城市定位与切换
- ✅ 当前天气展示（温度/体感/湿度/风速）
- ✅ 三日预报（昨天/今天/明天）
- ✅ 穿着搭配推荐
- ✅ 天气温馨提示
- ✅ 设置面板（温度单位切换）
- ✅ 城市搜索
- ✅ Toast 通知
- ✅ 底部弹层
- ✅ CSS 动画效果

## 文件结构

```
frontend/
├── index.html   # 主页面
├── styles.css   # 样式表
├── app.js       # 应用逻辑
└── README.md    # 说明文档
```

## 设计规范

- 品牌色：森林绿 `#8DB48E` + 奶油白 `#FFF8F0`
- 圆角：16-24px 大圆角风格
- 字体：PingFang SC / Nunito
- 动画：fadeInUp / scaleUp / 天气动画
