# 天气 API 接口文档

> **版本**：v1.0  
> **更新日期**：2026-04-05  
> **Base URL**：`http://localhost:3000/api`

---

## 1. 获取天气数据

### 请求

```
GET /api/weather/:cityName
```

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| cityName | string | 是 | 城市名称（中文），如"北京"、"上海" |

### 成功响应

**HTTP 200**

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
    "yesterday": {
      "high": 26,
      "low": 18,
      "weather": "晴",
      "icon": "☀️",
      "weatherCode": "sunny"
    },
    "today": {
      "high": 28,
      "low": 19,
      "weather": "晴",
      "icon": "🌤️",
      "weatherCode": "sunny"
    },
    "tomorrow": {
      "high": 25,
      "low": 17,
      "weather": "多云",
      "icon": "⛅",
      "weatherCode": "cloudy"
    }
  }
}
```

### 错误响应

**HTTP 404** - 城市不存在

```json
{
  "code": 404,
  "message": "城市不存在",
  "data": null
}
```

**HTTP 500** - 服务器错误

```json
{
  "code": 500,
  "message": "获取天气数据失败",
  "data": null
}
```

---

## 2. 数据字段说明

### 2.1 天气代码 (weatherCode)

| code | 说明 | 图标 |
|------|------|------|
| `sunny` | 晴天 | ☀️ |
| `cloudy` | 多云 | ⛅ |
| `overcast` | 阴天 | ☁️ |
| `rainy` | 小雨 | 🌧️ |
| `heavyRain` | 大雨 | 🌧️ |
| `thunder` | 雷阵雨/雷暴 | ⛈️ |

### 2.2 温度字段

- `temp`：当前温度（摄氏度）
- `feelsLike`：体感温度
- `high`：最高温度
- `low`：最低温度

### 2.3 其他字段

- `humidity`：湿度百分比
- `windSpeed`：风力等级（字符串，如"3级"）
- `updateTime`：数据更新时间（HH:mm 格式）

---

## 3. 支持的城市列表

| 城市 | cityId |
|------|--------|
| 北京 | 101010100 |
| 上海 | 101020100 |
| 广州 | 101280101 |
| 深圳 | 101280601 |
| 杭州 | 101210101 |
| 成都 | 101270101 |
| 重庆 | 101040100 |
| 武汉 | 101200101 |
| 西安 | 101110101 |
| 南京 | 101190101 |
| 苏州 | 101190401 |
| 天津 | 101030100 |

---

## 4. 前端对接说明

前端 `WeatherAPI.getWeather(cityName)` 返回的格式与本接口返回的 `data` 字段完全一致，可无缝切换。

### 切换步骤

1. 将 `app.js` 中的 `WeatherAPI.getWeather` 改为调用后端接口
2. 接口地址：`http://localhost:3000/api/weather/{cityName}`
3. 响应数据中取 `response.data` 即为原 `WeatherAPI.getWeather` 返回值

### 示例代码

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
