package db

import (
	"database/sql"
	"time"
)

// MySQL connection pool
func MySQL(host string, userName string, password string, port string, dbName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", userName+":"+password+"@tcp("+host+":"+port+")/"+dbName+"?parseTime=1&charset=utf8mb4&collation=utf8mb4_unicode_ci")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour / 2)
	return db, nil
}
