package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"main/cfg"
	"main/log"
	"strings"
	"time"
)

var db *sql.DB

var (
	Users       *KeyTable
	SingleUsers *SingleKeyTable

	EmailCodes                *SingleKeyTable
	PendingEmailConfirmations *SingleKeyTable
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
	db.SetConnMaxLifetime(2 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	Users = NewKeyTable("users")
	SingleUsers = &SingleKeyTable{
		KeyTable: Users,
		key:      "name",
	}
	EmailCodes = NewSingleKeyTable("emailcodes", "name")
	PendingEmailConfirmations = NewSingleKeyTable("pending_email_confirmations", "name")
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
