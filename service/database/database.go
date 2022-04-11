package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/misteeka/eplidr"
	"main/cfg"
	"main/log"
	"strings"
	"time"
)

var db *sql.DB

var (
	Users       *eplidr.Table
	SingleUsers *eplidr.SingleKeyTable

	EmailCodes                *eplidr.SingleKeyTable
	PendingEmailConfirmations *eplidr.SingleKeyTable
)

func InitDatabase() error {
	var err error
	dataSource := "{user}:{password}@tcp(localhost:41091)/{db}"
	dataSource = strings.Replace(dataSource, "{user}", cfg.GetString("user"), 1)
	dataSource = strings.Replace(dataSource, "{password}", cfg.GetString("password"), 1)
	dataSource = strings.Replace(dataSource, "{db}", cfg.GetString("dbName"), 1)
	db, err = sql.Open("mysql", dataSource)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetMaxIdleConns(cfg.GetInt("maxIdleConns"))
	db.SetMaxOpenConns(cfg.GetInt("maxOpenConns"))

	Users = eplidr.NewTable(
		"users",
		1,
		[]string{},
		db,
	)
	SingleUsers = eplidr.SingleKeyImplementation(Users, "name")
	EmailCodes = eplidr.NewSingleKeyTable(
		"emailcodes",
		"name",
		1,
		[]string{},
		db,
	)
	PendingEmailConfirmations = eplidr.NewSingleKeyTable(
		"pending_email_confirmations",
		"name",
		1,
		[]string{},
		db,
	)
	return nil
}

func Exec(sql string, v ...any) (sql.Result, error) {
	return db.Exec(sql, v)
}

func Query(sql string, v ...any) (*sql.Rows, error) {
	return db.Query(sql, v)
}

func ReleaseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.ErrorLogger.Print(err.Error())
		return
	}
}
