.SILENT:

TEST_CONTAINER_NAME=test-db
TEST_SERVER_NAME=test-server
TEST_POSTGRES_PASSWORD=qwerty
NETWORK_NAME=my-network


# если нужно запустить тесты локально, то в тесте в интеграционном тесте TestSignUp поменять ожидаемый id с 6 на 2
test.integration:
	echo "Starting PostgreSQL container for integration tests"
	docker run --rm -d \
 		  -p 5432:5432 \
 		  --name $(TEST_CONTAINER_NAME) \
 		  -e POSTGRES_PASSWORD=$(TEST_POSTGRES_PASSWORD) \
 		  postgres

	echo "Checking PostgreSQL availability"
	until docker exec $(TEST_CONTAINER_NAME) pg_isready -U postgres; do \
    	echo "Waiting for PostgreSQL to be ready..."; \
    	sleep 2; \
      done

	echo "Running migrations"
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up

	# Проверка завершения миграций
	if [ $$? -ne 0 ]; then \
      	echo "Migrations failed"; \
    	docker stop $(TEST_CONTAINER_NAME); \
    	exit 1; \
      fi

	echo "Running tests"
	go test -v ./tests/

	echo "Stopping Postgres container"
	docker stop $(TEST_CONTAINER_NAME)


unit-test:
	echo "Starting PostgreSQL container for unit tests"
	docker run --rm -d \
		  -p 5432:5432 \
		  --name $(TEST_CONTAINER_NAME) \
		  -e POSTGRES_PASSWORD=$(TEST_POSTGRES_PASSWORD) \
		  postgres

	until docker exec $(TEST_CONTAINER_NAME) pg_isready -U postgres; do \
  		echo "Waiting for PostgreSQL to be ready..."; \
  		sleep 2; \
  	  done

	echo "Running migrations"
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up

	# Проверка завершения миграций
	if [ $$? -ne 0 ]; then \
      	echo "Migrations failed"; \
      	docker stop $(TEST_CONTAINER_NAME); \
      	exit 1; \
      fi

	echo "Running unit tests"
	go test -v password_keeper/internal/server/handler  password_keeper/internal/server/service

	echo "Stopping PostgreSQL container"
	docker stop $(TEST_CONTAINER_NAME)

sender-test:
	@if ! docker network ls | grep -q $(NETWORK_NAME); then \
		echo "Creating network $(NETWORK_NAME)"; \
		docker network create $(NETWORK_NAME); \
	else \
		echo "Network $(NETWORK_NAME) already exists"; \
	fi

	# Запуск PostgreSQL контейнера
	echo "Starting PostgreSQL container for sender tests"
	docker run --rm -d \
		--network $(NETWORK_NAME) \
		--name $(TEST_CONTAINER_NAME) \
		-e POSTGRES_PASSWORD=$(TEST_POSTGRES_PASSWORD) \
		postgres

	# Проверка доступности контейнера PostgreSQL
	echo "Checking PostgreSQL availability"
	until docker exec $(TEST_CONTAINER_NAME) pg_isready -U postgres; do \
		echo "Waiting for PostgreSQL to be ready..."; \
		sleep 2; \
	done

	# Выполнение миграций
	echo "Running migrations"
	docker run --rm \
		--network $(NETWORK_NAME) \
		-v $(shell pwd)/schema:/schema \
		migrate/migrate \
		-path=/schema -database "postgres://postgres:$(TEST_POSTGRES_PASSWORD)@$(TEST_CONTAINER_NAME):5432/postgres?sslmode=disable" up

	# Проверка завершения миграций
	if [ $$? -ne 0 ]; then \
		echo "Migrations failed"; \
		docker stop $(TEST_CONTAINER_NAME); \
		docker network rm $(NETWORK_NAME); \
		exit 1; \
	fi

	# Запуск контейнера приложения
	echo "Starting application container"
	docker run --rm -d \
		--network $(NETWORK_NAME) \
		-e DB_HOST=$(TEST_CONTAINER_NAME) \
		-e DB_PORT=5432 \
		-e DB_DATABASE=postgres \
		-e DB_USERNAME=postgres \
		-e DB_PASSWORD=$(TEST_POSTGRES_PASSWORD) \
		-p 8000:8000 \
		--name $(TEST_SERVER_NAME) \
		password-keeper-server

	# Увеличиваем время ожидания для инициализации сервера
	echo "Waiting for application to start..."
	sleep 3

	# Проверка состояния контейнера приложения
	echo "Checking application container logs"
	docker logs $(TEST_SERVER_NAME)

	# Запуск тестов
	echo "Running tests"
	go test -v password_keeper/internal/client/sender

	# Остановка контейнеров после теста
	echo "Stopping containers"
	docker stop $(TEST_SERVER_NAME)
	docker stop $(TEST_CONTAINER_NAME)

	# Удаление сети
	@if docker network ls | grep -q $(NETWORK_NAME); then \
		echo "Removing network $(NETWORK_NAME)"; \
		docker network rm $(NETWORK_NAME); \
	fi

migration:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up

docker-stop:
	docker stop $(TEST_CONTAINER_NAME)