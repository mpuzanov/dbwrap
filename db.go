package dbwrap

import (
	"fmt"

	"errors"

	"github.com/jmoiron/sqlx"
)

// DBSQL ...
type DBSQL struct {
	DBX          *sqlx.DB
	timeoutQuery int // Second
}

// ErrBadConfigDB ошибка.
var ErrBadConfigDB = errors.New("не заполнены параметры подключения к БД")

// NewConnect Создание подключения к БД.
func NewConnect(cfg *Config) (*DBSQL, error) {
	if cfg.DriverName != "sqlite3" && cfg.DSN == "" {
		if cfg.Host == "" || cfg.Database == "" || cfg.User == "" {
			return nil, ErrBadConfigDB
		}
	}
	dsn := cfg.GetDatabaseURL()
	db, err := sqlx.Connect(cfg.DriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect driver %s dsn %s: %w", cfg.DriverName, dsn, err)
	}
	return &DBSQL{DBX: db, timeoutQuery: cfg.TimeoutQuery}, nil
}

// NewConnect Создание подключения к БД.
func NewConnectDSN(driver, dsn string) (*DBSQL, error) {
	if driver == "" || dsn == "" {
		return nil, ErrBadConfigDB
	}

	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect driver %s, dsn %s: %w", driver, dsn, err)
	}
	return &DBSQL{DBX: db, timeoutQuery: 600}, nil
}

// Close закрытие соединений.
func (d *DBSQL) Close() error {
	return d.DBX.Close()
}

// SetTimeout установка таймаута для выполнения запроса в секундах.
func (d *DBSQL) SetTimeoutQuery(timeout uint) {
	d.timeoutQuery = int(timeout)
}

// Timeout получение текущего таймаута для выполнения запроса в секундах.
func (d *DBSQL) TimeoutQuery() int {
	return d.timeoutQuery
}
