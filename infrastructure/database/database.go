package database

import (
	"database/sql"
	"date-apps-be/infrastructure/config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// DB component
type DB struct {
	Master *sql.DB
	Slave  *sql.DB
}

// InitializeDatabase to initialize database
func InitializeDatabase(conf *config.Config) (db *DB, err error) {
	if conf.DBMaster == nil || conf.DBSlave == nil {
		return nil, nil
	}
	db = &DB{}

	confMaster := conf.DBMaster
	db.Master, err = sql.Open("mysql", confMaster.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB master connection. %+v", err)
	}
	db.Master.SetMaxIdleConns(confMaster.MaxIdle)
	db.Master.SetMaxOpenConns(confMaster.MaxOpen)
	err = db.Master.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB master. %+v", err)
	}

	confSlave := conf.DBSlave
	db.Slave, err = sql.Open("mysql", confSlave.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB slave connection. %+v", err)
	}
	db.Slave.SetMaxIdleConns(confSlave.MaxIdle)
	db.Slave.SetMaxOpenConns(confSlave.MaxOpen)
	err = db.Slave.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB slave. %+v", err)
	}
	return db, err
}
