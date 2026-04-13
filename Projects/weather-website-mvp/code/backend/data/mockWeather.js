/**
 * 天气模拟数据
 * 支持的城市天气数据
 */

const mockWeatherData = {
  '北京': {
    cityId: '101010100',
    current: {
      temp: 25,
      feelsLike: 27,
      weather: '晴',
      weatherCode: 'sunny',
      humidity: 65,
      windSpeed: '3级',
      updateTime: '16:30'
    },
    yesterday: { high: 26, low: 18, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    today: { high: 28, low: 19, weather: '晴', icon: '🌤️', weatherCode: 'sunny' },
    tomorrow: { high: 25, low: 17, weather: '多云', icon: '⛅', weatherCode: 'cloudy' }
  },
  '上海': {
    cityId: '101020100',
    current: {
      temp: 22,
      feelsLike: 23,
      weather: '多云',
      weatherCode: 'cloudy',
      humidity: 72,
      windSpeed: '2级',
      updateTime: '16:20'
    },
    yesterday: { high: 23, low: 17, weather: '小雨', icon: '🌧️', weatherCode: 'rainy' },
    today: { high: 24, low: 18, weather: '多云', icon: '⛅', weatherCode: 'cloudy' },
    tomorrow: { high: 22, low: 16, weather: '雨', icon: '🌧️', weatherCode: 'rainy' }
  },
  '广州': {
    cityId: '101280101',
    current: {
      temp: 29,
      feelsLike: 32,
      weather: '雷阵雨',
      weatherCode: 'thunder',
      humidity: 85,
      windSpeed: '3级',
      updateTime: '16:00'
    },
    yesterday: { high: 31, low: 24, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    today: { high: 30, low: 25, weather: '雷阵雨', icon: '⛈️', weatherCode: 'thunder' },
    tomorrow: { high: 28, low: 23, weather: '大雨', icon: '🌧️', weatherCode: 'heavyRain' }
  },
  '深圳': {
    cityId: '101280601',
    current: {
      temp: 28,
      feelsLike: 31,
      weather: '阴',
      weatherCode: 'overcast',
      humidity: 80,
      windSpeed: '2级',
      updateTime: '16:15'
    },
    yesterday: { high: 30, low: 24, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    today: { high: 29, low: 24, weather: '阴', icon: '☁️', weatherCode: 'overcast' },
    tomorrow: { high: 27, low: 23, weather: '小雨', icon: '🌧️', weatherCode: 'rainy' }
  },
  '杭州': {
    cityId: '101210101',
    current: {
      temp: 20,
      feelsLike: 19,
      weather: '小雨',
      weatherCode: 'rainy',
      humidity: 88,
      windSpeed: '2级',
      updateTime: '16:25'
    },
    yesterday: { high: 24, low: 16, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    today: { high: 21, low: 15, weather: '小雨', icon: '🌧️', weatherCode: 'rainy' },
    tomorrow: { high: 19, low: 14, weather: '雨', icon: '🌧️', weatherCode: 'rainy' }
  },
  '成都': {
    cityId: '101270101',
    current: {
      temp: 18,
      feelsLike: 17,
      weather: '多云',
      weatherCode: 'cloudy',
      humidity: 75,
      windSpeed: '1级',
      updateTime: '16:10'
    },
    yesterday: { high: 20, low: 14, weather: '阴', icon: '☁️', weatherCode: 'overcast' },
    today: { high: 19, low: 14, weather: '多云', icon: '⛅', weatherCode: 'cloudy' },
    tomorrow: { high: 21, low: 15, weather: '晴', icon: '☀️', weatherCode: 'sunny' }
  },
  '重庆': {
    cityId: '101040100',
    current: {
      temp: 23,
      feelsLike: 24,
      weather: '晴',
      weatherCode: 'sunny',
      humidity: 68,
      windSpeed: '2级',
      updateTime: '16:30'
    },
    yesterday: { high: 25, low: 19, weather: '多云', icon: '⛅', weatherCode: 'cloudy' },
    today: { high: 24, low: 18, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    tomorrow: { high: 26, low: 19, weather: '晴', icon: '☀️', weatherCode: 'sunny' }
  },
  '武汉': {
    cityId: '101200101',
    current: {
      temp: 22,
      feelsLike: 23,
      weather: '晴',
      weatherCode: 'sunny',
      humidity: 65,
      windSpeed: '3级',
      updateTime: '16:20'
    },
    yesterday: { high: 24, low: 16, weather: '阴', icon: '☁️', weatherCode: 'overcast' },
    today: { high: 23, low: 17, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    tomorrow: { high: 25, low: 18, weather: '多云', icon: '⛅', weatherCode: 'cloudy' }
  },
  '西安': {
    cityId: '101110101',
    current: {
      temp: 19,
      feelsLike: 18,
      weather: '多云',
      weatherCode: 'cloudy',
      humidity: 55,
      windSpeed: '2级',
      updateTime: '16:15'
    },
    yesterday: { high: 21, low: 13, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    today: { high: 20, low: 14, weather: '多云', icon: '⛅', weatherCode: 'cloudy' },
    tomorrow: { high: 22, low: 15, weather: '晴', icon: '☀️', weatherCode: 'sunny' }
  },
  '南京': {
    cityId: '101190101',
    current: {
      temp: 21,
      feelsLike: 22,
      weather: '晴',
      weatherCode: 'sunny',
      humidity: 70,
      windSpeed: '2级',
      updateTime: '16:25'
    },
    yesterday: { high: 23, low: 15, weather: '阴', icon: '☁️', weatherCode: 'overcast' },
    today: { high: 22, low: 16, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    tomorrow: { high: 24, low: 17, weather: '多云', icon: '⛅', weatherCode: 'cloudy' }
  },
  '苏州': {
    cityId: '101190401',
    current: {
      temp: 22,
      feelsLike: 23,
      weather: '晴',
      weatherCode: 'sunny',
      humidity: 68,
      windSpeed: '2级',
      updateTime: '16:20'
    },
    yesterday: { high: 24, low: 16, weather: '多云', icon: '⛅', weatherCode: 'cloudy' },
    today: { high: 23, low: 17, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    tomorrow: { high: 25, low: 18, weather: '多云', icon: '⛅', weatherCode: 'cloudy' }
  },
  '天津': {
    cityId: '101030100',
    current: {
      temp: 24,
      feelsLike: 25,
      weather: '晴',
      weatherCode: 'sunny',
      humidity: 60,
      windSpeed: '3级',
      updateTime: '16:30'
    },
    yesterday: { high: 26, low: 17, weather: '多云', icon: '⛅', weatherCode: 'cloudy' },
    today: { high: 25, low: 18, weather: '晴', icon: '☀️', weatherCode: 'sunny' },
    tomorrow: { high: 27, low: 19, weather: '晴', icon: '☀️', weatherCode: 'sunny' }
  }
};

/**
 * 获取所有支持的城市列表
 */
function getSupportedCities() {
  return Object.keys(mockWeatherData).map(cityName => ({
    cityName,
    cityId: mockWeatherData[cityName].cityId
  }));
}

/**
 * 根据城市名获取天气数据
 * @param {string} cityName - 城市名
 * @returns {object|null} 天气数据
 */
function getWeatherByCity(cityName) {
  const data = mockWeatherData[cityName];
  if (!data) {
    return null;
  }

  // 更新一下时间戳（模拟真实场景）
  const now = new Date();
  const hours = now.getHours().toString().padStart(2, '0');
  const minutes = now.getMinutes().toString().padStart(2, '0');

  return {
    ...data,
    current: {
      ...data.current,
      updateTime: `${hours}:${minutes}`
    }
  };
}

module.exports = {
  mockWeatherData,
  getSupportedCities,
  getWeatherByCity
};
