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

func (d Db) InsertRecord(name string, s models.Record) string {
	var fields models.Event
	err := d.Session.DB("rbot").C("event").Find(bson.M{"id": name}).Select(bson.M{"fields": 1}).One(&fields)
	if err != nil && err.Error() == "not found" {
		return "No Such Event Exist"
	}
	if err != nil {
		log.Println(err)
		return "Something Went Wrong"
	}
	values := make(map[string]string)
	for _, v := range fields.Fields {
		values[v] = s[v]
	}
	err = d.Session.DB("rbot").C("event").Update(bson.M{"id": name}, bson.M{"$push": bson.M{"values": values}})
	if err != nil {
		log.Println(err)
		return "Something Went Wrong"
	}
	return "Record Posted Successfully"
}
