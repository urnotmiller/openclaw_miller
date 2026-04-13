/**
 * 森系天气 - 一键启动脚本
 * 同时启动前端静态服务器和后端服务
 */

const { spawn } = require('child_process');
const path = require('path');
const { execSync } = require('child_process');

const ROOT = __dirname;
const FRONTEND_PORT = 3001;
const BACKEND_PORT = 3000;

// 颜色输出
const green = (msg) => console.log(`\x1b[32m${msg}\x1b[0m`);
const red = (msg) => console.log(`\x1b[31m${msg}\x1b[0m`);
const cyan = (msg) => console.log(`\x1b[36m${msg}\x1b[0m`);

// 启动后端
green('[后端] 正在启动...');
const backend = spawn('node', ['server.js'], {
  cwd: path.join(ROOT, 'backend'),
  stdio: ['pipe', 'pipe', 'pipe'],
  shell: true,
  env: { ...process.env, PORT: BACKEND_PORT }
});

backend.stdout.on('data', (data) => {
  process.stdout.write(`[BE] ${data}`);
});

backend.stderr.on('data', (data) => {
  process.stderr.write(`[BE ERROR] ${data}`);
});

backend.on('error', (err) => {
  red(`[后端] 启动失败: ${err.message}`);
});

// 启动前端
cyan('[前端] 正在启动...');
const frontend = spawn('npx', ['serve', '.', '-p', String(FRONTEND_PORT), '-s'], {
  cwd: path.join(ROOT, 'frontend'),
  stdio: ['pipe', 'pipe', 'pipe'],
  shell: true
});

frontend.stdout.on('data', (data) => {
  process.stdout.write(`[FE] ${data}`);
});

frontend.stderr.on('data', (data) => {
  process.stderr.write(`[FE ERROR] ${data}`);
});

frontend.on('error', (err) => {
  red(`[前端] 启动失败: ${err.message}`);
});

// 等待后端就绪后自动打开浏览器
setTimeout(() => {
  const url = `http://localhost:${FRONTEND_PORT}`;
  green(`\n🌿 森系天气已启动`);
  console.log(`   前端: ${url}`);
  console.log(`   后端: http://localhost:${BACKEND_PORT}`);

  try {
    execSync(`start ${url}`, { shell: true, cwd: ROOT });
    green(`   ✅ 浏览器已打开`);
  } catch (e) {
    cyan(`   请手动打开浏览器访问: ${url}`);
  }
}, 4000);

// 优雅退出
const cleanup = () => {
  green('\n[关闭] 正在停止服务...');
  backend.kill();
  frontend.kill();
  process.exit(0);
};

process.on('SIGINT', cleanup);
process.on('SIGTERM', cleanup);
