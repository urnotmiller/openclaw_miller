package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// TushareAPI Tushare API配置
type TushareAPI struct {
	Token string
}

// NewTushareAPI 创建Tushare API实例
func NewTushareAPI(token string) *TushareAPI {
	return &TushareAPI{
		Token: token,
	}
}

// StockBasic 获取股票基础信息
func (api *TushareAPI) StockBasic(isListed bool) ([]*StockBasic, error) {
	params := url.Values{}
	params.Set("api_name", "stock_basic")
	params.Set("token", api.Token)
	params.Set("list_status", func() string {
		if isListed {
			return "L"
		}
		return "D"
	}())

	resp, err := http.Get("https://api.tushare.pro?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
		Data    struct {
			Fields []string          `json:"fields"`
			Items  [][]interface{}   `json:"items"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("Tushare API error: %s", result.Message)
	}

	var stocks []*StockBasic
	for _, item := range result.Data.Items {
		stock := &StockBasic{
			TSCode:     item[0].(string),
			Name:       item[1].(string),
			Industry:   item[2].(string),
			ListDate:   item[3].(string),
			Market:     item[4].(string),
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// Daily 获取每日行情数据
func (api *TushareAPI) Daily(startDate, endDate string) ([]*DailyQuote, error) {
	params := url.Values{}
	params.Set("api_name", "daily")
	params.Set("token", api.Token)
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)

	resp, err := http.Get("https://api.tushare.pro?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
		Data    struct {
			Fields []string        `json:"fields"`
			Items  [][]interface{} `json:"items"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("Tushare API error: %s", result.Message)
	}

	var quotes []*DailyQuote
	for _, item := range result.Data.Items {
		vol, _ := strconv.ParseFloat(item[5].(string), 64)
		amount, _ := strconv.ParseFloat(item[6].(string), 64)

		quote := &DailyQuote{
			TSCode:     item[0].(string),
			TradeDate:  item[1].(string),
			Open:       item[2].(float64),
			High:       item[3].(float64),
			Low:        item[4].(float64),
			Close:      item[5].(float64),
			Volume:     vol,
			Amount:     amount,
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

// GetStockDailyData 获取股票每日数据（包含技术指标）
func (api *TushareAPI) GetStockDailyData(startDate, endDate string) ([]*StockDailyData, error) {
	// 1. 获取每日行情数据
	dailyData, err := api.Daily(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 2. 获取股票基础信息
	stockBasic, err := api.StockBasic(true)
	if err != nil {
		return nil, err
	}

	// 3. 创建股票代码到名称和市场的映射
	stockMap := make(map[string]*StockBasic)
	for _, stock := range stockBasic {
		stockMap[stock.TSCode] = stock
	}

	// 4. 转换数据格式
	var data []*StockDailyData
	for _, quote := range dailyData {
		// 获取股票信息
		var stockInfo *StockBasic
		if info, ok := stockMap[quote.TSCode]; ok {
			stockInfo = info
		}

		// 转换时间格式
		date, err := time.Parse("20060102", quote.TradeDate)
		if err != nil {
			continue // 跳过无效日期
		}

		// 创建股票每日数据
		stockData := &StockDailyData{
			Code:       getStockCode(quote.TSCode),
			Name:       func() string { if stockInfo != nil { return stockInfo.Name } else { return "" } }(),
			Market:     func() string { if stockInfo != nil { return stockInfo.Market } else { return "A股" } }(),
			Date:       date,
			Close:      quote.Close,
			Open:       quote.Open,
			High:       quote.High,
			Low:        quote.Low,
			Volume:     quote.Volume,
			AvgVolume:  0, // 待计算平均成交量
		}

		data = append(data, stockData)
	}

	// 5. 计算平均成交量
	volumeMap := make(map[string][]float64)
	for _, stock := range data {
		volumeMap[stock.Code] = append(volumeMap[stock.Code], stock.Volume)
	}

	// 计算每个股票的平均成交量
	for _, stock := range data {
		if volumes, ok := volumeMap[stock.Code]; ok {
			var sum float64
			for _, vol := range volumes {
				sum += vol
			}
			stock.AvgVolume = sum / float64(len(volumes))
		}
	}

	// 6. 计算技术指标
	data = calculateTechnicalIndicators(data)

	return data, nil
}

// calculateTechnicalIndicators 计算技术指标
func calculateTechnicalIndicators(data []*StockDailyData) []*StockDailyData {
	// 按股票代码分组
	stockGroups := make(map[string][]*StockDailyData)
	for _, stock := range data {
		stockGroups[stock.Code] = append(stockGroups[stock.Code], stock)
	}

	// 计算每个股票的技术指标
	for code, stocks := range stockGroups {
		// 计算MA20
		ma20 := calculateMA(stocks, 20)

		// 计算MACD
		macd, dea := calculateMACD(stocks)

		// 计算RSI
		rsi := calculateRSI(stocks)

		// 将计算结果赋值给对应的数据点
		for i, stock := range stocks {
			if i >= 19 { // MA20需要至少20个数据点
				stock.MA20 = ma20[i-19]
			}

			if i >= 25 { // MACD需要至少26个数据点
				stock.MACD = macd[i-25]
				stock.DEA = dea[i-25]
			}

			if i >= 13 { // RSI需要至少14个数据点
				stock.RSI = rsi[i-13]
			}
		}
	}

	return data
}

// calculateMA 计算移动平均线
func calculateMA(data []*StockDailyData, period int) []float64 {
	var ma []float64

	for i := period - 1; i < len(data); i++ {
		var sum float64
		for j := 0; j < period; j++ {
			sum += data[i-j].Close
		}
		ma = append(ma, sum/float64(period))
	}

	return ma
}

// calculateMACD 计算MACD指标
func calculateMACD(data []*StockDailyData) ([]float64, []float64) {
	var macd, dea []float64

	// 计算EMA12和EMA26
	ema12 := calculateEMA(data, 12)
	ema26 := calculateEMA(data, 26)

	// 计算DIF（EMA12 - EMA26）
	var dif []float64
	for i, e12 := range ema12 {
		if i >= 14 { // EMA26需要至少26个数据点，EMA12需要12个
			dif = append(dif, e12-ema26[i-14])
		}
	}

	// 计算DEA（DIF的EMA9）
	var deaData []*StockDailyData
	for _, d := range dif {
		deaData = append(deaData, &StockDailyData{Close: d})
	}

	dea = calculateEMA(deaData, 9)

	// 计算MACD线
	for i, d := range dif {
		if i >= 8 { // DEA需要至少9个数据点
			macd = append(macd, d-dea[i-8])
		}
	}

	return macd, dea
}

// calculateEMA 计算指数移动平均线
func calculateEMA(data []*StockDailyData, period int) []float64 {
	var ema []float64
	multiplier := 2.0 / float64(period+1)

	// 初始值为简单移动平均
	var initialSum float64
	for i := 0; i < period; i++ {
		initialSum += data[i].Close
	}
	ema = append(ema, initialSum/float64(period))

	// 计算后续EMA值
	for i := period; i < len(data); i++ {
		ema = append(ema, data[i].Close*multiplier+ema[i-period]*(1-multiplier))
	}

	return ema
}

// calculateRSI 计算RSI指标
func calculateRSI(data []*StockDailyData, period ...int) []float64 {
	rsiPeriod := 14
	if len(period) > 0 {
		rsiPeriod = period[0]
	}

	var rsi []float64
	var gains, losses []float64

	// 计算价格变化
	for i := 1; i < len(data); i++ {
		change := data[i].Close - data[i-1].Close
		if change > 0 {
			gains = append(gains, change)
			losses = append(losses, 0)
		} else if change < 0 {
			gains = append(gains, 0)
			losses = append(losses, -change)
		} else {
			gains = append(gains, 0)
			losses = append(losses, 0)
		}
	}

	// 计算平均增益和平均损失
	var avgGain, avgLoss []float64
	var sumGain, sumLoss float64

	// 初始平均
	for i := 0; i < rsiPeriod; i++ {
		sumGain += gains[i]
		sumLoss += losses[i]
	}
	avgGain = append(avgGain, sumGain/float64(rsiPeriod))
	avgLoss = append(avgLoss, sumLoss/float64(rsiPeriod))

	// 后续平均
	for i := rsiPeriod; i < len(gains); i++ {
		avgGain = append(avgGain, (avgGain[i-rsiPeriod]*(float64(rsiPeriod)-1)+gains[i])/float64(rsiPeriod))
		avgLoss = append(avgLoss, (avgLoss[i-rsiPeriod]*(float64(rsiPeriod)-1)+losses[i])/float64(rsiPeriod))
	}

	// 计算RSI
	for i, gain := range avgGain {
		loss := avgLoss[i]
		if loss == 0 {
			rsi = append(rsi, 100)
		} else {
			rs := gain / loss
			rsi = append(rsi, 100-(100/(1+rs)))
		}
	}

	return rsi
}

// getStockCode 从TSCode中提取股票代码
func getStockCode(tsCode string) string {
	// TSCode格式为 "600000.SH" 或 "000001.SZ"
	parts := strings.Split(tsCode, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return tsCode
}

// StockBasic 股票基础信息
type StockBasic struct {
	TSCode     string
	Name       string
	Industry   string
	ListDate   string
	Market     string
}

// DailyQuote 每日行情
type DailyQuote struct {
	TSCode     string
	TradeDate  string
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     float64
	Amount     float64
}

// StockDailyData 包含指标的股票每日数据
type StockDailyData struct {
	Code       string
	Name       string
	Market     string
	Date       time.Time
	Close      float64
	Open       float64
	High       float64
	Low        float64
	Volume     float64
	AvgVolume  float64
	MA20       float64 // 20日均线
	MACD       float64 // MACD线
	DEA        float64 // MACD信号线
	RSI        float64 // RSI指标
}
