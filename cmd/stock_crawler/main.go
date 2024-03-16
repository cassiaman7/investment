package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cassiaman7/investment/meta"
	"github.com/cassiaman7/investment/pkg/logic"
	"github.com/cassiaman7/investment/pkg/model"
	stockapi "github.com/cassiaman7/investment/stock_api"
	"github.com/cassiaman7/investment/stock_api/variables"
)

const (
	crawlDuration = 0
	crawlParall   = 10
)

func main() {
	if err := initDataSource(); err != nil {
		panic(err)
	}
	if err := crawlAHistory(); err != nil {
		panic(err)
	}
}

func crawlAHistory() error {
	AStocks, err := stockapi.GetAStockList()
	if err != nil {
		return err
	}
	fmt.Printf("found len(AStocks) = %d\n", len(AStocks))

	stockChan := make(chan variables.Code, 1000)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, stock := range AStocks {
			code := variables.Code{
				Market: variables.GetMarketByNumber(stock.Number),
				Number: stock.Number,
				Name:   stock.Name,
			}
			stockChan <- code
		}
		close(stockChan)
	}()

	var total, done int64 = int64(len(AStocks)), 0
	for i := 0; i < crawlParall; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				code, ok := <-stockChan
				if !ok {
					return
				}
			RETRY_MARK:
				if err = logic.InitStockData(code); err != nil {
					fmt.Printf("err occur %s %s %v, start retry...\n", code.UniMark(), code.Name, err)
					time.Sleep(10 * time.Minute)
					goto RETRY_MARK
				}
				time.Sleep(crawlDuration)
				atomic.AddInt64(&done, 1)
				fmt.Printf("current process: %.2f %% [%d/%d]\n", float64(done)*100/float64(total), done, total)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("all succ...")

	return nil
}

func initDataSource() error {
	if err := meta.MetaDB.Init(meta.StockConfig); err != nil {
		return err
	}

	if err := model.MigrateTables(meta.MetaDB.Orm, model.StockTables); err != nil {
		return err
	}

	return nil
}

func fetchAllStocks() error {
	AStocks, err := stockapi.GetAStockList()
	if err != nil {
		return err
	}
	fmt.Printf("found len(AStocks) = %d\n", len(AStocks))

	for _, stock := range AStocks {
		code := variables.Code{
			Market: variables.GetMarketByNumber(stock.Number),
			Number: stock.Number,
			Name:   stock.Name,
		}
		fmt.Printf("==%s\n", code.UniMark())
	}

	return nil
}
