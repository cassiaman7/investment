package meta

import (
	"github.com/cassiaman7/investment/pkg/db"
	"gorm.io/gorm"
)

var MetaDB DBMeta

type DBMeta struct {
	Config *db.DBConfig
	Orm    *gorm.DB
}

func (c *DBMeta) Init(config *db.DBConfig) (err error) {
	c.Config = config
	c.Orm, err = config.GetORMConn()

	return err
}
