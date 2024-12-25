.ONESHELL:

test: down
	docker compose run -d --remove-orphans -p 5432:5432 postgres
	cd dte
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	cd ../dtegorm
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

build:
	cd dte
	go build -v ./...
	cd ../dtegorm
	go build -v ./...

lint:
	go vet
	go fmt
	golangci-lint run --fix ./...

down:
	docker compose down --remove-orphans

