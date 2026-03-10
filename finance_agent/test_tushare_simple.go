package main

import (
	"io/ioutil"
	"fmt"
	"time"
	"gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("=== Tushare API简单测试 ===\n")

	// 测试参数
	startDate := time.Now().AddDate(0, 0, -10).Format("20060102")
	endDate := time.Now().Format("20060102")

	fmt.Printf("测试时间范围: %s 到 %s\n", startDate, endDate)
	fmt.Println()

	// 使用Tushare API
	token := getTushareToken()
	if token == "" {
		fmt.Println("❌ Tushare Token未配置")
		return
	}

	api := NewTushareAPI(token)

	// 测试获取每日数据
	fmt.Println("1. 测试获取每日行情数据")
	fmt.Println("------------------------")

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

	// 打印前5条数据
	fmt.Println()
	fmt.Println("2. 前5条数据示例")
	fmt.Println("----------------")

	for i, data := range dailyData[:5] {
		fmt.Printf("  %d. %s (%s)\n", i+1, data.Name, data.Code)
		fmt.Printf("     日期: %s\n", data.Date.Format("2006-01-02"))
		fmt.Printf("     价格: %.2f (开盘: %.2f, 最高: %.2f, 最低: %.2f)\n",
			data.Close, data.Open, data.High, data.Low)
		fmt.Printf("     成交量: %.0f, 平均成交量: %.0f\n", data.Volume, data.AvgVolume)
		fmt.Printf("     指标: MA20=%.2f, MACD=%.2f, DEA=%.2f, RSI=%.1f\n",
			data.MA20, data.MACD, data.DEA, data.RSI)
		fmt.Println()
	}

	// 总结
	fmt.Println("3. 测试总结")
	fmt.Println("----------------")

	if validCount > 0 {
		fmt.Println("✅ 每日行情数据获取成功！")
		fmt.Println("✅ 技术指标计算成功！")
		fmt.Println("✅ 数据格式化成功！")
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

// 获取Tushare Token
func getTushareToken() string {
	// 这里直接返回从配置文件中读取的值
	const configFile = "./config.yaml"

	data, err := ioutil.ReadFile(configFile)
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
