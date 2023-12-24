package models

type App struct {
	Id     string `bson:"-"`
	Name   string `bson:"name"`
	Secret string `bson:"secret"`
}
