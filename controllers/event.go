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
	var err error
	flagch := make(chan int)
	name = strings.ToLower(s.Name)
	flag := strings.Split(s.Name, " ")
	name = strings.Join(flag, "")
	var id = name + time.Now().Month().String()
	event := models.Event{Id: id, Name: s.Name, Fields: s.Fields, Values: make([]models.Record, 0)}
	go func() {
		FindOne(d, "event", bson.M{"id": id}, &rf, &err)
		flagch <- 1
	}()
	if <-flagch == 1 {
		if err != nil && err.Error() != "not found" {
			log.Println(err)
			return "Something Went Wrong", ""
		}
		if rf.Id != "" {
			return "Event Already Exists", ""
		}
		go func() {
			Insert(d, "event", event, &err)
			log.Println("TEST")
			log.Println(err)
			flagch <- 2
		}()
		for i := range flagch {
			if i == 2 {
				if err != nil {
					log.Println(err)
					return "Something Went Wrong", ""
				}
				return "Event Created", ""
				break
			}
		}
	}
	return "", ""
}

func (d Db) InsertRecord(name string, s models.Record) string {
	var fields models.Event
	var err error
	flagch := make(chan int)
	go func() {
		FindOneProject(d, "event", bson.M{"id": name}, bson.M{"fields": 1}, &fields, &err)
		flagch <- 1
	}()
	if <-flagch == 1 {
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
		go func() {
			Update(d, "event", bson.M{"id": name}, bson.M{"$push": bson.M{"values": values}}, &err)
			flagch <- 2
		}()
		for i := range flagch {
			if i == 2 {
				if err != nil {
					log.Println(err)
					return "Something Went Wrong"
				}
				return "Record Posted Successfully"
			}
		}
	}
	return ""
}

func (d Db) CheckAdmin(username string) int {
	var rf models.Admin
	var err error
	flagch := make(chan int)
	go func() {
		FindOneAdmin(d, "admin", bson.M{username: username}, &rf, &err)
		flagch <- 1
	}()
	if <-flagch == 1 {
		if err.Error() == "not found" {
			return 2
		}
		return 1
	}
	return 2
}
