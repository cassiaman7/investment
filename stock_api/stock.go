package stockapi

import (
	"time"

	"github.com/cassiaman7/investment/pkg/str"
	"github.com/cassiaman7/investment/stock_api/dfcf"
	"github.com/cassiaman7/investment/stock_api/sina"
	"github.com/cassiaman7/investment/stock_api/variables"
)

// 获取交易日信息
func GetTradeDate(start, end time.Time) (dates []time.Time) {
	code := variables.Code{
		Market: variables.MarketSH,
		Number: "601766", // 中国中车上市时间非常久
	}
	data, err := GetStockTimeSeriesByCode(code)
	if err != nil {
		return
	}
	dates = make([]time.Time, len(data))
	for _, v := range data {
		if v.Day.After(start) && v.Day.Before(end) {
			dates = append(dates, v.Day)
		}
	}

	return
}

// 是否为交易日
func IsTradeDate(date time.Time) bool {
	dates := GetTradeDate(date.Add(7*24*time.Hour), date.Add(-7*24*time.Hour))
	for _, v := range dates {
		if str.ToDate(date) == str.ToDate(v) {
			return true
		}
	}

	return false
}

// 获取股票数据 时间序列数据
func GetStockTimeSeriesByCode(code variables.Code) (data []variables.Quote, err error) {
	return sina.NewClient().GetStockDataByCode(code)
}

// 获取最新的全A股票数据
func GetAStockList() (data []variables.SFull, err error) {
	return dfcf.NewClient().GetAStockList()
}
