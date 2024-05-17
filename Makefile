include dev.env

# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)


run:
	go run apis/blog-api/main.go

debug:
	curl http://localhost:8001

# ==============================================================================
# Define dependencies
GOLANG          := golang:1.21.3
POSTGRES        := postgres:15.4
DB_DSN 			:= postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}/${DB_NAME}?sslmode=disable


env:
	@echo ${DB_DSN}
## docker
docker/db/run:
	@echo "Running postgres db"
	docker run --name ${DB_NAME} -e POSTGRES_SECRET=${DB_PASSWORD} -p 5431:5432 ${POSTGRES}

docker/db/stop:
	@echo "stoping postgres db"
	docker stop ${DB_NAME} 

docker/db/start:
	@echo "starting postgres db"
	docker start ${DB_NAME}


## database 
db/psql:
	@psql ${DB_DSN} 
db/pgcli:
	@pgcli ${DB_DSN}

db/migrate/up:
	@migrate -path business/data/dbmigrate/migrations -database "${DB_DSN}" -verbose up
db/migrate/up/latest:
	@migrate -path business/data/dbmigrate/migrations -database "${DB_DSN}" -verbose up 1
	
db/migrate/new:
	@migrate create -ext sql -dir business/data/dbmigrate/migrations -seq ${name}

db/migrate/down:
	@migrate -path business/data/dbmigrate/migrations -database "${DB_DSN}" -verbose down
db/migrate/down/latest:
	@migrate -path business/data/dbmigrate/migrations -database "${DB_DSN}" -verbose down 1

db/init:
	@echo "Seeding database"
	go run apis/tooling/initdb/main.go
	
tidy:
	go mod tidy
	go mod vendor

test-race:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...
	
test: test-only lint vuln-check