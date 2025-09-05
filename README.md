# Обёртка для SQLX

Добавлены методы:

```golang
Exec(query string, args ...interface{}) (int64, error)
NamedExec(query string, arg interface{}) (int64, error)
Select(dest interface{}, query string, args ...interface{}) error
NamedSelect(dest interface{}, query string, arg interface{}) error
SelectMaps(query string, args ...interface{}) (ret []map[string]interface{}, err error)
NamedSelectMaps(query string, arg interface{}) (ret []map[string]interface{}, err error)
Get(dest interface{}, query string, args ...interface{}) error
NamedGet(dest interface{}, query string, arg interface{}) error
GetMap(query string, args ...interface{}) (ret map[string]interface{}, err error)
NamedGetMap(query string, arg interface{}) (ret map[string]interface{}, err error)
```

Протестировано для MSSQL, PostgreSQL, MySQL, SQLite

Установка `go get github.com/mpuzanov/dbwrap`

## Примеры

```golang

    config := dbwrap.NewConfig("sqlserver").WithPassword(password).WithDB("master").WithPort(port)
    db, err := dbwrap.NewConnect(config)
    if err != nil {
        panic(err)
    }
    log.Println("cfg.DB", config.String())


    query := fmt.Sprintf(`select last_name, email, created_at from %s where last_name=:Name`, tableName)
    person := Person{}
    err = db.NamedGet(&person, query, map[string]interface{}{"Name": "Иванов"})
```


## Драйвера БД

Для PostgreSQL:

>go get github.com/lib/pq

Для MSSQL Server:

>go get github.com/microsoft/go-mssqldb

Для MySQL:

>github.com/go-sql-driver/mysql

Для SQLite:

>go get github.com/mattn/go-sqlite3


