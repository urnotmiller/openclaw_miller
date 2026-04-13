/**
 * 天气路由
 */

const express = require('express');
const router = express.Router();
const weatherService = require('../services/weatherService');

/**
 * GET /api/weather/:cityName
 * 获取指定城市的天气数据
 */
router.get('/:cityName', async (req, res) => {
  try {
    const { cityName } = req.params;

    if (!cityName) {
      return res.status(400).json({
        code: 400,
        message: '城市名不能为空',
        data: null
      });
    }

    const weather = await weatherService.getWeather(cityName);

    res.json({
      code: 0,
      message: 'success',
      data: weather
    });
  } catch (error) {
    console.error('获取天气数据失败:', error);

    const code = error.code || 500;
    const message = error.message || '获取天气数据失败';

    res.status(code).json({
      code,
      message,
      data: null
    });
  }
});

/**
 * GET /api/weather/cities/list
 * 获取支持的城市列表
 */
router.get('/cities/list', async (req, res) => {
  try {
    const cities = await weatherService.getCities();

    res.json({
      code: 0,
      message: 'success',
      data: cities
    });
  } catch (error) {
    console.error('获取城市列表失败:', error);

    res.status(500).json({
      code: 500,
      message: '获取城市列表失败',
      data: null
    });
  }
});

module.exports = router;
