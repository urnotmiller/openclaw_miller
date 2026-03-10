package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// StockDataAlternative 备用股票数据获取方案
type StockDataAlternative struct{}

// NewStockDataAlternative 创建备用数据获取实例
func NewStockDataAlternative() *StockDataAlternative {
	return &StockDataAlternative{}
}

// GetStockBasicInfo 获取股票基础信息（使用备用数据源）
func (a *StockDataAlternative) GetStockBasicInfo() ([]*StockBasic, error) {
	// 使用模拟数据作为备用方案
	// 实际应用中可以使用其他免费API或数据库
	return []*StockBasic{
		{TSCode: "600000.SH", Name: "浦发银行", Industry: "银行", Market: "主板"},
		{TSCode: "600519.SH", Name: "贵州茅台", Industry: "饮料制造", Market: "主板"},
		{TSCode: "000001.SZ", Name: "平安银行", Industry: "银行", Market: "主板"},
		{TSCode: "000858.SZ", Name: "五粮液", Industry: "饮料制造", Market: "主板"},
		{TSCode: "000002.SZ", Name: "万科A", Industry: "房地产开发", Market: "主板"},
	}, nil
}

// GetStockDailyData 获取股票每日数据（使用备用数据源）
func (a *StockDataAlternative) GetStockDailyData(startDate, endDate string) ([]*StockDailyData, error) {
	// 使用模拟数据作为备用方案
	var data []*StockDailyData

	// 模拟股票代码列表
	stocks := []*StockBasic{
		{TSCode: "600000.SH", Name: "浦发银行", Industry: "银行", Market: "主板"},
		{TSCode: "600519.SH", Name: "贵州茅台", Industry: "饮料制造", Market: "主板"},
		{TSCode: "000001.SZ", Name: "平安银行", Industry: "银行", Market: "主板"},
	}

	// 生成模拟数据
	for _, stock := range stocks {
		start, err := time.Parse("20060102", startDate)
		if err != nil {
			return nil, err
		}

		end, err := time.Parse("20060102", endDate)
		if err != nil {
			return nil, err
		}

		// 计算日期范围
		days := int(end.Sub(start).Hours() / 24) + 1

		// 生成每日数据
		for i := 0; i < days; i++ {
			date := start.AddDate(0, 0, i)

			// 模拟价格波动
			basePrice := 100.0
			if strings.Contains(stock.TSCode, "600519") {
				basePrice = 2000.0
			}

			change := (rand.Float64() - 0.5) * 20
			price := basePrice + change

			// 模拟成交量
			volume := 100000 + rand.Intn(1000000)

			// 模拟指标数据
			ma20 := basePrice + (rand.Float64() - 0.5) * 10
			macd := (rand.Float64() - 0.5) * 2
			dea := (rand.Float64() - 0.5) * 2
			rsi := 40 + rand.Float64() * 20

			data = append(data, &StockDailyData{
				Code:       strings.Split(stock.TSCode, ".")[0],
				Name:       stock.Name,
				Market:     stock.Market,
				Date:       date,
				Close:      price,
				Open:       price - 2 + rand.Float64()*4,
				High:       price + 3,
				Low:        price - 3,
				Volume:     float64(volume),
				AvgVolume:  500000,
				MA20:       ma20,
				MACD:       macd,
				DEA:        dea,
				RSI:        rsi,
			})
		}
	}

	return data, nil
}

// GetRealTimeQuote 获取实时行情（使用腾讯财经API）
func (a *StockDataAlternative) GetRealTimeQuote(market, code string) (*RealTimeQuote, error) {
	// 使用腾讯财经API作为实时数据来源
	api := NewQQFinanceAPI()
	return api.GetStockQuote(market, code)
}

// GetStockRecommendations 获取股票推荐（使用模拟数据）
func (a *StockDataAlternative) GetStockRecommendations() (*RecommendationResult, error) {
	// 使用模拟数据生成推荐
	buyRecommendations := []*StockRecommendationResult{
		{Code: "600000", Name: "浦发银行", Market: "A股", Price: 9.85, Recommendation: "买入", Confidence: 85},
		{Code: "600519", Name: "贵州茅台", Market: "A股", Price: 1850.00, Recommendation: "买入", Confidence: 78},
		{Code: "000001", Name: "平安银行", Market: "A股", Price: 12.50, Recommendation: "买入", Confidence: 82},
	}

	sellRecommendations := []*StockRecommendationResult{
		{Code: "600837", Name: "海通证券", Market: "A股", Price: 13.50, Recommendation: "卖出", Confidence: 75},
		{Code: "000776", Name: "广发证券", Market: "A股", Price: 18.20, Recommendation: "卖出", Confidence: 70},
	}

	return &RecommendationResult{
		BuyRecommendations:  buyRecommendations,
		SellRecommendations: sellRecommendations,
		MarketTrend:         "震荡调整",
		RecommendationTime:  time.Now(),
	}, nil
}

// StockBasic 股票基础信息
type StockBasic struct {
	TSCode     string
	Name       string
	Industry   string
	ListDate   string
	Market     string
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

// RealTimeQuote 实时行情
type RealTimeQuote struct {
	Code          string
	Name          string
	Price         float64
	PreviousClose float64
	Open          float64
	High          float64
	Low           float64
	Volume        float64
	Amount        float64
	Change        float64
	ChangeRatio   float64
	Time          time.Time
}

// StockRecommendationResult 股票推荐结果
type StockRecommendationResult struct {
	Code             string
	Name             string
	Market           string
	Price            float64
	Recommendation   string // "买入"或"卖出"
	Confidence       float64 // 置信度（0-100）
}

// RecommendationResult 推荐结果汇总
type RecommendationResult struct {
	BuyRecommendations  []*StockRecommendationResult
	SellRecommendations []*StockRecommendationResult
	MarketTrend         string
	RecommendationTime  time.Time
}
