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
	cd dte
	go vet
	go fmt
	golangci-lint run --fix ./...

	cd ../dtegorm
	go vet
	go fmt
	golangci-lint run --fix ./...

down:
	docker compose down --remove-orphans

tag:
	@if [ -z "$(TAG)" ]; then echo "TAG variable is required."; exit 1; fi
	git tag $(TAG)
	git push origin $(TAG)

	git tag dte/$(TAG)
	git push origin dte/$(TAG)

	git tag dtegorm/$(TAG)
	git push origin dtegorm/$(TAG)
