package dbwrap

import (
	"fmt"

	"errors"

	"github.com/jmoiron/sqlx"
)

// DBSQL ...
type DBSQL struct {
	DBX *sqlx.DB
	Cfg *Config
}

// ErrBadConfigDB ошибка
var ErrBadConfigDB = errors.New("не заполнены параметры подключения к БД")

// NewConnect Создание подключения к БД
func NewConnect(cfg *Config) (*DBSQL, error) {

	if cfg.DriverName != "sqlite3" && (cfg.Host == "" || cfg.Database == "" || cfg.User == "") {
		return nil, ErrBadConfigDB
	}
	dsn := cfg.GetDatabaseURL()
	db, err := sqlx.Connect(cfg.DriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf(" sqlx.Connect : %w", err)
	}
	return &DBSQL{DBX: db, Cfg: cfg}, nil
}

// Close закрытие соединений
func (d *DBSQL) Close() error {
	return d.DBX.Close()
}

// SetTimeout установка таймаута для выполнения запроса в секундах
func (d *DBSQL) SetTimeout(timeout uint) {
	d.Cfg.TimeoutQuery = int(timeout)
}
