/**
 * 森系天气 - Forest Cute Weather App
 * Main Application JavaScript
 */

// ============================================
// Weather Data (API Integration)
// ============================================

const WeatherAPI = {
  // Backend API base URL
  BASE_URL: 'http://localhost:3000',

  // Weather code to icon/emoji mapping
  weatherIconMap: {
    'sunny': '☀️',
    'cloudy': '⛅',
    'overcast': '☁️',
    'rainy': '🌧️',
    'heavyRain': '🌧️',
    'thunder': '⛈️',
    'snow': '❄️',
    'fog': '🌫️',
    'windy': '💨'
  },

  // Weather code to frontend weatherCode mapping
  weatherCodeMap: {
    'CLEAR': 'sunny',
    'SUNNY': 'sunny',
    'PARTLY_CLOUDY': 'cloudy',
    'CLOUDY': 'cloudy',
    'OVERCAST': 'overcast',
    'RAIN': 'rainy',
    'LIGHT_RAIN': 'rainy',
    'HEAVY_RAIN': 'heavyRain',
    'THUNDER': 'thunder',
    'SNOW': 'snow',
    'FOG': 'fog',
    'WIND': 'windy'
  },

  // Map backend response to frontend expected format
  _mapBackendData(apiData) {
    const weatherCodeRaw = apiData.weatherCode || apiData.weather_code || apiData.current?.weatherCode || 'sunny';
    const weatherCode = this.weatherCodeMap[weatherCodeRaw] || weatherCodeRaw;
    const icon = this.weatherIconMap[weatherCode] || '☀️';

    // Extract update time from API or use current time
    const updateTime = apiData.updateTime || apiData.update_time || apiData.current?.updateTime || new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });

    return {
      cityId: apiData.cityId || apiData.city_id || '',
      current: {
        temp: Number(apiData.temp) || Number(apiData.current?.temp) || 0,
        feelsLike: Number(apiData.feelsLike) || Number(apiData.feels_like) || Number(apiData.current?.feelsLike) || 0,
        weather: apiData.weather || apiData.current?.weather || '晴',
        weatherCode: weatherCode,
        humidity: Number(apiData.humidity) || Number(apiData.current?.humidity) || 0,
        windSpeed: apiData.windSpeed || apiData.wind_speed || apiData.current?.windSpeed || '0级',
        updateTime: updateTime
      },
      yesterday: {
        high: Number(apiData.yesterday?.high) || Number(apiData.forecast?.[-1]?.high) || 0,
        low: Number(apiData.yesterday?.low) || Number(apiData.forecast?.[-1]?.low) || 0,
        weather: apiData.yesterday?.weather || apiData.forecast?.[-1]?.weather || '晴',
        icon: this.weatherIconMap[apiData.yesterday?.weatherCode] || this.weatherIconMap[weatherCode] || '☀️',
        weatherCode: this.weatherCodeMap[apiData.yesterday?.weatherCode] || weatherCode
      },
      today: {
        high: Number(apiData.today?.high) || Number(apiData.forecast?.[0]?.high) || 0,
        low: Number(apiData.today?.low) || Number(apiData.forecast?.[0]?.low) || 0,
        weather: apiData.today?.weather || apiData.forecast?.[0]?.weather || '晴',
        icon: this.weatherIconMap[apiData.today?.weatherCode] || icon,
        weatherCode: this.weatherCodeMap[apiData.today?.weatherCode] || weatherCode
      },
      tomorrow: {
        high: Number(apiData.tomorrow?.high) || Number(apiData.forecast?.[1]?.high) || 0,
        low: Number(apiData.tomorrow?.low) || Number(apiData.forecast?.[1]?.low) || 0,
        weather: apiData.tomorrow?.weather || apiData.forecast?.[1]?.weather || '多云',
        icon: this.weatherIconMap[apiData.tomorrow?.weatherCode] || this.weatherIconMap[weatherCode] || '⛅',
        weatherCode: this.weatherCodeMap[apiData.tomorrow?.weatherCode] || weatherCode
      }
    };
  },

  // Get weather by city name
  async getWeather(cityName) {
    try {
      const response = await fetch(`${this.BASE_URL}/api/weather/${encodeURIComponent(cityName)}`);

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('城市不存在');
        }
        throw new Error(`请求失败: ${response.status}`);
      }

      const apiData = await response.json();

      // Handle different response structures
      const rawData = apiData.data || apiData.result || apiData;

      // Map backend data to frontend expected format
      return this._mapBackendData(rawData);

    } catch (error) {
      if (error.message === '城市不存在') {
        throw error;
      }
      throw new Error(`网络错误: ${error.message}`);
    }
  },

  // Get user's current location
  async getCurrentLocation() {
    return new Promise((resolve, reject) => {
      if (!navigator.geolocation) {
        reject(new Error('您的浏览器不支持定位'));
        return;
      }

      navigator.geolocation.getCurrentPosition(
        async (position) => {
          try {
            // In a real app, we'd reverse geocode to get city name
            // For now, default to Beijing
            resolve('北京');
          } catch (error) {
            resolve('北京');
          }
        },
        (error) => {
          switch (error.code) {
            case error.PERMISSION_DENIED:
              reject(new Error('定位权限未开启'));
              break;
            case error.POSITION_UNAVAILABLE:
              reject(new Error('无法获取位置信息'));
              break;
            case error.TIMEOUT:
              reject(new Error('定位超时'));
              break;
            default:
              reject(new Error('定位失败'));
          }
        },
        { timeout: 10000 }
      );
    });
  },

  _delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
};

