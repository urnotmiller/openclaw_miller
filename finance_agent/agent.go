package main

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("=== 金融专家智能体 ===")
	fmt.Println("正在启动...")
	fmt.Println("=====================")
}

// 获取股票基础信息（供智能体调用）
func GetStockBasicInfo(isListed bool) ([]*StockBasic, error) {
	// 从配置中读取Tushare Token
	token := getTushareToken()
	if token == "" {
		return nil, fmt.Errorf("Tushare Token未配置")
	}

	api := NewTushareAPI(token)
	return api.StockBasic(isListed)
}

// 获取每日行情数据（供智能体调用）
func GetDailyQuotes(startDate, endDate string) ([]*DailyQuote, error) {
	// 从配置中读取Tushare Token
	token := getTushareToken()
	if token == "" {
		return nil, fmt.Errorf("Tushare Token未配置")
	}

	api := NewTushareAPI(token)
	return api.Daily(startDate, endDate)
}

// 获取股票每日数据（包含技术指标）（供智能体调用）
func GetStockDailyData(startDate, endDate string) ([]*StockDailyData, error) {
	// 尝试使用Tushare API
	token := getTushareToken()
	if token != "" {
		api := NewTushareAPI(token)
		data, err := api.GetStockDailyData(startDate, endDate)
		if err == nil {
			return data, nil
		}
		fmt.Printf("Tushare API错误: %v，将使用备用数据方案\n", err)
	}

	// 备用数据方案
	alt := NewStockDataAlternative()
	return alt.GetStockDailyData(startDate, endDate)
}

// 获取实时行情数据（供智能体调用）
func GetRealTimeQuote(market, code string) (*RealTimeQuote, error) {
	api := NewQQFinanceAPI()
	return api.GetStockQuote(market, code)
}

// 搜索股票（供智能体调用）
func SearchStock(keyword string) ([]*StockInfo, error) {
	searcher := NewStockSearch()
	return searcher.SearchStock(keyword)
}

// 根据股票代码获取股票信息（供智能体调用）
func GetStockInfoByCode(market, code string) (*StockInfo, error) {
	searcher := NewStockSearch()
	return searcher.GetStockInfoByCode(market, code)
}

// 获取股票推荐（供智能体调用）
func GetStockRecommendations(data []*StockDailyData) (*RecommendationResult, error) {
	rec := NewStockRecommendation()
	return rec.GetStockRecommendations(data)
}

// 获取买入推荐（供智能体调用）
func GetBuyRecommendations(data []*StockDailyData) ([]*StockRecommendationResult, error) {
	rec := NewStockRecommendation()
	return rec.RecommendBuyStocks(data)
}

// 获取卖出推荐（供智能体调用）
func GetSellRecommendations(data []*StockDailyData) ([]*StockRecommendationResult, error) {
	rec := NewStockRecommendation()
	return rec.RecommendSellStocks(data)
}

// 从配置文件中获取Tushare Token
func getTushareToken() string {
	// 这里使用简单的配置读取方法，实际项目中可以使用viper等配置库
	const configFile = "./config.yaml"

	data, err := os.ReadFile(configFile)
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
