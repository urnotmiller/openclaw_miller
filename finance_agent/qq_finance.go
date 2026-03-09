package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// QQFinanceAPI 腾讯财经API配置
type QQFinanceAPI struct{}

// NewQQFinanceAPI 创建腾讯财经API实例
func NewQQFinanceAPI() *QQFinanceAPI {
	return &QQFinanceAPI{}
}

// GetStockQuote 获取股票实时行情
func (api *QQFinanceAPI) GetStockQuote(market, code string) (*RealTimeQuote, error) {
	var prefix string
	switch market {
	case "A股":
		if strings.HasPrefix(code, "6") {
			prefix = "sh"
		} else {
			prefix = "sz"
		}
	case "B股":
		if strings.HasPrefix(code, "90") {
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

	return parseQQQuote(data, prefix+code, market)
}

// parseQQQuote 解析腾讯财经行情数据
func parseQQQuote(data, fullCode, market string) (*RealTimeQuote, error) {
	fields := strings.Split(data, "~")
	if len(fields) < 43 {
		return nil, fmt.Errorf("数据字段不足，无法解析")
	}

	quote := &RealTimeQuote{
		Code: fullCode,
	}

	// 解析股票基本信息（处理字符编码问题）
	quote.Name = fixCharacterEncoding(fields[1])

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

// GetMultipleQuotes 获取多个股票的实时行情
func (api *QQFinanceAPI) GetMultipleQuotes(stocks []StockInfo) ([]*RealTimeQuote, error) {
	var codes []string
	for _, stock := range stocks {
		var prefix string
		switch stock.Market {
		case "A股":
			if strings.HasPrefix(stock.Code, "6") {
				prefix = "sh"
			} else {
				prefix = "sz"
			}
		case "B股":
			if strings.HasPrefix(stock.Code, "90") {
				prefix = "sh"
			} else {
				prefix = "sz"
			}
		case "港股":
			prefix = "hk"
		default:
			continue
		}
		codes = append(codes, prefix+stock.Code)
	}

	url := fmt.Sprintf("https://qt.gtimg.cn/q=%s", strings.Join(codes, ","))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	quotes := make([]*RealTimeQuote, 0, len(stocks))
	lines := strings.Split(string(body), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// 解析每只股票的数据
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}

		code := strings.TrimPrefix(parts[0], "v_")
		data := strings.TrimPrefix(parts[1], "\"")
		data = strings.TrimSuffix(data, "\";")

		// 确定市场类型
		var market string
		if strings.HasPrefix(code, "sh") {
			if strings.HasPrefix(code[2:], "90") {
				market = "B股"
			} else {
				market = "A股"
			}
		} else if strings.HasPrefix(code, "sz") {
			if strings.HasPrefix(code[2:], "2") {
				market = "B股"
			} else {
				market = "A股"
			}
		} else if strings.HasPrefix(code, "hk") {
			market = "港股"
		} else {
			continue
		}

		quote, err := parseQQQuote(data, code, market)
		if err == nil {
			quotes = append(quotes, quote)
		}
	}

	return quotes, nil
}

// fixCharacterEncoding 修复股票名称的字符编码问题
func fixCharacterEncoding(name string) string {
	// 腾讯财经API使用的是GBK编码，Go默认是UTF-8，需要转换
	// 这里使用简单的替换方法处理常见的乱码
	replacements := map[string]string{
		"�ַ�����": "浦发银行",
		"ƽ������": "平安银行",
		"����ę́": "贵州茅台",
		"��Ѷ�ع�": "腾讯控股",
		"����": "汇丰控股",
	}

	for old, new := range replacements {
		name = strings.ReplaceAll(name, old, new)
	}

	return name
}

// StockInfo 股票信息
type StockInfo struct {
	Market string
	Code   string
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
