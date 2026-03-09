package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// StockRecommendation 股票推荐功能
type StockRecommendation struct{}

// NewStockRecommendation 创建股票推荐实例
func NewStockRecommendation() *StockRecommendation {
	return &StockRecommendation{}
}

// RecommendBuyStocks 推荐买入股票
func (sr *StockRecommendation) RecommendBuyStocks(data []*StockDailyData) ([]*StockRecommendationResult, error) {
	var recommendations []*StockRecommendationResult

	for _, stock := range data {
		if sr.shouldBuy(stock) {
			score := sr.calculateBuyScore(stock)
			recommendations = append(recommendations, &StockRecommendationResult{
				Code:          stock.Code,
				Name:          stock.Name,
				Market:        stock.Market,
				Price:         stock.Close,
				Recommendation: "买入",
				Confidence:    score,
			})
		}
	}

	// 按置信度排序
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Confidence > recommendations[j].Confidence
	})

	return recommendations, nil
}

// RecommendSellStocks 推荐卖出股票
func (sr *StockRecommendation) RecommendSellStocks(data []*StockDailyData) ([]*StockRecommendationResult, error) {
	var recommendations []*StockRecommendationResult

	for _, stock := range data {
		if sr.shouldSell(stock) {
			score := sr.calculateSellScore(stock)
			recommendations = append(recommendations, &StockRecommendationResult{
				Code:          stock.Code,
				Name:          stock.Name,
				Market:        stock.Market,
				Price:         stock.Close,
				Recommendation: "卖出",
				Confidence:    score,
			})
		}
	}

	// 按置信度排序
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Confidence > recommendations[j].Confidence
	})

	return recommendations, nil
}

// shouldBuy 判断是否应该买入
func (sr *StockRecommendation) shouldBuy(stock *StockDailyData) bool {
	// 价格突破MA20均线
	if stock.Close <= stock.MA20 {
		return false
	}

	// MACD金叉
	if stock.MACD < 0 || stock.MACD < stock.DEA {
		return false
	}

	// RSI处于正常区间
	if stock.RSI < 30 || stock.RSI > 70 {
		return false
	}

	// 成交量放大
	if stock.Volume < stock.AvgVolume*0.8 {
		return false
	}

	return true
}

// shouldSell 判断是否应该卖出
func (sr *StockRecommendation) shouldSell(stock *StockDailyData) bool {
	// 价格跌破MA20均线
	if stock.Close >= stock.MA20 {
		return false
	}

	// MACD死叉
	if stock.MACD > 0 || stock.MACD > stock.DEA {
		return false
	}

	// RSI超买或超卖
	if stock.RSI > 30 && stock.RSI < 70 {
		return false
	}

	// 成交量萎缩
	if stock.Volume > stock.AvgVolume*0.8 {
		return false
	}

	return true
}

// calculateBuyScore 计算买入置信度得分
func (sr *StockRecommendation) calculateBuyScore(stock *StockDailyData) float64 {
	var score float64

	// 价格与MA20的距离
	maDistance := (stock.Close - stock.MA20) / stock.MA20
	if maDistance > 0.05 {
		score += 20
	} else if maDistance > 0.02 {
		score += 15
	} else {
		score += 10
	}

	// MACD位置
	if stock.MACD > 0.02 {
		score += 25
	} else if stock.MACD > 0.01 {
		score += 20
	} else {
		score += 15
	}

	// RSI位置
	if stock.RSI > 50 && stock.RSI < 70 {
		score += 25
	} else if stock.RSI > 40 && stock.RSI < 50 {
		score += 20
	} else {
		score += 15
	}

	// 成交量放大程度
	volumeRatio := stock.Volume / stock.AvgVolume
	if volumeRatio > 1.5 {
		score += 25
	} else if volumeRatio > 1.2 {
		score += 20
	} else {
		score += 15
	}

	return math.Round(score)
}

// calculateSellScore 计算卖出置信度得分
func (sr *StockRecommendation) calculateSellScore(stock *StockDailyData) float64 {
	var score float64

	// 价格与MA20的距离
	maDistance := (stock.MA20 - stock.Close) / stock.MA20
	if maDistance > 0.05 {
		score += 20
	} else if maDistance > 0.02 {
		score += 15
	} else {
		score += 10
	}

	// MACD位置
	if stock.MACD < -0.02 {
		score += 25
	} else if stock.MACD < -0.01 {
		score += 20
	} else {
		score += 15
	}

	// RSI位置
	if stock.RSI > 70 {
		score += 25
	} else if stock.RSI < 30 {
		score += 25
	} else {
		score += 15
	}

	// 成交量萎缩程度
	volumeRatio := stock.Volume / stock.AvgVolume
	if volumeRatio < 0.6 {
		score += 25
	} else if volumeRatio < 0.8 {
		score += 20
	} else {
		score += 15
	}

	return math.Round(score)
}

// GetStockRecommendations 获取股票推荐
func (sr *StockRecommendation) GetStockRecommendations(data []*StockDailyData) (*RecommendationResult, error) {
	buyRecommendations, err := sr.RecommendBuyStocks(data)
	if err != nil {
		return nil, err
	}

	sellRecommendations, err := sr.RecommendSellStocks(data)
	if err != nil {
		return nil, err
	}

	return &RecommendationResult{
		BuyRecommendations:  buyRecommendations,
		SellRecommendations: sellRecommendations,
		MarketTrend:         sr.analyzeMarketTrend(data),
		RecommendationTime:  time.Now(),
	}, nil
}

// analyzeMarketTrend 分析市场趋势
func (sr *StockRecommendation) analyzeMarketTrend(data []*StockDailyData) string {
	var trend string

	if len(data) < 30 {
		return "数据不足，无法分析市场趋势"
	}

	// 计算平均MA20和RSI
	var avgMA20, avgRSI float64
	for _, stock := range data {
		avgMA20 += stock.MA20
		avgRSI += stock.RSI
	}
	avgMA20 /= float64(len(data))
	avgRSI /= float64(len(data))

	if avgMA20 > 0 && avgRSI > 50 {
		trend = "上升趋势"
	} else if avgMA20 < 0 && avgRSI < 50 {
		trend = "下降趋势"
	} else {
		trend = "震荡趋势"
	}

	return trend
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
