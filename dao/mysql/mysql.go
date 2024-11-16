package mysql

import (
	"BlueBell/settings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(ctg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		ctg.User,
		ctg.Password,
		ctg.Host,
		ctg.Port,
		ctg.DB,
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("connect database error:", err)
		return
	}

	db.SetMaxIdleConns(ctg.MaxIdleConns)
	db.SetMaxOpenConns(ctg.MaxOpenConns)

	return
}

func Close() {
	if db != nil {
		db.Close()
	}
}
