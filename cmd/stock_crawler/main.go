package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/cassiaman7/investment/meta"
	"github.com/cassiaman7/investment/pkg/logic"
	"github.com/cassiaman7/investment/pkg/model"
	stockapi "github.com/cassiaman7/investment/stock_api"
	"github.com/cassiaman7/investment/stock_api/variables"
	"github.com/oklog/run"
)

const (
	crawlDuration = 0
	crawlParall   = 1
)

func main() {
	if err := meta.MetaDB.Init(meta.StockConfig); err != nil {
		panic(err)
	}

	if err := model.MigrateTables(meta.MetaDB.Orm, model.StockTables); err != nil {
		panic(err)
	}

	AStocks, err := stockapi.GetAStockList()
	if err != nil {
		panic(err)
	}
	fmt.Printf("found len(AStocks) = %d\n", len(AStocks))

	var markets = []variables.Market{
		variables.MarketSH,
		variables.MarketSZ,
		variables.MarketBJ,
	}

	var stockChan = make(chan variables.Code, 1000)
	ctx, cancel := context.WithCancel(context.Background())

	var group run.Group
	var total, done int64 = int64(len(AStocks)), 0
	group.Add(func() error {
		for _, stock := range AStocks {
			code := variables.Code{
				Market: markets[stock.MarketType-1],
				Number: stock.Number,
				Name:   stock.Name,
			}

			select {
			case <-ctx.Done():
				return nil
			case stockChan <- code:
			}
		}
		close(stockChan)

		return nil
	}, func(err error) {
		cancel()
	})

	for i := 0; i < crawlParall; i++ {
		group.Add(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					code, ok := <-stockChan
					if !ok {
						return nil
					}
				RETRY_MARK:
					if err = logic.InitStockData(code); err != nil {
						fmt.Printf("err occur %s %s\n", code.UniMark(), code.Name)
						time.Sleep(10 * time.Minute)
						goto RETRY_MARK
					}
					time.Sleep(crawlDuration)
					atomic.AddInt64(&done, 1)
				}
			}
		}, func(err error) {
			cancel()
		})
	}

	group.Add(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				fmt.Printf("当前进度%.2f %% [%d/%d]\n", float64(done)*100/float64(total), done, total)
				time.Sleep(3 * time.Second)
			}
		}
	}, func(err error) {
		cancel()
	})

	if err = group.Run(); err != nil {
		fmt.Printf("fail...err: %v\n", err)
	} else {
		fmt.Println("succ...")
	}
}
