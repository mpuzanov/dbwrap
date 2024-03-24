.PHONY: lint test test-all test-mysql test-postgres test-mssql
.DEFAULT_GOAL = test

test:
	go test -v -count=1 ./tests/sqlite
test-mssql:
	go test -v -count=1 ./tests/mssql
test-postgres:
	go test -v -count=1 ./tests/postgres
test-mysql:
	go test -v -count=1 ./tests/mysql

test-all:
	go test -v -count=1 ./...

lint: 
	@gofmt -s -w ./...
	@go vet ./...
	@golint ./...

docker-mysql:
	docker-compose -f ./tests/mysql/docker-compose-mysql.yml up -d

docker-postgres:
	docker-compose -f ./tests/postgres/docker-compose-postgres.yml up -d

docker-mssql:
	docker-compose -f ./tests/mssql/docker-compose-mssql.yml up -d

docker-stop:
	@docker stop $(shell docker ps -a -q)
	@docker rm $(shell docker ps -a -q)

#docker-rm:
# 	for containers	
#	@docker rm $(shell docker ps -a -q)
# 	for volumes
#	@docker volume rm $(shell docker volume ls -q)	
# 	for images	
#	@docker rmi $(shell docker images -a -q)  