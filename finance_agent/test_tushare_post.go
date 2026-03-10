package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	fmt.Println("=== Tushare API POST请求测试 ===\n")

	// 测试参数
	token := "eee505dc939712a8cd0dfd3a7eb0271ca620b7824566b1ca8a3d6f4b"
	startDate := time.Now().AddDate(0, 0, -7).Format("20060102")
	endDate := time.Now().Format("20060102")

	fmt.Printf("Token: %s\n", token)
	fmt.Printf("时间范围: %s 到 %s\n\n", startDate, endDate)

	// 1. 测试stock_basic接口
	fmt.Println("1. 测试stock_basic接口")
	fmt.Println("--------------------------------------------------")

	stockBasic, err := testStockBasic(token)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功: 获取到 %d 条股票基础信息\n", len(stockBasic))
		for i, stock := range stockBasic[:3] {
			fmt.Printf("   %d: %s (%s) - %s\n", i+1, stock.Name, stock.TSCode, stock.Industry)
		}
	}
	fmt.Println()

	// 2. 测试daily接口
	fmt.Println("2. 测试daily接口")
	fmt.Println("--------------------------------------------------")

	dailyData, err := testDaily(token, startDate, endDate)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功: 获取到 %d 条每日行情数据\n", len(dailyData))
		for i, data := range dailyData[:3] {
			fmt.Printf("   %d: %s - %.2f (%.0f手)\n", i+1, data.TSCode, data.Close, data.Volume)
		}
	}
	fmt.Println()

	// 3. 测试API响应格式
	fmt.Println("3. 测试API响应格式")
	fmt.Println("--------------------------------------------------")

	err = testAPIFormat(token)
	if err != nil {
		fmt.Printf("❌ 响应格式有问题: %v\n", err)
	} else {
		fmt.Printf("✅ API响应格式正常\n")
	}
	fmt.Println()

	// 4. 网络连通性测试
	fmt.Println("4. 网络连通性测试")
	fmt.Println("--------------------------------------------------")

	err = testNetworkConnectivity(token)
	if err != nil {
		fmt.Printf("❌ 网络连通性有问题: %v\n", err)
	} else {
		fmt.Printf("✅ 网络连通性正常\n")
	}

	fmt.Println("\n=== 测试完成 ===\n")

	if err == nil {
		fmt.Println("建议：")
		fmt.Println("1. API连接正常，可以继续进行数据获取和处理")
		fmt.Println("2. 如果需要获取更多数据，可以增加时间范围或使用其他接口")
		fmt.Println("3. 可以继续优化数据处理和技术指标计算")
	} else {
		fmt.Println("问题排查：")
		fmt.Println("1. 检查网络连接是否正常")
		fmt.Println("2. 验证Tushare Token是否正确")
		fmt.Println("3. 检查是否有防火墙或代理设置")
		fmt.Println("4. 尝试使用其他网络或设备进行测试")
	}
}

func testStockBasic(token string) ([]*StockBasic, error) {
	reqData := map[string]interface{}{
		"api_name": "stock_basic",
		"token":    token,
		"params": map[string]string{
			"list_status": "L",
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

	fmt.Printf("响应长度: %d\n", len(body))
	fmt.Printf("响应内容: %s\n", string(body))

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
		return nil, fmt.Errorf("JSON解析失败: %v, 响应内容: %s", err, string(body))
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

func testDaily(token, startDate, endDate string) ([]*DailyQuote, error) {
	reqData := map[string]interface{}{
		"api_name": "daily",
		"token":    token,
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

	fmt.Printf("响应长度: %d\n", len(body))
	fmt.Printf("响应内容: %s\n", string(body))

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
		return nil, fmt.Errorf("JSON解析失败: %v, 响应内容: %s", err, string(body))
	}

	if result.Code != 0 {
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

func testAPIFormat(token string) error {
	reqData := map[string]interface{}{
		"api_name": "stock_basic",
		"token":    token,
		"params": map[string]string{
			"list_status": "L",
			"limit":       "1",
		},
		"fields": "ts_code,name",
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://api.tushare.pro", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 检查是否是有效的JSON格式
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("JSON格式错误: %v, 响应内容: %s", err, string(body))
	}

	return nil
}

func testNetworkConnectivity(token string) error {
	// 尝试访问其他常见网站
	websites := []struct {
		name string
		url  string
	}{
		{"Tushare", "https://api.tushare.pro"},
		{"百度", "https://www.baidu.com"},
		{"新浪", "https://www.sina.com.cn"},
		{"腾讯", "https://www.qq.com"},
	}

	for _, site := range websites {
		resp, err := http.Head(site.url)
		if err != nil {
			fmt.Printf("   %s: ❌ %v\n", site.name, err)
		} else {
			fmt.Printf("   %s: ✅ %s\n", site.name, resp.Status)
			resp.Body.Close()
		}
	}

	// 直接访问Tushare API
	reqData := map[string]interface{}{
		"api_name": "stock_basic",
		"token":    token,
		"params": map[string]string{
			"list_status": "L",
			"limit":       "1",
		},
		"fields": "ts_code,name",
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://api.tushare.pro", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
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
