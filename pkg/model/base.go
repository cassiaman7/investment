package model

import (
	"time"

	"gorm.io/gorm"
)

func MigrateTables(db *gorm.DB, tabs ...[]any) (err error) {
	for _, tablist := range tabs {
		for _, tab := range tablist {
			if !db.Migrator().HasTable(tab) {
				if err = db.AutoMigrate(tab); err != nil {
					return err
				}
			}
		}
	}

	return err
}

var StockTables = []interface{}{
	&Quote{},
}

type BaseColumns struct {
	ID         int64     `gorm:"primaryKey;comment:自增主键" json:"id"`
	CreateTime time.Time `gorm:"not null;default:'1970-01-01 00:00:00'" json:"createTime"`
}
