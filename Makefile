.SILENT:

export TEST_CONTAINER_NAME=test_project
export TEST_POSTGRES_PASSWORD=qwerty
export TEST_CLIENT_NAME=test-client

test.integration:
	docker run --rm -d -p 5444:5432 --name $$TEST_CONTAINER_NAME -e POSTGRES_PASSWORD=$$TEST_POSTGRES_PASSWORD postgres
	sleep 1
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5444/postgres?sslmode=disable' up
	go test -v ./tests/
	docker stop $$TEST_CONTAINER_NAME

unit-test:
	docker run --rm -d -p 5444:5432 --name $$TEST_CONTAINER_NAME -e POSTGRES_PASSWORD=$$TEST_POSTGRES_PASSWORD postgres
	sleep 1
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5444/postgres?sslmode=disable' up
	go test password_keeper/internal/server/handler  password_keeper/internal/server/service
	docker stop $$TEST_CONTAINER_NAME

docker-stop:
	docker stop $$TEST_CONTAINER_NAME