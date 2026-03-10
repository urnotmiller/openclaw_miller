package main

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
)

func main() {
	fmt.Println("=== Tushare API每日行情数据测试 ===\n")

	// 直接使用配置文件中的Token
	token := "eee505dc939712a8cd0dfd3a7eb0271ca620b7824566b1ca8a3d6f4b"

	// 测试参数 - 使用2023年的日期范围，这是一个已知有数据的时间段
	startDate := "20231201"
	endDate := "20231210"

	fmt.Printf("测试时间范围: %s 到 %s\n", startDate, endDate)
	fmt.Println()

	// 测试获取每日行情数据
	fmt.Println("1. 测试获取每日行情数据")
	fmt.Println("------------------------")

	dailyData, err := getDailyData(token, startDate, endDate)
	if err != nil {
		fmt.Printf("❌ 获取每日数据失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功获取 %d 条每日行情数据\n", len(dailyData))

	// 打印数据详情
	fmt.Println()
	fmt.Println("2. 数据详情")
	fmt.Println("----------------")

	if len(dailyData) > 0 {
		for i, data := range dailyData[:min(5, len(dailyData))] {
			fmt.Printf("  %d. TS代码: %s\n", i+1, data["ts_code"])
			fmt.Printf("     日期: %s\n", data["trade_date"])
			fmt.Printf("     价格: %.2f (开盘: %.2f, 最高: %.2f, 最低: %.2f)\n",
				data["close"], data["open"], data["high"], data["low"])
			fmt.Printf("     成交量: %.0f, 成交额: %.0f\n", data["vol"], data["amount"])
			fmt.Println()
		}
	} else {
		fmt.Println("   无数据返回")
	}

	// 总结
	fmt.Println("3. 测试总结")
	fmt.Println("----------------")

	if len(dailyData) > 0 {
		fmt.Println("✅ 每日行情数据获取成功！")
	} else {
		fmt.Println("⚠️  无数据返回，请检查日期范围和股票代码")
	}

	fmt.Println()
	fmt.Println("=== 测试完成 ===")
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

	fmt.Printf("响应状态: %s\n", resp.Status)
	fmt.Printf("响应内容: %s\n", string(body.Bytes()))

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
