package models

type Record map[string]string

type Event struct {
	Id     string
	Name   string
	Fields []string
	values []Record
}

type Admin struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}
