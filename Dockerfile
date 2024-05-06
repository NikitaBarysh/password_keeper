# не работате пока что
FROM golang

WORKDIR .

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8000

COPY . .

RUN go build -o serv ./cmd/server/main.go


CMD ["./serv"]