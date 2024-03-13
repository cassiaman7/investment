package model

import "time"

type Quote struct {
	BaseColumns
	Date          time.Time `gorm:"not null;index:idx_date;comment:时间" json:"day"`
	Market        string    `gorm:"size:2;not null;index:idx_market_number,priority:1;default:'';comment:市场标记1:上交所,2:深交所,3:北交所" json:"market"`
	Number        string    `gorm:"size:6;not null;index:idx_market_number,priority:2;default:'';comment:证券号码" json:"number"`
	Name          string    `gorm:"size:16;not null;default:'';comment:证券名称" json:"name"`
	PreClose      float64   `gorm:"not null;default:0.0" json:"preClose"`
	Open          float64   `gorm:"not null;default:0.0" json:"open"`
	High          float64   `gorm:"not null;default:0.0" json:"high"`
	Low           float64   `gorm:"not null;default:0.0" json:"low"`
	Close         float64   `gorm:"not null;default:0.0" json:"close"`
	Volume        float64   `gorm:"not null;default:0.0" json:"volume"`
	SMA5          float64   `gorm:"not null;default:0.0" json:"sma5"`
	SMA10         float64   `gorm:"not null;default:0.0" json:"sma10"`
	SMA20         float64   `gorm:"not null;default:0.0" json:"sma20"`
	SMA60         float64   `gorm:"not null;default:0.0" json:"sma60"`
	SMA120        float64   `gorm:"not null;default:0.0" json:"sma120"`
	SMA250        float64   `gorm:"not null;default:0.0" json:"sma250"`
	UpBoll        float64   `gorm:"not null;default:0.0" json:"upBoll"`
	MidBoll       float64   `gorm:"not null;default:0.0" json:"midBoll"`
	LowBoll       float64   `gorm:"not null;default:0.0" json:"lowBoll"`
	OutMACD       float64   `gorm:"not null;default:0.0" json:"outMacd"`
	OutMACDSignal float64   `gorm:"not null;default:0.0" json:"outMacdSignal"`
	OutMACDHist   float64   `gorm:"not null;default:0.0" json:"outMacdHist"`
}

func (c *Quote) TableName() string {
	return "quote"
}
