package sender

//go:generate mockgen -source ${GOFILE} -destination mock.go -package ${GOPACKAGE}

type SendInterface interface {
	PostUserRequest(login, password, path string) ([]byte, error)
	PostDataRequest(data, eventType string) error
	GetDataRequest(eventType string) ([]byte, error)
	DeleteDataRequest(eventType string) error
	ConnectWs()
}
