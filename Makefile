include dev.env

# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

run:
	go run apis/blog-api/main.go | go run apis/tooling/logfmt/main.go
run/migrate:
	go run apis/tooling/admin/main.go migrate
run/init:
	go run apis/tooling/admin/main.go migrate,seed

debug:
	curl http://localhost:8001

# ==============================================================================
# Define dependencies
GOLANG          := golang:1.21.3
POSTGRES        := postgres:15.4
DB_DSN 			:= postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}/${DB_NAME}?sslmode=disable

BASE_IMAGE_NAME := localhost/islamghany
VERSION         := 0.0.1
APP_NAME 		:= blog-api
NAMESPACE       := blog-system
BLOG_IMAGE_NAME := ${BASE_IMAGE_NAME}/${APP_NAME}:${VERSION}


env:
	@echo ${DB_DSN}
## docker

docker/build:
	@echo "Building docker image"
	docker build \
		-f infra/docker/dockerfile.blog \
		-t ${BLOG_IMAGE_NAME} \
		--build-arg BUILD_REF=$(BLOG_IMAGE_NAME) \
		--build-arg BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		.
docker/rmi:
	@echo "Removing docker image"
	docker rmi ${BLOG_IMAGE_NAME}

docker/run:
	@echo "Running docker image"
	docker run -p 8000:8000 ${BLOG_IMAGE_NAME}

docker/stop:
	@echo "Stopping docker image"
	docker stop ${BLOG_IMAGE_NAME}
docker/remove:
	@echo "Removing docker image"
	docker rm ${BLOG_IMAGE_NAME}

docker/db/run:
	@echo "Running postgres db"
	docker run --name ${DB_NAME} -e POSTGRES_SECRET=${DB_PASSWORD} -p 5431:5432 ${POSTGRES}

docker/db/stop:
	@echo "stoping postgres db"
	docker stop ${DB_NAME} 

docker/db/start:
	@echo "starting postgres db"
	docker start ${DB_NAME}

## k8s
dev-up:
	@echo "Making minikube look at the local docker daemon"
	eval $(minikube -p minikube docker-env)
	@echo "Building docker image for blog api in minikube"
	docker build \
		-f infra/docker/dockerfile.blog \
		-t ${BLOG_IMAGE_NAME} \
		--build-arg BUILD_REF=$(BLOG_IMAGE_NAME) \
		--build-arg BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		.
	@echo "pull postgres:15.4 to minikube"
	minikube image pull ${POSTGRES}

dev-apply:
	kustomize build infra/k8s/blog | kubectl apply -f -
	kustomize build infra/k8s/postgres | kubectl apply -f -
dev-down:
	kustomize build infra/k8s/blog | kubectl delete -f -
	kustomize build infra/k8s/postgres | kubectl delete -f -

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=blog --all-containers=true -f --tail=100 --max-log-requests=6 | go run apis/tooling/logfmt/main.go -service=$(BLOG_IMAGE_NAME)

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


## Metrics and Tracing

metrics-view-sc:
	expvarmon -ports="localhost:8001" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

metrics-view:
	expvarmon -ports="localhost:4020" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"


##  module setup

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all


## test
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