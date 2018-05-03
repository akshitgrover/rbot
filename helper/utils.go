package helper

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"rbot/controllers"
	"rbot/models"
	"time"
)

type MSession struct {
	Ses *mgo.Session
}

var UserSession = make(map[string]string)
var UserState = make(map[string]int)
var UserTimer = make(map[string]time.Time)

type Config struct {
	Token string            `json:"api_token"`
	Urls  map[string]string `json:"mongo_url"`
}

type Admin struct {
	Admin map[string]bool
}

func ReadConfigJson() Config {
	f, _ := os.Open("./config/config.json")
	st, _ := f.Stat()
	s := st.Size()
	defer f.Close()
	b := make([]byte, s)
	f.Read(b)
	var data Config
	json.Unmarshal(b, &data)
	return data
}

func AddURL(event_id string, url string) {
	data := ReadConfigJson()
	data.Urls[event_id] = url
	b, _ := json.Marshal(data)
	f, _ := os.Open("./config.config.json")
	f.Write(b)
	f.Close()
}

func MongoConnect(event_id string) (*mgo.Session, bool) {
	data := ReadConfigJson()
	url := data.Urls[event_id]
	if url == "" {
		return nil, false
	}
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, false
	}
	return session, true
}

func ReadAdminJson() Admin {
	f, _ := os.Open("./config/admin.json")
	st, _ := f.Stat()
	s := st.Size()
	defer f.Close()
	b := make([]byte, s)
	f.Read(b)
	var data Admin
	json.Unmarshal(b, &data)
	return data
}

func CheckAdmin(username string) bool {
	data := ReadAdminJson()
	if data.Admin[username] {
		return true
	}
	return false
}

func AddAdmin(username string) {
	data := ReadAdminJson()
	data.Admin[username] = true
	b, _ := json.Marshal(data)
	f, _ := os.Open("./config/admin.json")
	f.Write(b)
	f.Close()
}

func CheckSession(username string) bool {
	if UserSession[username] == "" {
		return false
	}
	return true
}

func AddSession(username string, event_id string) {
	UserSession[username] = event_id
}

func AddTimer(username string) {
	UserTimer[username] = time.Now()
}

func CheckTimer(username string) bool {
	log.Println(UserTimer[username].Sub(time.Now()))
	if UserTimer[username].IsZero() || time.Now().Sub(UserTimer[username]) > 5*time.Second {
		return false
	}
	return true
}

func GetEventDbUrl(event_id string) string {
	data := ReadConfigJson()
	return data.Urls[event_id]
}

func GetState(username string) int {
	return UserState[username]
}

func SetState(username string, state int) {
	UserState[username] = state
}

func (ses MSession) CheckEventValid(eventid string) bool {
	var err error
	var ev models.Event
	funcs.FindOne(ses.Ses, "event", bson.M{"id": bson.M{"$regex": "^" + eventid + "$", "$options": "i"}}, &ev, &err)
	if err != nil && err.Error() == "not found" {
		return false
	}
	return true
}

func DelMongoSession(eventid string) {
	timer := time.NewTimer(5 * time.Second)
	<-timer.C
	activeSessions[eventid].Close()
	activeSessions[eventid] = nil
}
