package models

type Message struct {
	Subject string `json:"subject"`
	EmailTo string `json:"email_to"`
	Body    []byte `json:"body"`
}