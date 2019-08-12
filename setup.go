package dbUtil

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
	Database string
	PoolSize int
}

func SetupMysqlDb(c *MySQLConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&autocommit=true&parseTime=True",
		c.UserName,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxIdleConns(c.PoolSize)
	db.SetMaxOpenConns(c.PoolSize)
	return db
}
