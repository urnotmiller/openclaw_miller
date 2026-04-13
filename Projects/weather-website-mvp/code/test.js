/**
 * 森系天气 - API 冒烟测试
 * 运行方式: node code/test.js
 */
const http = require('http');

const BASE = 'http://localhost:3000';
const PASS = '✅';
const FAIL = '❌';

function get(path) {
  return new Promise((resolve, reject) => {
    http.get(`${BASE}${path}`, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try { resolve(JSON.parse(data)); }
        catch { resolve(data); }
      });
    }).on('error', reject);
  });
}

async function test(name, fn) {
  try {
    const result = await fn();
    console.log(`${result ? PASS : FAIL} ${name}`);
    return result;
  } catch (e) {
    console.log(`${FAIL} ${name}: ${e.message}`);
    return false;
  }
}

async function run() {
  console.log('\n🌿 森系天气 API 冒烟测试\n');
  
  await test('健康检查', async () => {
    const r = await get('/health');
    return r.data && r.data.status === 'ok';
  });

  await test('获取北京天气', async () => {
    const r = await get('/api/weather/北京');
    return r.data && r.data.current && r.data.yesterday && r.data.today && r.data.tomorrow;
  });

  await test('获取上海天气', async () => {
    const r = await get('/api/weather/上海');
    return r.data && r.data.cityId === '101020100';
  });

  await test('无效城市返回错误', async () => {
    const r = await get('/api/weather/不存在城市');
    return r.code !== 0 || (r.message && r.message.includes('不存在'));
  });

  await test('城市列表', async () => {
    const r = await get('/api/weather/cities/list');
    return Array.isArray(r.data) && r.data.length === 12;
  });

  console.log('\n测试完成\n');
}

run();
