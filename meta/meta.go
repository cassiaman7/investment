package meta

import (
	"github.com/cassiaman7/investment/pkg/db"
	"gorm.io/gorm"
)

var MetaDB DBMeta

var StockConfig = &db.DBConfig{
	Driver:   db.MySQLDriver,
	Host:     "127.0.0.1",
	Port:     3306,
	User:     "stock_w",
	Password: "be1eb4f1a5ef5bf68fe84f40b19069d7",
	Database: "stock",
	Charset:  "utf8mb4",
}

type DBMeta struct {
	Config *db.DBConfig
	Orm    *gorm.DB
}

func (c *DBMeta) Init(config *db.DBConfig) (err error) {
	c.Config = config
	c.Orm, err = config.GetORMConn()

	return err
}
