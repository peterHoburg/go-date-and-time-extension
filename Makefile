.ONESHELL:

test: down
	docker compose run -d --remove-orphans -p 5432:5432 postgres
	go test -race ./...
	cd dte/gorm
	go test -race ./...


down:
	docker compose down --remove-orphans
