# Команды для работы с docker

Запуск докер-контейнеров

```bash
docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=Password123" -e "MSSQL_COLLATION=SQL_Latin1_General_CP1251_CI_AS" -p 1401:1433 --name sqlserver-test -d mcr.microsoft.com/mssql/server:2022-latest

```
Показать контейнеры

>docker container ls

Подключиться к контейнеру

```bash
docker container exec -it sqlserver-test bash
$ /opt/mssql-tools/bin/sqlcmd -S localhost -U SA -P Password123
$ exit
```

Stop all the containers:

```bash
docker stop $(docker ps -a -q)
docker stop sqlserver-test
```

Remove all the containers
>docker rm $(docker ps -a -q)