// ============================================
// Outfit Recommendation Engine
// ============================================

const OutfitEngine = {
  // Outfit recommendations based on weather conditions
  recommendations: {
    sunny: [
      {
        tags: ['☀️ 温暖', '🌬 微风'],
        icons: ['👕', '👖', '👟', '🕶️'],
        top: '短袖T恤<br>薄衬衫',
        bottom: '短裤<br>轻薄长裤',
        shoes: '凉鞋<br>运动鞋',
        accessories: '太阳镜<br>遮阳帽',
        quote: '今天阳光正好，适合浅色系穿搭，元气满满出门吧～'
      },
      {
        tags: ['☀️ 炎热', '💧 干燥'],
        icons: ['👕', '🩳', '🩴', '🧴'],
        top: '透气短袖<br>冰丝上衣',
        bottom: '宽松短裤<br>薄款长裤',
        shoes: '凉拖<br>帆布鞋',
        accessories: '防晒霜<br>大容量水壶',
        quote: '紫外线较强，记得做好防晒，多补充水分哦～'
      }
    ],
    cloudy: [
      {
        tags: ['🌤️ 舒适', '🌬 微风'],
        icons: ['👔', '👖', '👟', '🧣'],
        top: '薄长袖<br>针织开衫',
        bottom: '牛仔裤<br>休闲长裤',
        shoes: '运动鞋<br>平底鞋',
        accessories: '围巾<br>帽子',
        quote: '多云天气，温度适宜，穿搭可以轻薄又时尚～'
      },
      {
        tags: ['☁️ 阴天', '💧 微凉'],
        icons: ['👔', '👖', '🥾', '🎩'],
        top: '长袖衬衫<br>卫衣',
        bottom: '长裤<br>针织裙',
        shoes: '运动鞋<br>靴子',
        accessories: '薄外套<br>帽子',
        quote: '天气转凉，建议带件外套再出门，小心感冒～'
      }
    ],
    rainy: [
      {
        tags: ['🌧️ 小雨', '🌬 湿润'],
        icons: ['🧥', '👖', '🥾', '☂️'],
        top: '防水夹克<br>连帽卫衣',
        bottom: '速干长裤<br>牛仔裤',
        shoes: '防水鞋<br>洞洞鞋',
        accessories: '折叠伞<br>防水包',
        quote: '记得带伞！雨天路滑，走路开车都要小心哦～'
      },
      {
        tags: ['🌧️ 有雨', '💧 温差大'],
        icons: ['🧥', '👖', '🥾', '☂️'],
        top: '薄羽绒服<br>保暖内衣',
        bottom: '加绒长裤<br>休闲裤',
        shoes: '防水靴<br>运动鞋',
        accessories: '雨伞<br>暖宝宝',
        quote: '雨天降温明显，注意保暖别着凉啦～'
      }
    ],
    heavyRain: [
      {
        tags: ['🌧️ 大雨', '💨 大风'],
        icons: ['🧥', '👖', '🥾', '☂️'],
        top: '冲锋衣<br>保暖内衣',
        bottom: '防水长裤<br>加绒裤',
        shoes: '防水靴<br>雨靴',
        accessories: '大伞<br>防水包<br>暖宝宝',
        quote: '大雨来袭！尽量减少外出，安全第一哦～'
      }
    ],
    thunder: [
      {
        tags: ['⛈️ 雷暴', '🌧️ 阵雨'],
        icons: ['🧥', '👖', '🥾', '☂️'],
        top: '防水夹克<br>速干衣',
        bottom: '防水长裤<br>长裤',
        shoes: '橡胶鞋<br>防水靴',
        accessories: '折叠伞<br>防水包',
        quote: '雷暴天气！尽量待在室内，关好门窗注意安全～'
      }
    ],
    overcast: [
      {
        tags: ['☁️ 阴天', '🌬 舒适'],
        icons: ['👔', '👖', '👟', '🧢'],
        top: '薄长袖<br>针织衫',
        bottom: '牛仔裤<br>长裙',
        shoes: '运动鞋<br>帆布鞋',
        accessories: '薄外套<br>帽子',
        quote: '阴天也要有好心情，穿得舒适最重要～'
      }
    ]
  },

  // Get recommendation based on weather code
  getRecommendation(weatherCode, temp) {
    let category = weatherCode;

    // Map weather codes to categories
    if (weatherCode === 'sunny') {
      category = temp >= 28 ? 'sunny-hot' : 'sunny';
    }

    const recs = this.recommendations[weatherCode] || this.recommendations.cloudy;
    const index = Math.floor(Math.random() * recs.length);
    return recs[index];
  },

  // Get alternative recommendation
  getAlternative(currentRec, weatherCode) {
    const recs = this.recommendations[weatherCode] || this.recommendations.cloudy;
    const alternatives = recs.filter((_, i) => i !== recs.indexOf(currentRec));
    if (alternatives.length > 0) {
      return alternatives[Math.floor(Math.random() * alternatives.length)];
    }
    return recs[0];
  }
};

