package main

import (
	"fmt"
	"math"
	"time"
)

// StockAnalysis 股票分析功能
type StockAnalysis struct{}

// NewStockAnalysis 创建股票分析实例
func NewStockAnalysis() *StockAnalysis {
	return &StockAnalysis{}
}

// CalculateMA 计算移动平均线
func (sa *StockAnalysis) CalculateMA(data []*DailyQuote, period int) ([]float64, error) {
	if len(data) < period {
		return nil, fmt.Errorf("数据量不足，无法计算%d日均线", period)
	}

	var ma []float64
	for i := period - 1; i < len(data); i++ {
		var sum float64
		for j := 0; j < period; j++ {
			sum += data[i-j].Close
		}
		ma = append(ma, sum/float64(period))
	}

	return ma, nil
}

// CalculateMACD 计算MACD指标
func (sa *StockAnalysis) CalculateMACD(data []*DailyQuote) ([]float64, []float64, []float64, error) {
	// 计算EMA12和EMA26
	ema12, err := sa.CalculateEMA(data, 12)
	if err != nil {
		return nil, nil, nil, err
	}

	ema26, err := sa.CalculateEMA(data, 26)
	if err != nil {
		return nil, nil, nil, err
	}

	// 计算DIF（EMA12 - EMA26）
	dif := make([]float64, 0, len(data))
	for i := range data {
		if i >= 25 { // EMA26需要至少26个数据点
			dif = append(dif, ema12[i]-ema26[i])
		}
	}

	// 计算DEA（DIF的EMA9）
	deadata := make([]*DailyQuote, 0, len(dif))
	for _, d := range dif {
		deadata = append(deadata, &DailyQuote{Close: d})
	}

	dea, err := sa.CalculateEMA(deadata, 9)
	if err != nil {
		return nil, nil, nil, err
	}

	// 计算MACD柱状图（DIF - DEA）
	macd := make([]float64, 0, len(dif))
	for i, d := range dif {
		if i >= 8 { // DEA需要至少9个数据点
			macd = append(macd, (d-dea[i-8])*2)
		}
	}

	return dif, dea, macd, nil
}

// CalculateEMA 计算指数移动平均线
func (sa *StockAnalysis) CalculateEMA(data []*DailyQuote, period int) ([]float64, error) {
	if len(data) < period {
		return nil, fmt.Errorf("数据量不足，无法计算%d日EMA", period)
	}

	var ema []float64
	multiplier := 2.0 / float64(period+1)

	// 初始值为简单移动平均
	var initialSum float64
	for i := 0; i < period; i++ {
		initialSum += data[i].Close
	}
	initialEMA := initialSum / float64(period)
	ema = append(ema, initialEMA)

	for i := period; i < len(data); i++ {
		currentEMA := data[i].Close*multiplier + ema[i-period]*(1-multiplier)
		ema = append(ema, currentEMA)
	}

	return ema, nil
}

// CalculateRSI 计算RSI指标
func (sa *StockAnalysis) CalculateRSI(data []*DailyQuote, period int) ([]float64, error) {
	if len(data) < period+1 {
		return nil, fmt.Errorf("数据量不足，无法计算%d日RSI", period)
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
			losses = append(losses, math.Abs(change))
		} else {
			gains = append(gains, 0)
			losses = append(losses, 0)
		}
	}

	// 计算平均增益和平均损失
	avgGain := make([]float64, 0, len(gains))
	avgLoss := make([]float64, 0, len(losses))

	var sumGain float64
	var sumLoss float64
	for i := 0; i < period; i++ {
		sumGain += gains[i]
		sumLoss += losses[i]
	}
	avgGain = append(avgGain, sumGain/float64(period))
	avgLoss = append(avgLoss, sumLoss/float64(period))

	for i := period; i < len(gains); i++ {
		avgGain = append(avgGain, (avgGain[i-period]*(period-1)+gains[i])/float64(period))
		avgLoss = append(avgLoss, (avgLoss[i-period]*(period-1)+losses[i])/float64(period))
	}

	// 计算RSI
	for i, gain := range avgGain {
		if avgLoss[i] == 0 {
			rsi = append(rsi, 100)
		} else {
			rs := gain / avgLoss[i]
			rsi = append(rsi, 100-(100/(1+rs)))
		}
	}

	return rsi, nil
}

// GetTrendAnalysis 趋势分析
func (sa *StockAnalysis) GetTrendAnalysis(data []*DailyQuote, maPeriod int) string {
	if len(data) < maPeriod*2 {
		return "数据不足，无法进行趋势分析"
	}

	// 计算移动平均线
	ma, err := sa.CalculateMA(data, maPeriod)
	if err != nil {
		return "计算移动平均线失败"
	}

	// 趋势分析
	var trend string
	if ma[len(ma)-1] > ma[0] {
		trend = "上升趋势"
	} else if ma[len(ma)-1] < ma[0] {
		trend = "下降趋势"
	} else {
		trend = "震荡趋势"
	}

	// 价格与MA的关系
	if data[len(data)-1].Close > ma[len(ma)-1] {
		trend += "（价格在均线之上）"
	} else if data[len(data)-1].Close < ma[len(ma)-1] {
		trend += "（价格在均线之下）"
	} else {
		trend += "（价格与均线持平）"
	}

	return trend
}
