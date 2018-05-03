package models

type Record map[string]string

type Event struct {
	Id     string
	Name   string
	Fields []string
}
