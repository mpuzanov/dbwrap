package mssql_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/mpuzanov/dbwrap"

	_ "github.com/microsoft/go-mssqldb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Person человек
type Person struct {
	LastName    string     `faker:"lang=rus" db:"last_name"`
	Birthdate   *time.Time `db:"birthdate"`
	Salary      *float64   `faker:"amount"`
	IsOwnerFlat *bool      `db:"is_owner_flat"` // признак владельца помещения
	Email       string     `faker:"email"`
	CreatedAt   time.Time  `faker:"-" db:"created_at"`
}

var (
	dbName        = "db_test"
	tableName     = fmt.Sprintf("%s.dbo.people", dbName)
	password      = "Password123"
	port          = 1401
	ctxDefault, _ = context.WithTimeout(context.Background(), 5*time.Second)
)

// DBSuite структура для набора тестов с БД
type TestDBSuite struct {
	suite.Suite
	db *dbwrap.DBSQL // коннект к БД master
}

func TestDBNewConnect(t *testing.T) {
	config := dbwrap.NewConfig("sqlserver").WithDB("").WithPort(port)

	_, err := dbwrap.NewConnect(config)
	assert.ErrorIs(t, err, dbwrap.ErrBadConfigDB)

	config.WithPassword("12345").WithDB("master")
	_, err = dbwrap.NewConnect(config)
	assert.ErrorContains(t, err, "sqlx.Connect")

	config.WithPassword(password)
	db, err := dbwrap.NewConnect(config)
	assert.NoError(t, err)

	db.SetTimeout(100)
	assert.Equal(t, 100, db.Cfg.TimeoutQuery)
}

func TestTestDBSuite(t *testing.T) {
	suite.Run(t, &TestDBSuite{})
}

func (ts *TestDBSuite) SetupSuite() {

	config := dbwrap.NewConfig("sqlserver").WithPassword(password).WithDB("master").WithPort(port)
	db, err := dbwrap.NewConnect(config)
	if err != nil {
		ts.T().Fatalf("cannot connect db : %v", err)
	}
	ts.db = db
	setupDatabase(ts)
}

func (ts *TestDBSuite) TearDownSuite() {
	tearDownDatabase(ts)
}

func setupDatabase(ts *TestDBSuite) {
	ts.T().Log("setting up database")

	_, err := ts.db.DBX.ExecContext(ctxDefault, fmt.Sprintf(`DROP DATABASE IF EXISTS %s; CREATE DATABASE %s`, dbName, dbName))
	if err != nil {
		ts.FailNowf("unable to create database", err.Error())
	}
	ts.T().Logf("База [%s] создана\n", dbName)

	query := fmt.Sprintf(`CREATE TABLE %s (
		last_name varchar(50) PRIMARY KEY,
		birthdate datetime,
		salary decimal(15,2),
		is_owner_flat bit,
		email varchar(100) UNIQUE,
		created_at datetime NOT NULL DEFAULT current_timestamp
	)`, tableName)

	_, err = ts.db.DBX.ExecContext(ctxDefault, query)
	if err != nil {
		ts.FailNowf("unable to create table", err.Error())
	}
	ts.T().Logf("Таблица [%s] создана\n", tableName)

}

func tearDownDatabase(ts *TestDBSuite) {
	ts.T().Log("tearing down database")

	_, err := ts.db.DBX.ExecContext(ctxDefault, fmt.Sprintf(`DROP TABLE %s`, tableName))
	if err != nil {
		ts.FailNowf("unable to drop table", err.Error())
	}

	_, err = ts.db.DBX.ExecContext(ctxDefault, fmt.Sprintf(`DROP DATABASE %s`, dbName))
	if err != nil {
		ts.FailNowf("unable to drop database", err.Error())
	}

	err = ts.db.Close()
	if err != nil {
		ts.FailNowf("unable to close database", err.Error())
	}
}

