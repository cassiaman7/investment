package sina

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cassiaman7/investment/pkg/myhttp"
	"github.com/cassiaman7/investment/stock_api/variables"
)

const (
	SinaEndpoint   = "http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData"
	SinaDefaultLen = 73000
)

type Client struct {
	DataLen  int
	Endpoint string
}

func NewClient() *Client {
	return &Client{
		DataLen:  SinaDefaultLen,
		Endpoint: SinaEndpoint,
	}
}

type RespStockQuote struct {
	Day    string `json:"day"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

func (c RespStockQuote) toCell() variables.Quote {
	date, _ := time.Parse("2006-01-02", c.Day)
	open, _ := strconv.ParseFloat(c.Open, 64)
	high, _ := strconv.ParseFloat(c.High, 64)
	low, _ := strconv.ParseFloat(c.Low, 64)
	close, _ := strconv.ParseFloat(c.Close, 64)
	vol, _ := strconv.ParseFloat(c.Volume, 64)

	return variables.Quote{
		Day:    date,
		Open:   open,
		High:   high,
		Low:    low,
		Close:  close,
		Volume: vol,
	}
}

/*
	GetStockData 获取股票数据

参数：股票代码，开始时间，结束时间
返回：股票数据
备注：
- 股票数据格式为：时间，开盘价，最高价，最低价，收盘价，成交量
- 当前采用Sina接口，后续考虑使用其他接口
*/
func (c *Client) GetStockDataByCode(code variables.Code) (data []variables.Quote, err error) {
	var respData []RespStockQuote
	url := fmt.Sprintf("%s?symbol=%s&scale=240&ma=no&datalen=%d&fq=0", c.Endpoint, code.UniMark(), c.DataLen)
	if err = myhttp.NewClient().HTTPGet(url, &respData); err != nil {
		return
	}
	data = make([]variables.Quote, 0, len(respData))
	for _, cell := range respData {
		data = append(data, cell.toCell())
	}

	return
}
