/**
 * 天气服务层
 * 负责天气数据的获取和处理
 */

const { getWeatherByCity, getSupportedCities } = require('../data/mockWeather');

class WeatherService {
  /**
   * 获取城市天气
   * @param {string} cityName - 城市名
   * @returns {Promise<object>} 天气数据
   */
  async getWeather(cityName) {
    const weather = getWeatherByCity(cityName);

    if (!weather) {
      const error = new Error('城市不存在');
      error.code = 404;
      throw error;
    }

    return weather;
  }

  /**
   * 获取支持的城市列表
   * @returns {Array} 城市列表
   */
  async getCities() {
    return getSupportedCities();
  }
}

module.exports = new WeatherService();