// ============================================
// Tips Engine
// ============================================

const TipsEngine = {
  // Tip definitions
  tipTypes: {
    umbrella: {
      icon: '🌂',
      title: '带伞',
      bgClass: 'umbrella',
      getDesc: (weather) => {
        if (weather.includes('雨') || weather.includes('雷')) return '傍晚有雨';
        return '午后可能有阵雨';
      },
      action: '查看详情'
    },
    clothing: {
      icon: '👕',
      title: '加件外套',
      bgClass: 'clothing',
      getDesc: (weather, temp) => {
        if (temp < 20) return '早晚温差大';
        return '气温略低，注意保暖';
      },
      action: '查看穿搭'
    },
    sunscreen: {
      icon: '🧴',
      title: '防晒',
      bgClass: 'sunscreen',
      getDesc: (weather) => {
        return '紫外线较强';
      },
      action: '了解详情'
    },
    water: {
      icon: '💧',
      title: '补水',
      bgClass: 'water',
      getDesc: (weather, temp) => {
        if (temp > 28) return '天气干燥，多喝水';
        return '空气湿度较低';
      },
      action: '健康提示'
    },
    outdoor: {
      icon: '🌿',
      title: '户外活动',
      bgClass: 'outdoor',
      getDesc: (weather) => {
        return '空气清新，适合出行';
      },
      action: '查看推荐'
    }
  },

  // Get tips based on weather
  getTips(weatherCode, temp, weatherDesc) {
    const tips = [];
    const weatherStr = weatherDesc || '';

    // Rain-related tips
    if (weatherStr.includes('雨') || weatherStr.includes('雷')) {
      tips.push({
        ...this.tipTypes.umbrella,
        desc: this.tipTypes.umbrella.getDesc(weatherStr)
      });
    }

    // Temperature-related tips
    if (temp < 18) {
      tips.push({
        ...this.tipTypes.clothing,
        desc: this.tipTypes.clothing.getDesc(weatherStr, temp)
      });
    } else if (temp > 26 && !weatherStr.includes('雨')) {
      tips.push({
        ...this.tipTypes.sunscreen,
        desc: this.tipTypes.sunscreen.getDesc(weatherStr)
      });
    }

    // Dry weather
    if (weatherCode === 'sunny' && temp > 25) {
      tips.push({
        ...this.tipTypes.water,
        desc: this.tipTypes.water.getDesc(weatherStr, temp)
      });
    }

    // Good weather
    if (weatherCode === 'sunny' && temp >= 20 && temp <= 26) {
      tips.push({
        ...this.tipTypes.outdoor,
        desc: this.tipTypes.outdoor.getDesc(weatherStr)
      });
    }

    // If no tips matched, add outdoor tip
    if (tips.length === 0) {
      tips.push({
        ...this.tipTypes.outdoor,
        desc: this.tipTypes.outdoor.getDesc(weatherStr)
      });
    }

    return tips.slice(0, 4); // Max 4 tips
  }
};

