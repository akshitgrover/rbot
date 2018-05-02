package funcs

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"rbot/models"
)

type Db struct {
	Session *mgo.Session
}

func (d Db) CreateAdmin(s map[string]string) string {
	var rf models.Admin
	admin := models.Admin{Username: s["username"], Password: s["password"]}
	err := d.Session.DB("rbot").C("admin").Find(bson.M{"username": s["username"]}).One(&rf)
	if err != nil && err.Error() != "not found" {
		log.Println(err)
		return "Something Went Wrong"
	}
	if rf.Username != "" {
		return "Username Already Exists"
	}
	bs, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.MinCost)
	admin.Password = string(bs)
	err = d.Session.DB("rbot").C("admin").Insert(admin)
	if err != nil {
		log.Println(err)
		return "Something Went Wrong"
	}
	return "Inserted"
}

func (d Db) AdminAuth(s map[string]string) string {
	var rf models.Admin
	err := d.Session.DB("rbot").C("admin").Find(bson.M{"username": s["username"]}).One(&rf)
	if err != nil && err.Error() == "not found" {
		return "Unauthenticated"
	}
	if err != nil {
		return "Something Went Wrong"
	}
	return "Authentication Successful"
}
