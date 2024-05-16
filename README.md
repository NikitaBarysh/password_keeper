Сборка клиента:
go build -ldflags="-X main.Version=v1.0.0 -X 'main.BuildTime=$(date +'%Y/%m/%d %H:%M:%S')'"


