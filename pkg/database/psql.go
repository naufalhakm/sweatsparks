package database

import (
	"database/sql"
	"fmt"
	"sweatsparks/internal/config"
	"time"

	_ "github.com/lib/pq"
)

func NewPqSQLClient() (*sql.DB, error) {
	var (
		DB_User   = config.ENV.DBUserName
		DB_Pass   = config.ENV.DBUserPassword
		DB_Host   = config.ENV.DBHost
		DB_Port   = config.ENV.DBPort
		DB_DbName = config.ENV.DBName
	)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		DB_User, DB_Pass, DB_Host, DB_Port, DB_DbName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 60)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
