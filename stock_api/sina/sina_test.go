package sina

import (
	"testing"

	"github.com/cassiaman7/investment/stock_api/variables"
)

func TestGetStockDataBy(t *testing.T) {
	// code := variables.Code{
	// 	Market: variables.MarketSH,
	// 	Number: "510050",
	// }
	code := variables.Code{
		Market: variables.MarketBJ,
		Number: "873833",
	}

	data, err := NewClient().GetStockDataByCode(code)
	if err != nil {
		t.Error(err)
	}
	t.Logf("fetch %s KLine length[%d] succ", code.UniMark(), len(data))
}
