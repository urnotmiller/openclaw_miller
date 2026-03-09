package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SinaFinanceAPI 新浪财经API配置
type SinaFinanceAPI struct{}

// NewSinaFinanceAPI 创建新浪财经API实例
func NewSinaFinanceAPI() *SinaFinanceAPI {
	return &SinaFinanceAPI{}
}

// GetStockQuote 获取股票实时行情
func (api *SinaFinanceAPI) GetStockQuote(market, code string) (*RealTimeQuote, error) {
	var url string
	switch market {
	case "A股":
		if strings.HasPrefix(code, "6") {
			url = fmt.Sprintf("https://hq.sinajs.cn/list=sh%s", code)
		} else {
			url = fmt.Sprintf("https://hq.sinajs.cn/list=sz%s", code)
		}
	case "B股":
		if strings.HasPrefix(code, "90") {
			url = fmt.Sprintf("https://hq.sinajs.cn/list=sh%s", code)
		} else {
			url = fmt.Sprintf("https://hq.sinajs.cn/list=sz%s", code)
		}
	case "港股":
		url = fmt.Sprintf("https://hq.sinajs.cn/list=hk%s", code)
	default:
		return nil, fmt.Errorf("不支持的市场类型: %s", market)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body string
	_, err = fmt.Fscanln(resp.Body, &body)
	if err != nil {
		return nil, err
	}

	// 解析响应数据
	data := strings.TrimPrefix(body, fmt.Sprintf("var hq_str_%s=\"", getSinaStockCode(market, code)))
	data = strings.TrimSuffix(data, "\";")

	if data == body {
		return nil, fmt.Errorf("股票代码不存在或数据获取失败")
	}

	return parseSinaQuote(data)
}

// getSinaStockCode 获取新浪财经格式的股票代码
func getSinaStockCode(market, code string) string {
	switch market {
	case "A股":
		if strings.HasPrefix(code, "6") {
			return "sh" + code
		}
		return "sz" + code
	case "B股":
		if strings.HasPrefix(code, "90") {
			return "sh" + code
		}
		return "sz" + code
	case "港股":
		return "hk" + code
	default:
		return ""
	}
}

// parseSinaQuote 解析新浪财经行情数据
func parseSinaQuote(data string) (*RealTimeQuote, error) {
	reader := csv.NewReader(strings.NewReader(data))
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	fields, err := reader.Read()
	if err != nil {
		return nil, err
	}

	if len(fields) < 32 {
		return nil, fmt.Errorf("数据字段不足，无法解析")
	}

	quote := &RealTimeQuote{}

	// 解析股票基本信息
	quote.Name = fields[0]
	quote.Code = fields[29]

	// 解析价格信息
	if price, err := strconv.ParseFloat(fields[1], 64); err == nil {
		quote.Price = price
	}
	if preClose, err := strconv.ParseFloat(fields[2], 64); err == nil {
		quote.PreviousClose = preClose
	}
	if open, err := strconv.ParseFloat(fields[3], 64); err == nil {
		quote.Open = open
	}
	if high, err := strconv.ParseFloat(fields[4], 64); err == nil {
		quote.High = high
	}
	if low, err := strconv.ParseFloat(fields[5], 64); err == nil {
		quote.Low = low
	}

	// 解析成交量和成交额
	if volume, err := strconv.ParseFloat(fields[8], 64); err == nil {
		quote.Volume = volume
	}
	if amount, err := strconv.ParseFloat(fields[9], 64); err == nil {
		quote.Amount = amount
	}

	// 解析涨跌幅
	if change, err := strconv.ParseFloat(fields[31], 64); err == nil {
		quote.Change = change
	}
	if changeRatio, err := strconv.ParseFloat(fields[32], 64); err == nil {
		quote.ChangeRatio = changeRatio
	}

	// 解析时间
	if timeStr := fields[30]; timeStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", timeStr); err == nil {
			quote.Time = t
		}
	}

	return quote, nil
}

// GetMultipleQuotes 获取多个股票的实时行情
func (api *SinaFinanceAPI) GetMultipleQuotes(stocks []StockInfo) ([]*RealTimeQuote, error) {
	var codes []string
	for _, stock := range stocks {
		codes = append(codes, getSinaStockCode(stock.Market, stock.Code))
	}

	url := fmt.Sprintf("https://hq.sinajs.cn/list=%s", strings.Join(codes, ","))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body string
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			body += string(buf[:n])
		}
		if err != nil {
			break
		}
	}

	quotes := make([]*RealTimeQuote, 0, len(stocks))
	lines := strings.Split(body, "\n")

	for i, line := range lines {
		if line == "" {
			continue
		}

		data := strings.TrimPrefix(line, fmt.Sprintf("var hq_str_%s=\"", codes[i]))
		data = strings.TrimSuffix(data, "\";")

		if data == line {
			continue
		}

		quote, err := parseSinaQuote(data)
		if err == nil {
			quotes = append(quotes, quote)
		}
	}

	return quotes, nil
}

// StockInfo 股票信息
type StockInfo struct {
	Market string
	Code   string
}

// RealTimeQuote 实时行情
type RealTimeQuote struct {
	Code            string
	Name            string
	Price           float64
	PreviousClose   float64
	Open            float64
	High            float64
	Low             float64
	Volume          float64
	Amount          float64
	Change          float64
	ChangeRatio     float64
	Time            time.Time
}
