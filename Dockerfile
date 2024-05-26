FROM golang
ARG PORT=8080
ARG SERVICE=app
ARG DB_HOST=localhost
ARG DB_PORT=5432
ARG DB_DATABASE=postgres
ARG DB_USERNAME=postgres
ARG DB_PASSWORD=qwerty

# Устанавливаем рабочую директорию
WORKDIR /app

# Определяем переменные окружения
ENV PORT=${PORT}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_DATABASE=${DB_DATABASE}
ENV DB_USERNAME=${DB_USERNAME}
ENV DB_PASSWORD=${DB_PASSWORD}


COPY go.mod .
RUN go mod download

EXPOSE ${PORT}

COPY private.rsa private.rsa

COPY . .

RUN go build -o app ./cmd/${SERVICE}
CMD ["./app"]


