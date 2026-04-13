/**
 * 森系天气后端服务
 * 基于 Express.js 构建
 */

const express = require('express');
const cors = require('cors');
const weatherRoutes = require('./routes/weather');

const app = express();
const PORT = process.env.PORT || 3000;

// 中间件
app.use(cors());
app.use(express.json());

// 请求日志
app.use((req, res, next) => {
  console.log(`[${new Date().toISOString()}] ${req.method} ${req.url}`);
  next();
});

// 路由
app.use('/api/weather', weatherRoutes);

// 健康检查
app.get('/health', (req, res) => {
  res.json({
    code: 0,
    message: 'success',
    data: {
      status: 'ok',
      timestamp: new Date().toISOString()
    }
  });
});

// 404 处理
app.use((req, res) => {
  res.status(404).json({
    code: 404,
    message: '接口不存在',
    data: null
  });
});

// 错误处理
app.use((err, req, res, next) => {
  console.error('服务器错误:', err);
  res.status(500).json({
    code: 500,
    message: '服务器内部错误',
    data: null
  });
});

// 启动服务器
app.listen(PORT, () => {
  console.log(`
╔════════════════════════════════════════════════╗
║     🌿 森系天气后端服务已启动                  ║
╠════════════════════════════════════════════════╣
║  服务地址: http://localhost:${PORT}              ║
║  健康检查: http://localhost:${PORT}/health       ║
║  天气接口: http://localhost:${PORT}/api/weather ║
╚════════════════════════════════════════════════╝
  `);
});

module.exports = app;
