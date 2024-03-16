package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cassiaman7/investment/meta"
	"github.com/cassiaman7/investment/pkg/db"
	"github.com/cassiaman7/investment/pkg/logic"
	"github.com/cassiaman7/investment/pkg/model"
	stockapi "github.com/cassiaman7/investment/stock_api"
	"github.com/cassiaman7/investment/stock_api/variables"
)

const (
	toolName      = "stock_crawler"
	toolVersion   = "1.0.0"
	crawlDuration = time.Duration(0)
	crawlParall   = int(10)
)

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Printf("err occur: %v\n", err)
		}
	}()

	dbConfig, err := initConfig()
	if err != nil || dbConfig == nil {
		return
	}

	if err = initDataSource(dbConfig); err != nil {
		return
	}
	if err = crawlAHistory(); err != nil {
		return
	}
}

func initConfig() (*db.DBConfig, error) {
	var dbType, host, user, passwd, dbName, schema string
	var port int64
	var help, version bool
	flag.BoolVar(&version, "version", false, "show version info")
	flag.BoolVar(&help, "help", false, "show help info")
	flag.StringVar(&dbType, "db_type", "mysql", "set dbType of dataSource, for example: mysql")
	flag.StringVar(&host, "host", "127.0.0.1", "set db host of dataSource")
	flag.StringVar(&user, "user", "stock_w", "set db user of dataSource")
	flag.StringVar(&passwd, "pass", "123456", "set passwd of dataSource")
	flag.StringVar(&dbName, "db", "stock", "set db name of dataSource")
	flag.StringVar(&schema, "schema", "public", "set db schema for dataSource")
	flag.Int64Var(&port, "port", 3306, "set db port for dataSource")
	flag.Parse()

	if help || len(os.Args) < 2 {
		flag.PrintDefaults()
		return nil, nil
	}
	if version {
		fmt.Printf("%s version %s\n", toolName, toolVersion)
		return nil, nil
	}

	dbConfig := &db.DBConfig{
		Driver:   dbType,
		Host:     host,
		Port:     port,
		User:     user,
		Password: passwd,
		Database: dbName,
		Schema:   schema,
	}
	if err := dbConfig.CheckValid(); err != nil {
		return nil, err
	}

	return dbConfig, nil
}

func initDataSource(dbConfig *db.DBConfig) error {
	if err := meta.MetaDB.Init(dbConfig); err != nil {
		return err
	}

	if err := model.MigrateTables(meta.MetaDB.Orm, model.StockTables); err != nil {
		return err
	}

	return nil
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
