package funcs

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"rbot/models"
	"strings"
	"time"
)

func (d Db) CreateEvent(s models.Event) (string, string) {
	var rf models.Event
	var name string
	name = strings.ToLower(s.Name)
	flag := strings.Split(s.Name, " ")
	name = strings.Join(flag, "")
	var id = name + time.Now().Month().String()
	event := models.Event{Id: id, Name: s.Name, Fields: s.Fields, Values: make([]models.Record, 0)}
	err := d.Session.DB("rbot").C("event").Find(bson.M{"id": id}).One(&rf)
	if err != nil && err.Error() != "not found" {
		log.Println(err)
		return "Something Went Wrong", ""
	}
	if rf.Id != "" {
		return "Event Already Exists", ""
	}
	err = d.Session.DB("rbot").C("event").Insert(event)
	if err != nil {
		log.Println(err)
		return "Something Went Wrong", ""
	}
	return "Event Created", ""
}
