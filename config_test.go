package dbwrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseUrl(t *testing.T) {

	var tests = []struct {
		in   Config
		want string
	}{
		{Config{DriverName: "sqlserver"}, "sqlserver://:@:0?database="},
		{Config{Host: "host", Port: 1433, Database: "database_name", User: "user", Password: "password", DriverName: "sqlserver"},
			"sqlserver://user:password@host:1433?database=database_name"},
		{Config{Host: "host", Port: 1433, Database: "database_name", User: "user", Password: "password", APPName: "APP", DriverName: "sqlserver"},
			"sqlserver://user:password@host:1433?app+name=APP&database=database_name"},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.in.GetDatabaseURL())
	}
}

func TestConfigString(t *testing.T) {
	expected := "DriverName=sqlserver, Host=127.0.0.1, Port=1433, User=sa, Password=<REMOVED>, Database=master, TimeoutQuery=300"
	got := NewConfig("sqlserver").String()
	assert.Equal(t, expected, got)
}
