package main

import (
	//"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"rbot/controllers"
	"rbot/models"
)

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		log.Println(err)
	} else {
		d := funcs.Db{Session: session}
		log.Println("Connected To MongoDb")
		rf := d.CreateAdmin(map[string]string{"username": "akshitgrover", "password": "1516"})
		log.Println(rf)
		rf = d.AdminAuth(map[string]string{"username": "akshit", "password": "1516"})
		log.Println(rf)
		rf, flag := d.CreateEvent(models.Event{Name: "test event", Fields: []string{"name", "email"}})
		log.Println(rf, flag)
		rf = d.InsertRecord("testeventMarch", models.Record{"name": "akshitgrover", "email": "a@a.com"})
		log.Println(rf)
		go http.HandleFunc("/event", d.InsertApi)
		http.ListenAndServe(":8080", nil)
	}
}
