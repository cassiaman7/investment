package sina

import (
	"testing"

	"github.com/cassiaman7/investment/pkg/str"
	"github.com/cassiaman7/investment/stock_api/variables"
)

func TestGetStockDataBy(t *testing.T) {
	cli := NewClient()
	data, err := cli.GetStockDataByCode(variables.Code{
		Market: variables.MarketSH,
		Number: "510050",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(str.ToJSON(data))
}
