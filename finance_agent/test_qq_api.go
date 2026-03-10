package main

import (
	"fmt"
	"strings"
	"time"
)

// 测试腾讯财经API对接（使用正式代码中的函数）
func main() {
	fmt.Println("=== 腾讯财经API正式代码测试 ===")
	fmt.Println()

	// 创建腾讯财经API实例
	api := NewQQFinanceAPI()

	// 测试股票列表（A股、B股、港股）
	testStocks := []struct {
		Market string
		Code   string
		Name   string
	}{
		{"A股", "600000", "浦发银行"},
		{"A股", "000001", "平安银行"},
		{"A股", "600519", "贵州茅台"},
		{"B股", "900956", "东贝B股"},
		{"港股", "00001", "汇丰控股"},
		{"港股", "00700", "腾讯控股"},
	}

	// 测试获取单个股票实时行情
	fmt.Println("1. 测试单个股票实时行情获取:")
	fmt.Println(strings.Repeat("-", 50))

	for i, stock := range testStocks {
		fmt.Printf("测试 %d: %s (%s) - %s\n", i+1, stock.Market, stock.Code, stock.Name)

		quote, err := api.GetStockQuote(stock.Market, stock.Code)
		if err != nil {
			fmt.Printf("❌ 失败: %v\n\n", err)
			continue
		}

		// 打印详细行情信息
		fmt.Printf("✅ 成功: %s (%s)\n", quote.Name, quote.Code)
		fmt.Printf("   当前价格: %.2f\n", quote.Price)
		fmt.Printf("   开盘价: %.2f | 最高价: %.2f | 最低价: %.2f\n", quote.Open, quote.High, quote.Low)
		fmt.Printf("   前收盘价: %.2f\n", quote.PreviousClose)
		fmt.Printf("   涨跌幅: %.2f%% | 涨跌额: %.2f\n", quote.ChangeRatio, quote.Change)
		fmt.Printf("   成交量: %.0f手 | 成交额: %.2f万元\n", quote.Volume, quote.Amount)
		fmt.Printf("   更新时间: %s\n\n", quote.Time.Format("2006-01-02 15:04:05"))

		// 测试每3只股票后暂停一下，避免请求过快
		if (i+1)%3 == 0 {
			fmt.Println("--- 暂停1秒，避免请求过快 ---")
			<-time.After(1 * time.Second)
		}
	}

	// 测试批量获取股票实时行情
	fmt.Println("2. 测试批量获取股票实时行情:")
	fmt.Println(strings.Repeat("-", 50))

	testBatchStocks := []struct {
		Market string
		Code   string
	}{
		{"A股", "600000"},
		{"A股", "000001"},
		{"A股", "600519"},
	}

	var stockInfos []StockInfo
	for _, stock := range testBatchStocks {
		stockInfos = append(stockInfos, StockInfo{
			Market: stock.Market,
			Code:   stock.Code,
		})
	}

	quotes, err := api.GetMultipleQuotes(stockInfos)
	if err != nil {
		fmt.Printf("❌ 批量获取失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功获取 %d 只股票的实时行情\n", len(quotes))
		for _, quote := range quotes {
			fmt.Printf("   %s (%s): %.2f (%.2f%%)\n", quote.Name, quote.Code, quote.Price, quote.ChangeRatio)
		}
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("提示: 如果测试中出现'股票代码不存在'等错误，可能是该股票已退市或代码错误")
}
