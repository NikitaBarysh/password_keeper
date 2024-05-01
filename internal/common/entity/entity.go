package entity

type User struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login" required:"true"`
	Password string `json:"password" required:"true"`
}

type Data struct {
	UserId    int    `json:"id,omitempty"`
	Payload   []byte `json:"payload"`
	EventType string `json:"event_type"`
}

type WebDataMSG struct {
	Action    string `json:"action"`
	EventType string `json:"event_type"`
	Payload   []byte `json:"payload,omitempty"`
}
