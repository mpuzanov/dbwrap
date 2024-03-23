package dbwrap

import (
	"fmt"
	"net/url"
)

// Config .
type Config struct {
	Host         string `yaml:"host" env:"DB_HOST" env-required:"true"`
	Port         int    `yaml:"port" env:"DB_PORT" env-default:"1433" env-description:"sql server port"`
	User         string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password     string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Database     string `yaml:"database" env:"DB_DATABASE" env-required:"true"`
	TimeoutQuery int    `yaml:"timeout_query" env:"TIMEOUT_QUERY" env-default:"300"` // Second
	APPName      string `yaml:"app_name" env:"APP_NAME"`
	DriverName   string `yaml:"driver_name" env:"DRIVER_NAME"  env-default:"sqlserver"`
}

// NewConfig создание конфига по умолчанию
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

// WithPassword установка пароля
func (c *Config) WithPassword(pwd string) *Config {
	c.Password = pwd
	return c
}

// WithDriverName задания наименования драйвера БД
func (c *Config) WithDriverName(driverName string) *Config {
	c.DriverName = driverName
	return c
}

// WithPort установка порта БД
func (c *Config) WithPort(port int) *Config {
	c.Port = port
	return c
}

// WithDB установка БД
func (c *Config) WithDB(dbname string) *Config {
	c.Database = dbname
	return c
}

// GetDatabaseURL
// "sqlserver://user:password@host:port?database=database_name"
// "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
// "mysql://username:password@protocol(address)/dbname?param=value"  "user:password@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
// driver sqlserver || postgres
func (c *Config) GetDatabaseURL() string {
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

		if c.DriverName == "postgres" {
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

// String выводо полей в строку
func (c *Config) String() string {
	return fmt.Sprintf("DriverName=%s, Host=%s, Port=%d, User=%s, Password=<REMOVED>, Database=%s, TimeoutQuery=%d",
		c.DriverName, c.Host, c.Port, c.User, c.Database, c.TimeoutQuery)
}