func (ts *TestDBSuite) TestData1() {

	dataInsert := map[string]any{
		"LastName": "Иванов",
		"Email":    "ivan@example.com",
		//"is_owner_flat": true,
	}
	//ts.T().Logf("dataInsert: %#v", dataInsert)

	ts.Suite.Run("insert Test", func() {
		query := fmt.Sprintf(`INSERT INTO %s (last_name, Email) VALUES (:LastName, :Email)`, tableName)
		count, err := ts.db.NamedExecContext(ctxDefault, query, dataInsert)
		ts.NoError(err)
		ts.Equal(int64(1), count)
	})

	//===========================================================
	ts.Suite.Run("select Test", func() {
		query := fmt.Sprintf(`select * from %s`, tableName)
		var people []Person
		err := ts.db.SelectContext(ctxDefault, &people, query)
		ts.NoError(err)
		ts.Len(people, 1)
		ts.Equal("Иванов", people[0].LastName)
		//ts.T().Logf("%+v", people)

		// ===========================================================
		ts.T().Log("select Test where")
		query = fmt.Sprintf(`select last_name, email, created_at from %s where last_name=:Name`, tableName)
		var people2 []Person
		err = ts.db.NamedSelectContext(ctxDefault, &people2, query, map[string]any{"Name": "Иванов"})
		ts.NoError(err)
		ts.Len(people2, 1)
		ts.Equal("ivan@example.com", people2[0].Email)
		//ts.T().Logf("%+v", people2)
	})

	//===========================================================
	ts.Suite.Run("Get Test", func() {
		query := fmt.Sprintf(`select count(*) from %s`, tableName)
		var count int
		err := ts.db.GetContext(ctxDefault, &count, query)
		ts.NoError(err)
		ts.Equal(1, count)
		//ts.T().Logf("count: %d", count)

		// ===========================================================
		ts.T().Log("Get Test where")
		query = fmt.Sprintf(`select last_name, email, created_at from %s where last_name=:Name`, tableName)
		person := Person{}
		err = ts.db.NamedGetContext(ctxDefault, &person, query, map[string]any{"Name": "Иванов"})
		ts.NoError(err)
		//ts.T().Logf("person: %+v", person)
		ts.Equal("ivan@example.com", person.Email)
	})

	//===========================================================
	ts.Suite.Run("Get Map", func() {
		query := fmt.Sprintf(`select last_name, email, created_at from %s where last_name=:Name`, tableName)
		person, err := ts.db.NamedGetMapContext(ctxDefault, query, map[string]any{"Name": "Иванов"})
		ts.NoError(err)
		ts.Equal("ivan@example.com", person["email"])
		//ts.T().Logf("person: %+v", person)
	})

	//===========================================================
	ts.Suite.Run("update Test", func() {
		query := fmt.Sprintf(`UPDATE %s SET email=:Email where last_name=:Name`, tableName)
		dataUpdate := map[string]any{
			"Name":  "Иванов",
			"Email": "email_update@example.com",
		}
		//ts.T().Logf("dataUpdate: %+v", dataUpdate)
		count, err := ts.db.NamedExecContext(ctxDefault, query, dataUpdate)
		ts.Require().NoError(err)
		ts.Equal(int64(1), count)
	})

	//===========================================================
	ts.Suite.Run("delete Test", func() {
		query := fmt.Sprintf(`delete from %s where last_name=:Name`, tableName)
		dataDelete := map[string]any{
			"Name": "Иванов",
		}
		//ts.T().Logf("dataDelete: %+v", dataDelete)
		count, err := ts.db.NamedExecContext(ctxDefault, query, dataDelete)
		ts.Require().NoError(err)
		ts.Equal(int64(1), count, "удалено")
	})
}

func (ts *TestDBSuite) TestDataSelect() {

	// batch insert with maps
	dtIns := []map[string]any{
		{"LastName": "Сидоров", "Email": "sidr@example.com", "Birthdate": time.Date(2000, 2, 21, 0, 0, 0, 0, time.UTC)},
		{"LastName": "Кузнецов", "Email": "kuz@gmail.com", "Birthdate": nil},
		{"LastName": "Петров", "Email": "peter@example.com", "Birthdate": nil},
	}

	//ts.T().Logf("данные для вставки: %+v", dtIns)
	query := fmt.Sprintf(`INSERT INTO %s (last_name, Email, Birthdate) VALUES (:LastName, :Email, :Birthdate)`, tableName)
	_, err := ts.db.NamedExecContext(ctxDefault, query, dtIns)
	ts.NoError(err)

	//==================================================
	query = fmt.Sprintf(`select * from %s`, tableName)
	resultMap, err := ts.db.SelectMapsContext(ctxDefault, query)
	ts.NoError(err)
	ts.Len(resultMap, 3)
	//ts.T().Logf("получим в map: %+v", resultMap)

	var resultSlice []Person
	err = ts.db.SelectContext(ctxDefault, &resultSlice, query)
	ts.Require().NoError(err)
	ts.Len(resultSlice, 3)
	//ts.T().Logf("получим в slice: %+v", resultSlice)
	//==================================================

	query = fmt.Sprintf(`select * from %s where last_name=:Name`, tableName)
	p := []map[string]any{{"Name": "Кузнецов"}}
	resultMap2, err := ts.db.NamedSelectMapsContext(ctxDefault, query, p)
	ts.NoError(err)
	ts.Len(resultMap2, 1)
	//ts.T().Logf("получим в map: %+v", resultMap2)
	ts.Equal("Кузнецов", resultMap2[0]["last_name"])
}

func (ts *TestDBSuite) TestEmptyData() {
	ts.Suite.Run("Get EmptyData", func() {
		query := fmt.Sprintf(`select last_name, email, created_at from %s where last_name=:Name`, tableName)
		person := Person{}
		err := ts.db.NamedGetContext(ctxDefault, &person, query, map[string]any{"Name": "Большунов"})
		ts.ErrorIs(err, sql.ErrNoRows)
		//ts.T().Log(err)
		//ts.T().Logf("person: %+v", person)
	})
}
