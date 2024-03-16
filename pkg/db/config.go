package db

import (
	"database/sql"
	"fmt"
	"net"
	"time"

	"golang.org/x/exp/slices"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MySQLDriver = "mysql"
	// PostgresDriver = "postgres"
)

// DbConfig - dsn info for mysql database connection
type DBConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Schema   string `json:"schema"`
	Charset  string `json:"charset"`
	SSLMode  string `json:"sslMode"`
}

// GetConn - get db connection using DbConfig
func (dbc *DBConfig) GetConn() (db *sql.DB, err error) {
	dsn := ""
	if dsn, err = dbc.GenerateDSN(); err != nil {
		return
	}

	if db, err = sql.Open(dbc.Driver, dsn); err != nil {
		return
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxIdleTime(3600 * time.Second)

	return db, nil
}

// GetORMConn - get db connection using DbConfig
func (dbc *DBConfig) GetORMConn() (db *gorm.DB, err error) {
	dsn := ""
	if dsn, err = dbc.GenerateDSN(); err != nil {
		return
	}

	gormConfig := &gorm.Config{
		DisableAutomaticPing: false,
	}
	switch dbc.Driver {
	// case PostgresDriver:
	// 	db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	case MySQLDriver:
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	}

	sqlDB, err := db.DB()
	if sqlDB != nil {
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetConnMaxIdleTime(3600 * time.Second)
	}

	return
}

// DSN - connStr for connect to db
/*
* - postgres: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
 */
func (dbc *DBConfig) GenerateDSN() (connStr string, err error) {
	if err = dbc.CheckValid(); err != nil {
		return
	}
	switch dbc.Driver {
	// case PostgresDriver:
	// 	connStr = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
	// 		PostgresDriver, dbc.User, dbc.Password, dbc.Host, dbc.Port, dbc.Database)
	// 	if dbc.Schema != "" && strings.ToLower(dbc.Schema) != "public" {
	// 		connStr = fmt.Sprintf("%s&search_path=%s", connStr, dbc.Schema)
	// 	}
	case MySQLDriver:
		connStr = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?autocommit=1&parseTime=true",
			dbc.User, dbc.Password, "tcp", dbc.Host, dbc.Port, dbc.Database)
		if dbc.Charset != "" {
			connStr = fmt.Sprintf("%s&charset=%s", connStr, dbc.Charset)
		}
	default:
		err = fmt.Errorf("not support for driver (%s)", dbc.Driver)
		return
	}

	return
}

func (dbc *DBConfig) CheckValid() (err error) {
	supportDrivers := []string{MySQLDriver}
	if !slices.Contains(supportDrivers, dbc.Driver) {
		return fmt.Errorf("not support for driver (%s)", dbc.Driver)
	}

	if IP := net.ParseIP(dbc.Host); IP == nil {
		return fmt.Errorf("host(%s) is invalid", dbc.Host)
	}
	if dbc.Port > 100000 || dbc.Port <= 0 {
		return fmt.Errorf("port(%d) is not in (0, 100000)", dbc.Port)
	}
	if dbc.User == "" {
		return fmt.Errorf("user(%s) is empty str", dbc.User)
	}
	if dbc.Password == "" {
		return fmt.Errorf("password(%s) is empty str", dbc.Password)
	}
	// if dbc.Driver == PostgresDriver && dbc.Database == "" {
	// 	return fmt.Errorf("you must set db, where driver is postgres")
	// }

	return
}
