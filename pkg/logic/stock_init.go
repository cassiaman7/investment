package logic

import (
	"fmt"
	"time"

	"github.com/cassiaman7/investment/meta"
	"github.com/cassiaman7/investment/pkg/model"
	stockapi "github.com/cassiaman7/investment/stock_api"
	"github.com/cassiaman7/investment/stock_api/variables"
	talib "github.com/markcheno/go-talib"
	"gorm.io/gorm/clause"
)

func InitStockData(ts variables.Code) error {
	fmt.Printf("start %s\n", ts.UniMark())
	data, err := stockapi.GetStockTimeSeriesByCode(ts)
	if err != nil {
		return err
	}

	close := closeSeries(data)
	sma5, sma10, sma20 := smaFn(close, 5), smaFn(close, 10), smaFn(close, 20)
	sma60, sma120, sma250 := smaFn(close, 60), smaFn(close, 120), smaFn(close, 250)
	up, mid, lower := bollBandFn(close, 20)
	outMACD, outMACDSignal, outMACDHist := macdFn(close)

	QuoteRows := make([]model.Quote, len(data))
	for i, v := range data {
		QuoteRows[i] = model.Quote{
			BaseColumns: model.BaseColumns{
				CreateTime: time.Now(),
			},
			Date:   v.Day,
			Market: string(ts.Market),
			Number: ts.Number,
			Name:   ts.Name,
			Open:   v.Open,
			High:   v.High,
			Low:    v.Low,
			Close:  v.Close,
			Volume: v.Volume,
			SMA5:   sma5[i], SMA10: sma10[i], SMA20: sma20[i],
			SMA60: sma60[i], SMA120: sma120[i], SMA250: sma250[i],
			UpBoll: up[i], MidBoll: mid[i], LowBoll: lower[i],
			OutMACD: outMACD[i], OutMACDSignal: outMACDSignal[i],
			OutMACDHist: outMACDHist[i],
		}
		if i > 0 {
			QuoteRows[i].PreClose = QuoteRows[i-1].Close
		}

		err = meta.MetaDB.Orm.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&QuoteRows[i]).Error
		if err != nil {
			return fmt.Errorf("write [%s%s]-%s fail, err: %v",
				ts.Market, ts.Number, ts.Name, err)
		}
	}

	return nil
}

func closeSeries(data []variables.Quote) (b []float64) {
	b = make([]float64, len(data))
	for i, v := range data {
		b[i] = v.Close
	}

	return
}

func smaFn(closeSeries []float64, timePeriod int) (sma []float64) {
	if len(closeSeries) < timePeriod {
		return make([]float64, len(closeSeries))
	}
	return talib.Sma(closeSeries, timePeriod)
}

func bollBandFn(closeSeries []float64, timePeriod int) (up []float64, mid []float64, lower []float64) {
	if len(closeSeries) < timePeriod {
		return make([]float64, len(closeSeries)), make([]float64, len(closeSeries)), make([]float64, len(closeSeries))
	}
	return talib.BBands(closeSeries, timePeriod, 2, 2, talib.SMA)
}

func macdFn(closeSeries []float64) (outMACD []float64, outMACDSignal []float64, outMACDHist []float64) {
	if len(closeSeries) < 26 {
		return make([]float64, len(closeSeries)), make([]float64, len(closeSeries)), make([]float64, len(closeSeries))
	}
	return talib.Macd(closeSeries, 12, 26, 9)
}
