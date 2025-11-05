# Обёртка для SQLX

Добавлены методы:

```golang
ExecContext(ctx context.Context, query string, args ...any) (int64, error)
NamedExecContext(ctx context.Context,query string, arg any) (int64, error)
SelectContext(ctx context.Context,dest any, query string, args ...any) error
NamedSelectContext(ctx context.Context,dest any, query string, arg any) error
SelectMapsContext(ctx context.Context,query string, args ...any) (ret []map[string]any, err error)
NamedSelectMapsContext(ctx context.Context,query string, arg any) (ret []map[string]any, err error)
GetContext(ctx context.Context,dest any, query string, args ...any) error
NamedGetContext(ctx context.Context,dest any, query string, arg iany) error
GetMapContext(ctx context.Context,query string, args ...any) (ret map[string]any, err error)
NamedGetMapContext(ctx context.Context,query string, arg any) (ret map[string]any, err error)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

    query := fmt.Sprintf(`select last_name, email, created_at from %s where last_name=:Name`, tableName)
    person := Person{}
    err = db.NamedGetContext(ctx, &person, query, map[string]any{"Name": "Иванов"})
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
