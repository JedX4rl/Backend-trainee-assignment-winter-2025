package storage

import (
	cfg "Backend-trainee-assignment-winter-2025/internal/config/storageConfig"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

//const (
//	driver = "postgres"
//	host     = "host="
//	port     = "port="
//	user     = "user="
//	dbname   = "dbname="
//	password = "password="
//	sslmode  = "sslmode="
//)

func NewStorage(dbConfig cfg.StorageConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.DBName, dbConfig.Password, dbConfig.SSLMode))
	if err != nil {
		return nil, err //TODO add logs
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
