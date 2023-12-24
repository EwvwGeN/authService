package models

type User struct {
	Id       string `bson:"-"`
	Email    string `bson:"email"`
	PassHash []byte `bson:"passHash"`
	IsAdmin  bool   `bson:"admin"`
}