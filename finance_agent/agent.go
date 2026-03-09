package main

import (
	"context"
	"fmt"
	"log"
	"os"

	veagent "github.com/volcengine/veadk-go/agent/llmagent"
	"github.com/volcengine/veadk-go/common"
	"github.com/volcengine/veadk-go/tool/builtin_tools/web_search"
	"gopkg.in/yaml.v3"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
)

func main() {
	ctx := context.Background()
	cfg := veagent.Config{
		Config: llmagent.Config{
			Name:        "financeExpert",
			Description: "金融专家智能体，提供股票、理财等金融服务",
			Instruction: `你是一位严谨的金融专家，拥有股票和理财领域的专业知识。你的任务是：
1. 提供股票市场分析和投资建议
2. 解答理财相关问题
3. 提供金融产品的评估和比较
4. 风险提示和投资策略建议

你需要保持严谨的风格，使用专业术语，提供准确的信息。所有建议仅供参考，不构成投资建议。`,
		},
		ModelName: "doubao-seed-code",
	}

	// 创建搜索工具
	webSearch, err := web_search.NewWebSearchTool(&web_search.Config{})
	if err != nil {
		fmt.Printf("NewWebSearchTool failed: %v", err)
		return
	}

	cfg.Tools = []tool.Tool{webSearch}

	// 创建智能体
	a, err := veagent.New(&cfg)
	if err != nil {
		fmt.Printf("NewLLMAgent failed: %v", err)
		return
	}

	// 启动应用
	config := &launcher.Config{
		AgentLoader:    agent.NewSingleLoader(a),
		SessionService: session.InMemoryService(),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
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

// 获取实时行情数据（供智能体调用）
func GetRealTimeQuote(market, code string) (*RealTimeQuote, error) {
	api := NewQQFinanceAPI()
	return api.GetStockQuote(market, code)
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
