# Обёртка для SQLX

Протестировано для MSSQL, PostgreSQL, MySQL, SQLite

## Драйвера БД

Для PostgreSQL:

>go get github.com/lib/pq

Для MSSQL Server:

>go get github.com/denisenkom/go-mssqldb

Для MySQL:

>github.com/go-sql-driver/mysql

Для SQLite:

>go get github.com/mattn/go-sqlite3

## Запуск докер-контейнеров

```bash
docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=Password123" -e "MSSQL_COLLATION=SQL_Latin1_General_CP1251_CI_AS" -p 1401:1433 --name sqlserver-test -d mcr.microsoft.com/mssql/server:2022-latest

docker container ls

docker container exec -it sqlserver-test bash
$ /opt/mssql-tools/bin/sqlcmd -S localhost -U SA -P Password123
$ exit

docker stop sqlserver-test

```
