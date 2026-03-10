package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 测试腾讯财经API对接
func main() {
	fmt.Println("=== 腾讯财经API简单测试 ===")
	fmt.Println()

	// 测试股票代码
	testCases := []struct {
		market string
		code   string
		name   string
	}{
		{"A股", "600000", "浦发银行"},
		{"A股", "000001", "平安银行"},
		{"A股", "600519", "贵州茅台"},
		{"港股", "00001", "汇丰控股"},
		{"港股", "00700", "腾讯控股"},
	}

	for i, tc := range testCases {
		fmt.Printf("测试 %d: %s (%s) - %s\n", i+1, tc.market, tc.code, tc.name)

		quote, err := getTencentQuote(tc.market, tc.code)
		if err != nil {
			fmt.Printf("❌ 失败: %v\n", err)
		} else {
			fmt.Printf("✅ 成功: %s (%s)\n", quote.Name, quote.Code)
			fmt.Printf("   价格: %.2f | 涨跌幅: %.2f%%\n", quote.Price, quote.ChangeRatio)
			fmt.Printf("   成交量: %.0f手 | 成交额: %.2f万元\n", quote.Volume, quote.Amount)
		}
		fmt.Println()

		// 避免请求过快
		if i < len(testCases)-1 {
			time.Sleep(1 * time.Second)
		}
	}
}

// getTencentQuote 从腾讯财经API获取实时行情
func getTencentQuote(market, code string) (*RealTimeQuote, error) {
	var prefix string
	switch market {
	case "A股":
		if strings.HasPrefix(code, "6") {
			prefix = "sh"
		} else {
			prefix = "sz"
		}
	case "港股":
		prefix = "hk"
	default:
		return nil, fmt.Errorf("不支持的市场类型: %s", market)
	}

	url := fmt.Sprintf("https://qt.gtimg.cn/q=%s%s", prefix, code)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应数据
	data := string(body)
	data = strings.TrimPrefix(data, fmt.Sprintf("v_%s%s=\"", prefix, code))
	data = strings.TrimSuffix(data, "\";")

	if data == string(body) {
		return nil, fmt.Errorf("股票代码不存在或数据获取失败")
	}

	return parseTencentQuote(data, prefix+code, market)
}

// parseTencentQuote 解析腾讯财经行情数据
func parseTencentQuote(data, fullCode, market string) (*RealTimeQuote, error) {
	fields := strings.Split(data, "~")
	if len(fields) < 43 {
		return nil, fmt.Errorf("数据字段不足，无法解析")
	}

	quote := &RealTimeQuote{
		Code: fullCode,
	}

	// 解析股票基本信息
	quote.Name = fields[1]

	// 解析价格信息
	if price, err := strconv.ParseFloat(fields[3], 64); err == nil {
		quote.Price = price
	}
	if preClose, err := strconv.ParseFloat(fields[4], 64); err == nil {
		quote.PreviousClose = preClose
	}
	if open, err := strconv.ParseFloat(fields[5], 64); err == nil {
		quote.Open = open
	}
	if high, err := strconv.ParseFloat(fields[33], 64); err == nil {
		quote.High = high
	}
	if low, err := strconv.ParseFloat(fields[34], 64); err == nil {
		quote.Low = low
	}

	// 解析成交量和成交额
	if volume, err := strconv.ParseFloat(fields[8], 64); err == nil {
		quote.Volume = volume
	}
	if amount, err := strconv.ParseFloat(fields[36], 64); err == nil {
		quote.Amount = amount / 10000 // 转换为万元
	}

	// 解析涨跌幅
	if change, err := strconv.ParseFloat(fields[31], 64); err == nil {
		quote.Change = change
	}
	if changeRatio, err := strconv.ParseFloat(fields[32], 64); err == nil {
		quote.ChangeRatio = changeRatio
	}

	// 解析时间
	if timeStr := fields[30]; timeStr != "" && len(timeStr) == 14 {
		if t, err := time.Parse("20060102150405", timeStr); err == nil {
			quote.Time = t
		}
	}

	return quote, nil
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
