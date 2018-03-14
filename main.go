package main

import (
	//"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	//"net/http"
	"log"
	"rbot/controllers"
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
	}
}
