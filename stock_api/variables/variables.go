package variables

import "time"

const (
	MarketUnknown Market = "unknown"
	MarketSZ      Market = "sz" // 深圳交易所
	MarketSH      Market = "sh" // 上海交易所
	MarketBJ      Market = "bj"
)

type Market string

// 股票数据最小单元，一只股票一田
type Quote struct {
	Day    time.Time `json:"day"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume float64   `json:"volume"`
}

type Code struct {
	Market Market `json:"market"`
	Number string `json:"number"`
	Name   string `json:"name"`
}

func (c *Code) UniMark() string {
	return string(c.Market) + c.Number
}

func GetMarketByNumber(num string) Market {
	switch num[0] {
	case '6':
		return MarketSH
	case '0', '3':
		return MarketSZ
	case '4', '8':
		return MarketBJ
	}

	return MarketUnknown
}
