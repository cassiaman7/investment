package dfcf

import (
	"fmt"

	"github.com/cassiaman7/investment/pkg/myhttp"
	"github.com/cassiaman7/investment/stock_api/variables"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetAStockList() (data []variables.SFull, err error) {
	var resp struct {
		Data struct {
			Diff []variables.SFull `json:"diff"`
		} `json:"data"`
	}
	token := "bd1d9ddb04089700cf9c27f6f7426281"
	// https://11.push2.eastmoney.com/api/qt/clist/get?pn=1&pz=10000&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&wbp2u=|0|0|0|web&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23,m:0+t:81+s:2048&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152
	url := fmt.Sprintf("https://11.push2.eastmoney.com/api/qt/clist/get?pn=1&pz=10000&po=1&np=1&ut=%s&fltt=2&invt=2&wbp2u=|0|0|0|web&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23,m:0+t:81+s:2048&fields=f1,f12,f14",
		token)
	if err = myhttp.NewClient().HTTPGet(url, &resp); err != nil {
		return data, err
	}

	return resp.Data.Diff, err
}
