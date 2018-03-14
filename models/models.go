package models

type Record map[string]string

type Event struct {
	Id     string
	Name   string
	Fields []string
	Values []Record
}

type Admin struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}
