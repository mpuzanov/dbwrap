package dbwrap

import (
	"fmt"
	"net/url"
)

// Config структура для параметров соединения с БД.
type Config struct {
	DriverName   string `json:"driver_name" yaml:"driver_name" env:"DB_DRIVER_NAME" env-default:"sqlserver" envDefault:"sqlserver"`
	Host         string `json:"host" yaml:"host" env:"DB_HOST"`
	Port         int    `json:"port" yaml:"port" env:"DB_PORT" env-default:"1433" envDefault:"1433" env-description:"sql server port"`
	User         string `json:"user" yaml:"user" env:"DB_USER"`
	Password     string `json:"password" yaml:"password" env:"DB_PASSWORD"`
	Database     string `json:"database" yaml:"database" env:"DB_DATABASE"`
	APPName      string `json:"app_name" yaml:"app_name" env:"APP_NAME"`
	DSN          string `json:"dsn" yaml:"dsn" env:"DB_DSN"`
	Encrypt      string `json:"encrypt" yaml:"encrypt" env:"DB_ENCRYPT"`
	TimeoutQuery int    `json:"timeout_query" yaml:"timeout_query" env:"TIMEOUT_QUERY" env-default:"300" envDefault:"300"` // Second
}

// NewConfig создание конфига по умолчанию.
func NewConfig(driverName string) *Config {
	c := &Config{Host: "127.0.0.1",
		TimeoutQuery: 300,
		DriverName:   driverName,
	}
	switch c.DriverName {
	case "sqlserver":
		c.Port = 1433
		c.User = "sa"
		c.Database = "master"
	case "postgres":
		c.Port = 5432
		c.User = "postgres"
		c.Database = "postgres"
	case "mysql":
		c.Port = 3306
		c.User = "root"
		c.Database = "mysql"
	}

	return c
}

// WithPassword установка пароля.
func (c *Config) WithPassword(pwd string) *Config {
	c.Password = pwd
	return c
}

// WithDriverName задания наименования драйвера БД.
func (c *Config) WithDriverName(driverName string) *Config {
	c.DriverName = driverName
	return c
}

// WithPort установка порта БД.
func (c *Config) WithPort(port int) *Config {
	c.Port = port
	return c
}

// WithDB установка БД.
func (c *Config) WithDB(dbname string) *Config {
	c.Database = dbname
	return c
}

// WithDSN установка строки подключения к БД.
func (c *Config) WithDSN(dsn string) *Config {
	c.DSN = dsn
	return c
}

// GetDatabaseURL
// "sqlserver://user:password@host:port?database=database_name"
// "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
// "mysql://username:password@protocol(address)/dbname?param=value"  "user:password@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
// sqlite3 или имя БД или :memory:
// driver sqlserver || postgres
func (c *Config) GetDatabaseURL() string {
	if c.DSN != "" {
		return c.DSN
	}

	switch c.DriverName {
	case "sqlite3":
		if c.Database != "" {
			return c.Database
		}
		return ":memory:"
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4,utf8&parseTime=true&loc=Local", c.User, c.Password, c.Host, c.Database)
	default:
		v := url.Values{}
		v.Set("database", c.Database)
		if c.APPName != "" {
			v.Add("app name", c.APPName)
		}
		switch c.DriverName {
		case "sqlserver":
			if c.Encrypt != "" {
				v.Set("encrypt", c.Encrypt)
			}
		case "postgres":
			v.Set("sslmode", "disable")
		}
		var u = url.URL{
			Scheme:   c.DriverName,
			User:     url.UserPassword(c.User, c.Password),
			Host:     fmt.Sprintf("%s:%d", c.Host, c.Port),
			RawQuery: v.Encode(),
		}
		return u.String()
	}
}

// String вывод полей в строку
func (c *Config) String() string {
	return fmt.Sprintf("DriverName=%s, Host=%s, Port=%d, User=%s, Password=<REMOVED>, Database=%s, TimeoutQuery=%d",
		c.DriverName, c.Host, c.Port, c.User, c.Database, c.TimeoutQuery)
}
