package dfcf

import (
	"testing"
)

func TestGetAStockList(t *testing.T) {
	data, err := NewClient().GetAStockList()
	if err != nil {
		t.Error(err)
	}
	t.Logf("fetch A stocks[%d] succ", len(data))
}