// ============================================
// UI Controller
// ============================================

class WeatherApp {
  constructor() {
    this.currentCity = '北京';
    this.weatherData = null;
    this.currentOutfit = null;
    this.tempUnit = 'c';
    this.isLoading = true;

    this.init();
  }

  async init() {
    // Bind event listeners
    this.bindEvents();

    // Load initial data
    await this.loadWeatherData(this.currentCity);

    // Hide loading screen
    this.hideLoadingScreen();
  }

  bindEvents() {
    // City selector
    document.getElementById('location-btn').addEventListener('click', () => {
      this.openCitySheet();
    });

    // City sheet close
    document.getElementById('city-sheet-overlay').addEventListener('click', (e) => {
      if (e.target.id === 'city-sheet-overlay') {
        this.closeCitySheet();
      }
    });
    document.getElementById('city-sheet-close').addEventListener('click', () => {
      this.closeCitySheet();
    });

    // City search
    document.getElementById('city-search-input').addEventListener('input', (e) => {
      this.filterCities(e.target.value);
    });

    // City selection
    document.querySelectorAll('.city-item').forEach(item => {
      item.addEventListener('click', () => {
        const city = item.dataset.city;
        this.selectCity(city);
      });
    });

    // Settings
    document.getElementById('settings-btn').addEventListener('click', () => {
      this.openSettingsSheet();
    });

    document.getElementById('settings-sheet-overlay').addEventListener('click', (e) => {
      if (e.target.id === 'settings-sheet-overlay') {
        this.closeSettingsSheet();
      }
    });
    document.getElementById('settings-sheet-close').addEventListener('click', () => {
      this.closeSettingsSheet();
    });

    // Temperature unit toggle
    document.querySelectorAll('.settings-toggle').forEach(btn => {
      btn.addEventListener('click', () => {
        document.querySelectorAll('.settings-toggle').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        this.tempUnit = btn.dataset.unit;
        this.updateTemperatureDisplay();
      });
    });

    // Change outfit button
    document.getElementById('change-outfit-btn').addEventListener('click', () => {
      console.log('按钮点击');
      this.changeOutfit();
    });

    // Header scroll effect
    window.addEventListener('scroll', () => {
      const header = document.getElementById('app-header');
      if (window.scrollY > 10) {
        header.classList.add('scrolled');
      } else {
        header.classList.remove('scrolled');
      }
    });

    // Keyboard navigation
    document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape') {
        this.closeAllSheets();
      }
    });

    // Bottom sheet touch swipe
    this.bindSheetSwipe();
  }

  bindSheetSwipe() {
    const sheets = document.querySelectorAll('.bottom-sheet');
    sheets.forEach(sheet => {
      let startY = 0;
      let startTranslateY = 0;

      sheet.addEventListener('touchstart', (e) => {
        startY = e.touches[0].clientY;
        const transform = getComputedStyle(sheet).transform;
        if (transform !== 'none') {
          const matrix = new DOMMatrix(transform);
          startTranslateY = matrix.m42;
        }
      });

      sheet.addEventListener('touchmove', (e) => {
        const currentY = e.touches[0].clientY;
        const diff = currentY - startY;
        if (diff > 0) {
          sheet.style.transform = `translateX(-50%) translateY(${diff}px)`;
        }
      });

      sheet.addEventListener('touchend', (e) => {
        const endY = e.changedTouches[0].clientY;
        const diff = endY - startY;
        const velocity = Math.abs(diff) / 200;

        if (diff > 100 && velocity > 0.5) {
          this.closeAllSheets();
        } else {
          sheet.style.transform = 'translateX(-50%) translateY(0)';
        }
      });
    });
  }

  async loadWeatherData(city) {
    try {
      // Show loading state
      this.showLoadingState();

      // Get weather data
      const data = await WeatherAPI.getWeather(city);
      this.weatherData = data;

      // Update UI
      this.updateWeatherUI(data);
      this.updateForecastUI(data);
      this.updateOutfitUI(data);
      this.updateTipsUI(data);

      // Update city name in header
      document.getElementById('city-name').textContent = city;

      // Show success toast
      this.showToast(`${city}天气已更新`, 'success');

    } catch (error) {
      // Show error toast
      this.showToast(error.message || '获取天气失败', 'error');

      // Use default city
      if (city !== '北京') {
        await this.loadWeatherData('北京');
      }
    }
  }

  updateWeatherUI(data) {
    const { current } = data;

    // Temperature (with animation)
    const tempEl = document.getElementById('temp-number');
    this.animateNumber(tempEl, 0, current.temp, 1000);

    // Weather details
    document.getElementById('weather-desc').textContent = current.weather;
    document.getElementById('feels-like').textContent = `体感 ${current.feelsLike}°`;
    document.getElementById('humidity').textContent = `湿度 ${current.humidity}%`;
    document.getElementById('wind-speed').textContent = `风速 ${current.windSpeed}`;
    document.getElementById('update-time').textContent = `更新时间: ${current.updateTime}`;

    // Date
    const now = new Date();
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];
    const months = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'];

    document.getElementById('date-weekday').textContent = '今天';
    document.getElementById('date-full').textContent = `${months[now.getMonth()]}${now.getDate()}日`;

    // Weather animation
    this.updateWeatherAnimation(current.weatherCode);
  }

  updateForecastUI(data) {
    const { yesterday, today, tomorrow } = data;

    // Yesterday
    document.getElementById('icon-yesterday').textContent = yesterday.icon;
    document.getElementById('high-yesterday').textContent = `${yesterday.high}°`;
    document.getElementById('low-yesterday').textContent = `${yesterday.low}°`;

    // Today
    document.getElementById('icon-today').textContent = today.icon;
    document.getElementById('high-today').textContent = `${today.high}°`;
    document.getElementById('low-today').textContent = `${today.low}°`;

    // Tomorrow
    document.getElementById('icon-tomorrow').textContent = tomorrow.icon;
    document.getElementById('high-tomorrow').textContent = `${tomorrow.high}°`;
    document.getElementById('low-tomorrow').textContent = `${tomorrow.low}°`;
  }

  // Helper: set text content with <br> rendered as line breaks
  _setTextWithLineBreak(el, text) {
    el.innerHTML = '';
    const parts = text.split('<br>');
    parts.forEach((part, i) => {
      if (i > 0) el.appendChild(document.createElement('br'));
      el.appendChild(document.createTextNode(part));
    });
  }

  updateOutfitUI(data) {
    const { current } = data;
    const rec = OutfitEngine.getRecommendation(current.weatherCode, current.temp);
    this.currentOutfit = rec;

    // Update tags (使用 textContent 避免 XSS)
    const tagsEl = document.getElementById('outfit-tags');
    tagsEl.innerHTML = '';
    rec.tags.forEach(tag => {
      let tagClass = 'outfit-tag warm';
      if (tag.includes('雨')) tagClass = 'outfit-tag rain';
      else if (tag.includes('风')) tagClass = 'outfit-tag windy';
      const span = document.createElement('span');
      span.className = tagClass;
      span.textContent = tag;
      tagsEl.appendChild(span);
    });

    // Update outfit items (按 <br> 分割并渲染换行)
    this._setTextWithLineBreak(document.getElementById('outfit-top'), rec.top);
    this._setTextWithLineBreak(document.getElementById('outfit-bottom'), rec.bottom);
    this._setTextWithLineBreak(document.getElementById('outfit-shoes'), rec.shoes);
    this._setTextWithLineBreak(document.getElementById('outfit-accessories'), rec.accessories);

    // Update outfit icons (each recommendation has its own icons)
    if (rec.icons) {
      const iconEls = document.querySelectorAll('.outfit-card .outfit-icon');
      rec.icons.forEach((icon, i) => {
        if (iconEls[i]) iconEls[i].textContent = icon;
      });
    }

    // Update quote
    document.getElementById('outfit-quote').textContent = rec.quote;
  }

  // 换装时直接使用传入的 outfit，不重新查询
  updateOutfitUIWithOutfit(data, outfit) {
    const { current } = data;
    const rec = outfit;
    this.currentOutfit = rec;

    // Update tags
    const tagsEl = document.getElementById('outfit-tags');
    tagsEl.innerHTML = '';
    rec.tags.forEach(tag => {
      let tagClass = 'outfit-tag warm';
      if (tag.includes('雨')) tagClass = 'outfit-tag rain';
      else if (tag.includes('风')) tagClass = 'outfit-tag windy';
      const span = document.createElement('span');
      span.className = tagClass;
      span.textContent = tag;
      tagsEl.appendChild(span);
    });

    // Update outfit items
    this._setTextWithLineBreak(document.getElementById('outfit-top'), rec.top);
    this._setTextWithLineBreak(document.getElementById('outfit-bottom'), rec.bottom);
    this._setTextWithLineBreak(document.getElementById('outfit-shoes'), rec.shoes);
    this._setTextWithLineBreak(document.getElementById('outfit-accessories'), rec.accessories);

    // Update outfit icons
    if (rec.icons) {
      const iconEls = document.querySelectorAll('.outfit-card .outfit-icon');
      rec.icons.forEach((icon, i) => {
        if (iconEls[i]) iconEls[i].textContent = icon;
      });
    }

    // Update quote
    document.getElementById('outfit-quote').textContent = rec.quote;
  }

  updateTipsUI(data) {
    const { current } = data;
    const tips = TipsEngine.getTips(current.weatherCode, current.temp, current.weather);

    const tipsGrid = document.getElementById('tips-grid');
    tipsGrid.innerHTML = '';
    tips.forEach((tip, index) => {
      const card = document.createElement('div');
      card.className = 'tip-card';
      card.style.animationDelay = `${index * 100}ms`;

      const iconWrapper = document.createElement('div');
      iconWrapper.className = `tip-icon-wrapper ${tip.bgClass || ''}`;
      iconWrapper.textContent = tip.icon || '';
      card.appendChild(iconWrapper);

      const title = document.createElement('div');
      title.className = 'tip-title';
      title.textContent = tip.title || '';
      card.appendChild(title);

      const desc = document.createElement('div');
      desc.className = 'tip-desc';
      desc.textContent = tip.desc || '';
      card.appendChild(desc);

      const action = document.createElement('button');
      action.className = 'tip-action';
      action.textContent = tip.action || '';
      card.appendChild(action);

      tipsGrid.appendChild(card);
    });
  }

  updateWeatherAnimation(weatherCode) {
    const container = document.getElementById('weather-anim-container');
    container.className = 'weather-animation';

    // Update sun and clouds visibility based on weather
    const sun = container.querySelector('.sun');
    const sunRays = container.querySelector('.sun-rays');
    const clouds = container.querySelectorAll('.cloud');

    if (weatherCode === 'rainy' || weatherCode === 'heavyRain' || weatherCode === 'thunder') {
      // Rainy weather - show rain drops
      container.classList.add('rainy');
      if (sun) sun.style.display = 'none';
      if (sunRays) sunRays.style.display = 'none';
      clouds.forEach(c => c.style.display = 'none');

      // Add rain drops
      if (!container.querySelector('.rain-drop')) {
        for (let i = 0; i < 6; i++) {
          const drop = document.createElement('div');
          drop.className = 'rain-drop';
          container.appendChild(drop);
        }
      }
    } else if (weatherCode === 'cloudy' || weatherCode === 'overcast') {
      // Cloudy - show more clouds
      container.classList.add('cloudy');
      if (sun) sun.style.display = 'none';
      if (sunRays) sunRays.style.display = 'none';
      clouds.forEach(c => c.style.display = 'block');
    } else {
      // Sunny - default
      container.classList.add('sunny');
      if (sun) sun.style.display = 'block';
      if (sunRays) sunRays.style.display = 'block';
      clouds.forEach(c => c.style.display = 'block');
    }
  }

  updateTemperatureDisplay() {
    if (!this.weatherData) return;

    const celsius = this.weatherData.current.temp;
    let displayTemp = celsius;

    if (this.tempUnit === 'f') {
      displayTemp = Math.round(celsius * 9 / 5 + 32);
    }

    document.getElementById('temp-number').textContent = displayTemp;
  }

  animateNumber(element, start, end, duration) {
    const startTime = performance.now();
    const diff = end - start;

    const animate = (currentTime) => {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);

      // Easing function
      const eased = 1 - Math.pow(1 - progress, 3);
      const current = Math.round(start + diff * eased);

      element.textContent = current;

      if (progress < 1) {
        requestAnimationFrame(animate);
      }
    };

    requestAnimationFrame(animate);
  }

  changeOutfit() {
    console.log('changeOutfit called', {
      hasWeatherData: !!this.weatherData,
      hasCurrentOutfit: !!this.currentOutfit,
      animating: document.getElementById('change-outfit-btn').dataset.animating
    });
    if (!this.weatherData || !this.currentOutfit) return;

    // Guard: prevent re-entrancy during animation
    const btn = document.getElementById('change-outfit-btn');
    if (btn.dataset.animating === 'true') return;
    btn.disabled = true;
    btn.dataset.animating = 'true';

    const alt = OutfitEngine.getAlternative(this.currentOutfit, this.weatherData.current.weatherCode);
    this.currentOutfit = alt;

    // Animate transition
    const outfitSection = document.getElementById('outfit-section');
    outfitSection.style.opacity = '0';
    outfitSection.style.transform = 'translateY(10px)';

    setTimeout(() => {
      console.log('setTimeout callback firing');
      this.updateOutfitUIWithOutfit(this.weatherData, alt);
      outfitSection.style.opacity = '1';
      outfitSection.style.transform = 'translateY(0)';
      // Re-enable button after animation completes
      btn.disabled = false;
      btn.dataset.animating = 'false';
    }, 200);
  }

  showLoadingState() {
    // Could add skeleton loading here
  }

  showLoadingScreen() {
    const loader = document.getElementById('loading-screen');
    loader.classList.remove('hidden');
    loader.setAttribute('aria-hidden', 'false');
  }

  hideLoadingScreen() {
    setTimeout(() => {
      const loader = document.getElementById('loading-screen');
      loader.classList.add('hidden');
      loader.setAttribute('aria-hidden', 'true');
    }, 800);
  }

  // City Sheet Methods
  openCitySheet() {
    const overlay = document.getElementById('city-sheet-overlay');
    const arrow = document.getElementById('dropdown-arrow');
    overlay.classList.add('active');
    arrow.classList.add('rotated');
    document.body.style.overflow = 'hidden';

    // Highlight current city
    document.querySelectorAll('.city-item').forEach(item => {
      item.classList.toggle('selected', item.dataset.city === this.currentCity);
    });
  }

  closeCitySheet() {
    const overlay = document.getElementById('city-sheet-overlay');
    const arrow = document.getElementById('dropdown-arrow');
    overlay.classList.remove('active');
    arrow.classList.remove('rotated');
    document.body.style.overflow = '';
  }

  filterCities(query) {
    const items = document.querySelectorAll('.city-item');
    const lowerQuery = query.toLowerCase();

    items.forEach(item => {
      const cityName = item.dataset.city.toLowerCase();
      const province = item.querySelector('.city-item-province').textContent.toLowerCase();
      const match = cityName.includes(lowerQuery) || province.includes(lowerQuery);
      item.style.display = match ? 'flex' : 'none';
    });
  }

  async selectCity(city) {
    if (city === this.currentCity) {
      this.closeCitySheet();
      return;
    }

    this.currentCity = city;

    // Close sheet with animation
    this.closeCitySheet();

    // Show loading briefly
    await new Promise(resolve => setTimeout(resolve, 300));

    // Load new weather
    await this.loadWeatherData(city);
  }

  // Settings Sheet Methods
  openSettingsSheet() {
    const overlay = document.getElementById('settings-sheet-overlay');
    overlay.classList.add('active');
    document.body.style.overflow = 'hidden';
  }

  closeSettingsSheet() {
    const overlay = document.getElementById('settings-sheet-overlay');
    overlay.classList.remove('active');
    document.body.style.overflow = '';
  }

  closeAllSheets() {
    this.closeCitySheet();
    this.closeSettingsSheet();
  }

  // Toast Notifications
  showToast(message, type = 'info', duration = 3000) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;

    const icons = {
      success: '✓',
      warning: '⚠',
      error: '✕',
      info: 'ℹ'
    };

    toast.innerHTML = `
      <span class="toast-icon">${icons[type]}</span>
      <span class="toast-message">${message}</span>
      <button class="toast-close" aria-label="关闭">×</button>
    `;

    container.appendChild(toast);

    // Close button
    toast.querySelector('.toast-close').addEventListener('click', () => {
      this.hideToast(toast);
    });

    // Auto hide
    setTimeout(() => {
      this.hideToast(toast);
    }, duration);
  }

  hideToast(toast) {
    if (!toast || toast.classList.contains('hiding')) return;
    toast.classList.add('hiding');
    setTimeout(() => {
      if (toast.parentNode) {
        toast.parentNode.removeChild(toast);
      }
    }, 300);
  }
}

// ============================================
// Initialize App
// ============================================

document.addEventListener('DOMContentLoaded', () => {
  window.weatherApp = new WeatherApp();
});

// Apply CSS transitions for outfit section
const outfitSection = document.getElementById('outfit-section');
if (outfitSection) {
  outfitSection.style.transition = 'opacity 0.2s ease, transform 0.2s ease';
}
