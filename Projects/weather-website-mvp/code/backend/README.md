# 森系天气后端服务

> 天气网站 MVP 后端 API 服务

## 快速启动

```bash
# 安装依赖
npm install

# 开发模式（热重载）
npm run dev

# 生产模式
npm start
```

服务启动后访问 http://localhost:3000

## API 接口

### 获取天气数据

```
GET http://localhost:3000/api/weather/:cityName
```

**示例**

```bash
curl http://localhost:3000/api/weather/北京
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "cityId": "101010100",
    "current": {
      "temp": 25,
      "feelsLike": 27,
      "weather": "晴",
      "weatherCode": "sunny",
      "humidity": 65,
      "windSpeed": "3级",
      "updateTime": "16:30"
    },
    "yesterday": { "high": 26, "low": 18, "weather": "晴", "icon": "☀️", "weatherCode": "sunny" },
    "today": { "high": 28, "low": 19, "weather": "晴", "icon": "🌤️", "weatherCode": "sunny" },
    "tomorrow": { "high": 25, "low": 17, "weather": "多云", "icon": "⛅", "weatherCode": "cloudy" }
  }
}
```

### 获取支持城市列表

```
GET http://localhost:3000/api/weather/cities/list
```

### 健康检查

```
GET http://localhost:3000/health
```

## 项目结构

```
code/backend/
├── server.js              # 服务入口
├── package.json
├── routes/
│   └── weather.js        # 天气路由
├── services/
│   └── weatherService.js # 业务逻辑
└── data/
    └── mockWeather.js    # 模拟数据
```

## 支持的城市

北京、上海、广州、深圳、杭州、成都、重庆、武汉、西安、南京、苏州、天津

## 前端对接

详见 `docs/apis/weather-api.md`

将前端 `app.js` 中的 `WeatherAPI.getWeather` 改为调用：

```javascript
async getWeather(cityName) {
  const response = await fetch(`http://localhost:3000/api/weather/${encodeURIComponent(cityName)}`);
  const result = await response.json();
  if (result.code !== 0) {
    throw new Error(result.message);
  }
  return result.data;
}
```
