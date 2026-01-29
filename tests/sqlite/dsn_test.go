package sqlite_test

import (
	"testing"

	"github.com/mpuzanov/dbwrap"
)

func TestNewConnectDSN(t *testing.T) {
	db, err := dbwrap.NewConnectDSN("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("cannot connect db: %v", err)
	}

	query := `CREATE TABLE user (name varchar(50) PRIMARY KEY, age int, email varchar(100));`
	_, err = db.DBX.ExecContext(ctxDefault, query)
	if err != nil {
		t.Errorf("unable to create table [user] %s", err.Error())
	}

	_, err = db.DBX.ExecContext(ctxDefault, "DROP TABLE user;")
	if err != nil {
		t.Errorf("unable to drop table %s", err.Error())
	}

	err = db.Close()
	if err != nil {
		t.Errorf("unable to close database %s", err.Error())
	}
}
