package models

type User struct {
	Id       string
	Email    string
	PassHash []byte
	IsAdmin  bool
}