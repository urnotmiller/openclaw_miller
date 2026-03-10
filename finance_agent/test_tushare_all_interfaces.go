package main

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
	"time"
)

func main() {
	fmt.Println("=== Tushare API所有接口测试 ===\n")

	// 直接使用配置文件中的Token
	token := "eee505dc939712a8cd0dfd3a7eb0271ca620b7824566b1ca8a3d6f4b"

	// 测试参数
	startDate := time.Now().AddDate(0, 0, -30).Format("20060102")
	endDate := time.Now().Format("20060102")

	fmt.Printf("测试时间范围: %s 到 %s\n", startDate, endDate)
	fmt.Println()

	// 测试获取股票基础信息
	fmt.Println("1. 测试获取股票基础信息")
	fmt.Println("------------------------")

	stockBasic, err := getStockBasic(token)
	if err != nil {
		fmt.Printf("❌ 获取股票基础信息失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功获取 %d 条股票基础信息\n", len(stockBasic))
		if len(stockBasic) > 0 {
			for i, data := range stockBasic[:min(5, len(stockBasic))] {
				fmt.Printf("  %d. TS代码: %s\n", i+1, data["ts_code"])
				fmt.Printf("     股票名称: %s\n", data["name"])
				fmt.Printf("     行业: %s\n", data["industry"])
				fmt.Printf("     市场: %s\n", data["market"])
				fmt.Println()
			}
		}
	}

	// 测试获取每日行情数据
	fmt.Println("2. 测试获取每日行情数据")
	fmt.Println("------------------------")

	dailyData, err := getDailyData(token, startDate, endDate)
	if err != nil {
		fmt.Printf("❌ 获取每日数据失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功获取 %d 条每日行情数据\n", len(dailyData))
		if len(dailyData) > 0 {
			for i, data := range dailyData[:min(5, len(dailyData))] {
				fmt.Printf("  %d. TS代码: %s\n", i+1, data["ts_code"])
				fmt.Printf("     日期: %s\n", data["trade_date"])
				fmt.Printf("     价格: %.2f (开盘: %.2f, 最高: %.2f, 最低: %.2f)\n",
					data["close"], data["open"], data["high"], data["low"])
				fmt.Printf("     成交量: %.0f, 成交额: %.0f\n", data["vol"], data["amount"])
				fmt.Println()
			}
		}
	}

	// 测试获取大盘指数数据
	fmt.Println("3. 测试获取大盘指数数据")
	fmt.Println("------------------------")

	indexDailyData, err := getIndexDailyData(token, startDate, endDate)
	if err != nil {
		fmt.Printf("❌ 获取大盘指数数据失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功获取 %d 条大盘指数数据\n", len(indexDailyData))
		if len(indexDailyData) > 0 {
			for i, data := range indexDailyData[:min(5, len(indexDailyData))] {
				fmt.Printf("  %d. TS代码: %s\n", i+1, data["ts_code"])
				fmt.Printf("     日期: %s\n", data["trade_date"])
				fmt.Printf("     价格: %.2f (开盘: %.2f, 最高: %.2f, 最低: %.2f)\n",
					data["close"], data["open"], data["high"], data["low"])
				fmt.Printf("     成交量: %.0f, 成交额: %.0f\n", data["vol"], data["amount"])
				fmt.Println()
			}
		}
	}

	// 测试获取股票名单数据
	fmt.Println("4. 测试获取股票名单数据")
	fmt.Println("------------------------")

	stockListData, err := getStockListData(token)
	if err != nil {
		fmt.Printf("❌ 获取股票名单数据失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功获取 %d 条股票名单数据\n", len(stockListData))
		if len(stockListData) > 0 {
			for i, data := range stockListData[:min(5, len(stockListData))] {
				fmt.Printf("  %d. TS代码: %s\n", i+1, data["ts_code"])
				fmt.Printf("     股票名称: %s\n", data["name"])
				fmt.Printf("     行业: %s\n", data["industry"])
				fmt.Printf("     市场: %s\n", data["market"])
				fmt.Println()
			}
		}
	}

	// 总结
	fmt.Println("5. 测试总结")
	fmt.Println("----------------")

	if stockBasic != nil && dailyData != nil && indexDailyData != nil && stockListData != nil {
		fmt.Println("✅ 所有接口测试成功！")
	} else {
		fmt.Println("⚠️  部分接口测试失败或无数据返回")
	}

	fmt.Println()
	fmt.Println("=== 测试完成 ===")
}

func getStockBasic(token string) ([]map[string]interface{}, error) {
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

	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
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

	err = json.Unmarshal(body.Bytes(), &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("Tushare API error: %s", result.Message)
	}

	var data []map[string]interface{}
	for _, item := range result.Data.Items {
		row := make(map[string]interface{})
		for j, field := range result.Data.Fields {
			row[field] = item[j]
		}
		data = append(data, row)
	}

	return data, nil
}

func getDailyData(token, startDate, endDate string) ([]map[string]interface{}, error) {
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

	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
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

	err = json.Unmarshal(body.Bytes(), &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("Tushare API error: %s", result.Message)
	}

	var data []map[string]interface{}
	for _, item := range result.Data.Items {
		row := make(map[string]interface{})
		for j, field := range result.Data.Fields {
			row[field] = item[j]
		}
		data = append(data, row)
	}

	return data, nil
}

func getIndexDailyData(token, startDate, endDate string) ([]map[string]interface{}, error) {
	reqData := map[string]interface{}{
		"api_name": "index_daily",
		"token":    token,
		"params": map[string]string{
			"ts_code":    "000001.SH",
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

	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
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

	err = json.Unmarshal(body.Bytes(), &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("Tushare API error: %s", result.Message)
	}

	var data []map[string]interface{}
	for _, item := range result.Data.Items {
		row := make(map[string]interface{})
		for j, field := range result.Data.Fields {
			row[field] = item[j]
		}
		data = append(data, row)
	}

	return data, nil
}

func getStockListData(token string) ([]map[string]interface{}, error) {
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

	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
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

	err = json.Unmarshal(body.Bytes(), &result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("Tushare API error: %s", result.Message)
	}

	var data []map[string]interface{}
	for _, item := range result.Data.Items {
		row := make(map[string]interface{})
		for j, field := range result.Data.Fields {
			row[field] = item[j]
		}
		data = append(data, row)
	}

	return data, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
