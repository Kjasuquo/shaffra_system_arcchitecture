run:
	go run cmd/main.go

mock-repo:
	mockgen -source=internal/repository/repository.go -destination=internal/repository/mocks/repository_mock.go -package=mocks

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=donotshare -e POSTGRES_DB=shaffra -d postgres:latest

docker-up:
	docker compose up

test: mock-repo
	go test ./... -coverprofile=coverage.out

buggy-db:
	docker run --name postgres -p 5432:5432 \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=donotshare \
	-e POSTGRES_DB=test \
	-d postgres:latest