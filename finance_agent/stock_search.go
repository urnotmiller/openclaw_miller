package main

import (
	"fmt"
	"strings"
)

// StockSearch 股票搜索功能
type StockSearch struct{}

// NewStockSearch 创建股票搜索实例
func NewStockSearch() *StockSearch {
	return &StockSearch{}
}

// SearchStock 搜索股票
func (s *StockSearch) SearchStock(keyword string) ([]*StockInfo, error) {
	// 这里使用简单的模拟数据，实际应该从股票数据库中搜索
	// 常见的股票代码和名称
	stockDatabase := []*StockInfo{
		{"A股", "600000", "浦发银行"},
		{"A股", "000001", "平安银行"},
		{"A股", "600519", "贵州茅台"},
		{"A股", "000858", "五粮液"},
		{"A股", "000002", "万科A"},
		{"A股", "600036", "招商银行"},
		{"B股", "900956", "东贝B股"},
		{"港股", "00001", "汇丰控股"},
		{"港股", "00700", "腾讯控股"},
		{"港股", "00388", "香港交易所"},
	}

	keyword = strings.ToLower(keyword)
	var results []*StockInfo

	for _, stock := range stockDatabase {
		if strings.Contains(strings.ToLower(stock.Name), keyword) ||
			strings.Contains(stock.Code, keyword) ||
			strings.Contains(strings.ToLower(stock.Market), keyword) {
			results = append(results, stock)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("未找到与\"%s\"相关的股票", keyword)
	}

	return results, nil
}

// StockInfo 股票信息（包含名称）
type StockInfo struct {
	Market string
	Code   string
	Name   string
}

// GetStockInfoByCode 根据股票代码获取股票信息
func (s *StockSearch) GetStockInfoByCode(market, code string) (*StockInfo, error) {
	stockDatabase := []*StockInfo{
		{"A股", "600000", "浦发银行"},
		{"A股", "000001", "平安银行"},
		{"A股", "600519", "贵州茅台"},
		{"B股", "900956", "东贝B股"},
		{"港股", "00001", "汇丰控股"},
		{"港股", "00700", "腾讯控股"},
	}

	for _, stock := range stockDatabase {
		if stock.Market == market && stock.Code == code {
			return stock, nil
		}
	}

	return nil, fmt.Errorf("未找到%s市场的股票代码%s", market, code)
}
