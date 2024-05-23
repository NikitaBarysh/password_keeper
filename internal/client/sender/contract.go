// Package sender - пакет для взаимодействия с сервером
package sender

//go:generate mockgen -source ${GOFILE} -destination mock.go -package ${GOPACKAGE}

// SendInterface - интерфейс для мок тестирования пакета sender
type SendInterface interface {
	PostUserRequest(login, password, path string) error
	PostDataRequest(data, eventType string) error
	GetDataRequest(eventType string) ([]byte, error)
	DeleteDataRequest(eventType string) error
	ConnectWs() error
}
