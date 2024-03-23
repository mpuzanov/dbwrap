package mssql_test

import (
	"dbwrap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBNew(t *testing.T) {
	config := dbwrap.NewConfig("sqlserver").WithDB("").WithPort(1401)

	_, err := dbwrap.New(config)
	assert.ErrorIs(t, err, dbwrap.ErrBadConfigDB)

	config.WithPassword("12345").WithDB("master")
	_, err = dbwrap.New(config)
	assert.ErrorContains(t, err, "sqlx.Connect")

	config.WithPassword("Password123")
	db, err := dbwrap.New(config)
	assert.NoError(t, err)

	db.SetTimeout(100)
	assert.Equal(t, 100, db.Cfg.TimeoutQuery)
}
