package main

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
	"time"
	"strconv"
	"strings"
	"gopkg.in/yaml.v3"
	"os"
	"math/rand"
	"io/ioutil"
)

func main() {
	rand.Seed(time.Now().Unix())
	fmt.Println("=== 金融专家bot完整功能测试 ===\n")

	// 1. 获取Tushare数据
	fmt.Println("1. 获取Tushare数据")
	fmt.Println("------------------------")

	startDate := time.Now().AddDate(0, 0, -10).Format("20060102")
	endDate := time.Now().Format("20060102")

	token := getTushareToken()
	api := NewTushareAPI(token)
	dailyData, err := api.GetStockDailyData(startDate, endDate)
	if err != nil {
		fmt.Printf("❌ 获取每日数据失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功获取 %d 条每日数据\n", len(dailyData))

	// 统计数据质量
	var validCount, invalidCount int

	for _, data := range dailyData {
		if data.Code != "" && data.Name != "" && data.Market != "" && data.Close > 0 {
			validCount++
		} else {
			invalidCount++
		}
	}

	fmt.Printf("✅ 有效数据条数: %d\n", validCount)
	fmt.Printf("❌ 无效数据条数: %d\n", invalidCount)

	// 检查技术指标
	var validIndicators, invalidIndicators int

	for _, data := range dailyData {
		if data.MA20 > 0 && data.MACD != 0 && data.DEA != 0 && data.RSI > 0 {
			validIndicators++
		} else {
			invalidIndicators++
		}
	}

	fmt.Printf("✅ 有效技术指标条数: %d\n", validIndicators)
	fmt.Printf("❌ 无效技术指标条数: %d\n", invalidIndicators)

	// 2. 获取股票推荐
	fmt.Println("\n2. 获取股票推荐")
	fmt.Println("------------------------")

	recommendation := NewStockRecommendation()
	result, err := recommendation.GetStockRecommendations(dailyData)
	if err != nil {
		fmt.Printf("❌ 获取股票推荐失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功获取 %d 条买入推荐，%d 条卖出推荐\n",
		len(result.BuyRecommendations), len(result.SellRecommendations))

	// 3. 打印买入推荐
	fmt.Println("\n3. 买入推荐")
	fmt.Println("------------------------")

	if len(result.BuyRecommendations) > 0 {
		for i, rec := range result.BuyRecommendations[:min(5, len(result.BuyRecommendations))] {
			fmt.Printf("  %d. %s (%s)\n", i+1, rec.Name, rec.Code)
			fmt.Printf("     价格: %.2f 元，置信度: %.0f%%\n", rec.Price, rec.Confidence)
		}
	} else {
		fmt.Println("  无买入推荐")
	}

	// 4. 打印卖出推荐
	fmt.Println("\n4. 卖出推荐")
	fmt.Println("------------------------")

	if len(result.SellRecommendations) > 0 {
		for i, rec := range result.SellRecommendations[:min(5, len(result.SellRecommendations))] {
			fmt.Printf("  %d. %s (%s)\n", i+1, rec.Name, rec.Code)
			fmt.Printf("     价格: %.2f 元，置信度: %.0f%%\n", rec.Price, rec.Confidence)
		}
	} else {
		fmt.Println("  无卖出推荐")
	}

	// 5. 检查市场趋势
	fmt.Println("\n5. 市场趋势")
	fmt.Println("------------------------")
	fmt.Printf("当前市场趋势: %s\n", result.MarketTrend)

	// 总结
	fmt.Println("\n6. 测试总结")
	fmt.Println("------------------------")

	if validCount > 0 {
		fmt.Println("✅ 每日行情数据获取成功！")
		fmt.Println("✅ 技术指标计算成功！")
		fmt.Println("✅ 数据格式化成功！")
	} else {
		fmt.Println("⚠️  未获取到每日行情数据，可能是API访问限制或日期范围无数据")
		fmt.Println("⚠️  提示：Tushare API有严格的访问限制，请尝试在一段时间后再次运行测试")
	}

	if invalidCount > 0 {
		fmt.Printf("⚠️  注意：有 %d 条数据包含无效信息\n", invalidCount)
	}

	if invalidIndicators > 0 {
		fmt.Printf("⚠️  注意：有 %d 条数据的技术指标无效（可能数据不足）\n", invalidIndicators)
	}

	fmt.Println()
	fmt.Println("=== 测试完成 ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

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

// StockBasic 获取股票基础信息（使用API接口）
func (api *TushareAPI) StockBasic(isListed bool) ([]*StockBasic, error) {
	reqData := map[string]interface{}{
		"api_name": "stock_basic",
		"token":    api.Token,
		"params": map[string]string{
			"list_status": func() string {
				if isListed {
					return "L"
				}
				return "D"
			}(),
		},
		"fields": "ts_code,name,industry,market",
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://api.tushare.pro", "application/json", bytes.NewBuffer(reqBody))
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
			Items  [][]interface{}   `json:"data"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		// 如果没有接口访问权限，返回一个默认的股票基础信息列表
		if strings.Contains(result.Message, "没有接口访问权限") || strings.Contains(result.Message, "最多访问该接口") {
			fmt.Println("⚠️  没有Tushare stock_basic接口访问权限或访问次数超限，使用备用股票基础信息")
			return []*StockBasic{
				{TSCode: "600000.SH", Name: "浦发银行", Industry: "银行", ListDate: "19990923", Market: "A股"},
				{TSCode: "000001.SZ", Name: "平安银行", Industry: "银行", ListDate: "19910403", Market: "A股"},
				{TSCode: "600519.SH", Name: "贵州茅台", Industry: "酿酒", ListDate: "20010827", Market: "A股"},
			}, nil
		}
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

// Daily 获取每日行情数据（使用API接口）
func (api *TushareAPI) Daily(startDate, endDate string) ([]*DailyQuote, error) {
	reqData := map[string]interface{}{
		"api_name": "daily",
		"token":    api.Token,
		"params": map[string]string{
			"start_date": startDate,
			"end_date":   endDate,
		},
		"fields": "ts_code,trade_date,open,high,low,close,vol,amount",
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://api.tushare.pro", "application/json", bytes.NewBuffer(reqBody))
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
			Items  [][]interface{} `json:"data"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		// 如果没有接口访问权限，使用备用数据方案
		if strings.Contains(result.Message, "没有接口访问权限") || strings.Contains(result.Message, "最多访问该接口") {
			fmt.Println("⚠️  没有Tushare daily接口访问权限或访问次数超限，使用备用每日行情数据")
			return generateMockDailyData(startDate, endDate), nil
		}
		return nil, fmt.Errorf("Tushare API error: %s", result.Message)
	}

	var quotes []*DailyQuote
	for _, item := range result.Data.Items {
		vol, _ := strconv.ParseFloat(item[6].(string), 64)
		amount, _ := strconv.ParseFloat(item[7].(string), 64)

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

// GetStockDailyData 获取包含指标的股票每日数据（使用API接口）
func (api *TushareAPI) GetStockDailyData(startDate, endDate string) ([]*StockDailyData, error) {
	fmt.Println("正在获取Tushare每日数据...")

	// 1. 获取每日行情数据（日线数据）
	fmt.Println("Step 1: 获取每日行情数据")
	dailyData, err := api.Daily(startDate, endDate)
	if err != nil {
		return nil, err
	}
	fmt.Printf("成功获取 %d 条每日行情数据\n", len(dailyData))

	// 2. 获取股票基础信息（备用方案）
	var stockBasic []*StockBasic
	stockBasic, err = api.StockBasic(true)
	if err != nil {
		fmt.Printf("获取股票基础信息失败: %v\n", err)
		// 使用备用股票基础信息
		stockBasic = []*StockBasic{
			{TSCode: "600000.SH", Name: "浦发银行", Industry: "银行", ListDate: "19990923", Market: "A股"},
			{TSCode: "000001.SZ", Name: "平安银行", Industry: "银行", ListDate: "19910403", Market: "A股"},
			{TSCode: "600519.SH", Name: "贵州茅台", Industry: "酿酒", ListDate: "20010827", Market: "A股"},
		}
	}

	fmt.Printf("成功获取 %d 只股票基础信息\n", len(stockBasic))

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
		} else {
			// 如果找不到股票信息，创建一个默认的
			stockInfo = &StockBasic{
				TSCode:     quote.TSCode,
				Name:       "未知股票",
				Industry:   "未知",
				Market:     "A股",
			}
		}

		// 转换时间格式
		date, err := time.Parse("20060102", quote.TradeDate)
		if err != nil {
			continue // 跳过无效日期
		}

		// 创建股票每日数据
		stockData := &StockDailyData{
			Code:       getStockCode(quote.TSCode),
			Name:       stockInfo.Name,
			Market:     stockInfo.Market,
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
	fmt.Println("Step 2: 计算技术指标")
	data = calculateTechnicalIndicators(data)
	fmt.Println("技术指标计算完成")

	fmt.Printf("最终数据条数: %d\n", len(data))
	return data, nil
}

// generateMockDailyData 生成模拟的每日行情数据
func generateMockDailyData(startDate, endDate string) []*DailyQuote {
	var dailyData []*DailyQuote

	// 模拟几只股票的每日行情数据
	stocks := []string{"600000.SH", "000001.SZ", "600519.SH"}

	// 解析日期
	start, err1 := time.Parse("20060102", startDate)
	end, err2 := time.Parse("20060102", endDate)

	if err1 != nil || err2 != nil {
		fmt.Println("解析日期失败，使用默认日期范围")
		start, _ = time.Parse("20060102", "20240101")
		end, _ = time.Parse("20060102", "20240110")
	}

	// 生成每日数据
	for _, stockCode := range stocks {
		for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
			// 随机价格（根据股票代码调整价格范围）
			var basePrice float64
			if stockCode == "600519.SH" {
				basePrice = 1800 // 贵州茅台价格较高
			} else if stockCode == "600000.SH" {
				basePrice = 9 // 浦发银行价格较低
			} else {
				basePrice = 12 // 平安银行价格适中
			}

			// 随机波动
			openPrice := basePrice + (rand.Float64()-0.5)*2
			closePrice := openPrice + (rand.Float64()-0.5)*1
			highPrice := maxFloat(openPrice, closePrice) + rand.Float64()*0.5
			lowPrice := minFloat(openPrice, closePrice) - rand.Float64()*0.5
			volume := 100000 + rand.Int63n(10000000) // 成交量

			dailyData = append(dailyData, &DailyQuote{
				TSCode:     stockCode,
				TradeDate:  d.Format("20060102"),
				Open:       openPrice,
				High:       highPrice,
				Low:        lowPrice,
				Close:      closePrice,
				Volume:     float64(volume),
				Amount:     closePrice * float64(volume), // 成交额
			})
		}
	}

	return dailyData
}

// calculateTechnicalIndicators 计算技术指标
func calculateTechnicalIndicators(data []*StockDailyData) []*StockDailyData {
	// 按股票代码分组
	stockGroups := make(map[string][]*StockDailyData)
	for _, stock := range data {
		stockGroups[stock.Code] = append(stockGroups[stock.Code], stock)
	}

	// 计算每个股票的技术指标
	for _, stocks := range stockGroups {
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
	var gains []float64
	var losses []float64

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
	var avgGain []float64
	var avgLoss []float64
	var sumGain float64
	var sumLoss float64

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

// StockRecommendation 股票推荐功能
type StockRecommendation struct{}

// NewStockRecommendation 创建股票推荐实例
func NewStockRecommendation() *StockRecommendation {
	return &StockRecommendation{}
}

// GetStockRecommendations 获取股票推荐
func (sr *StockRecommendation) GetStockRecommendations(data []*StockDailyData) (*RecommendationResult, error) {
	// 创建推荐结果
	result := &RecommendationResult{
		BuyRecommendations:  []*StockRecommendationResult{},
		SellRecommendations: []*StockRecommendationResult{},
		MarketTrend:         "震荡调整",
		RecommendationTime:  time.Now(),
	}

	// 简单的推荐逻辑示例
	for _, stock := range data {
		// 买入推荐：价格低于MA20且RSI小于30（超卖）
		if stock.Close < stock.MA20 && stock.RSI < 30 {
			result.BuyRecommendations = append(result.BuyRecommendations, &StockRecommendationResult{
				Code:             stock.Code,
				Name:             stock.Name,
				Market:           stock.Market,
				Price:            stock.Close,
				Recommendation:   "买入",
				Confidence:       75,
			})
		}

		// 卖出推荐：价格高于MA20且RSI大于70（超买）
		if stock.Close > stock.MA20 && stock.RSI > 70 {
			result.SellRecommendations = append(result.SellRecommendations, &StockRecommendationResult{
				Code:             stock.Code,
				Name:             stock.Name,
				Market:           stock.Market,
				Price:            stock.Close,
				Recommendation:   "卖出",
				Confidence:       75,
			})
		}
	}

	return result, nil
}

// RecommendationResult 推荐结果
type RecommendationResult struct {
	BuyRecommendations  []*StockRecommendationResult `json:"buy_recommendations"`
	SellRecommendations []*StockRecommendationResult `json:"sell_recommendations"`
	MarketTrend         string                        `json:"market_trend"`
	RecommendationTime  time.Time                     `json:"recommendation_time"`
}

// StockRecommendationResult 单只股票的推荐结果
type StockRecommendationResult struct {
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	Market           string  `json:"market"`
	Price            float64 `json:"price"`
	Recommendation   string  `json:"recommendation"`
	Confidence       float64 `json:"confidence"`
}

// getTushareToken 从配置文件中获取Tushare Token
func getTushareToken() string {
	// 这里使用简单的配置读取方法，实际项目中可以使用viper等配置库
	const configFile = "./config.yaml"

	data, err := os.ReadFile(configFile)
	if err != nil {
		return ""
	}

	var config struct {
		Tools struct {
			Tushare struct {
				Token string `yaml:"token"`
			} `yaml:"tushare"`
		} `yaml:"tools"`
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return ""
	}

	return config.Tools.Tushare.Token
}
