package helper

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"rbot/controllers"
	"strconv"
	"strings"
)

var activeSessions = make(map[string]*mgo.Session)

var Instants = map[string]bool{
	"get event":    true,
	"add event":    true,
	"remove event": true,
	"make admin":   true,
}

var StateInstants = map[string]int{
	"get event":    2,
	"add event":    3,
	"remove event": 4,
	"make admin":   5,
}

var ReqActiveEidInstants = map[string]bool{
	"get event":    true,
	"remove event": true,
}

var InverStateInstants = map[int]string{
	2: "get event",
	3: "add event",
	4: "remove event",
	5: "make admin",
}

var Texts = map[string]string{
	"default":  "Sorry! I Can't Understand Your Text",
	"hello":    "Hello, How May I Help You!",
	"timerexp": "Would You Like To Continue With The Event Id Provided Earlier",
	"dberror":  "Couldn't Connect To Database At The Moment!",
	"1+":       "Before Proceeding, I Couldn't Find Out An Active Session, Help Me By Providing An Event Id",
	"1--":      "Sorry! Unfortunately Event Id Is Invalid! Try Again!",
	"1-+":      "Hooray! Your Session Is Active!",
	"2+":       "Cool, Would Like To Mention Some EventId Please!",
}

func StateTwo(ses MSession, eventid string) string {
	if !ses.CheckEventValid(eventid) {
		return Texts["1--"]
	}
	data := ReadConfigJson()
	url := data.Urls[eventid]
	flag := strings.Split(url, "/")
	dbname := flag[len(flag)-1]
	var session *mgo.Session
	var err error
	if activeSessions[eventid] == nil {
		log.Println("Db Connect Call: " + eventid)
		session, err = mgo.Dial(url)
		if err != nil {
			log.Println(err)
			return Texts["dberror"]
		}
		activeSessions[eventid] = session
		go DelMongoSession(eventid)
	} else {
		session = activeSessions[eventid]
	}
	var ev []map[string]string
	funcs.FindAll(session, "records", bson.M{}, &ev, &err, dbname)
	var res = "Event Id: " + eventid + "\r\n" + "Count: " + strconv.Itoa(len(ev))
	log.Println(ses.CheckEventValid(eventid))
	return res
}

func StateSeven(ses MSession, username string, res string) string {
	if res == "yes" {
		AddTimer(username)
		state := TimerState[username]
		SetState(username, state)
		return RecvResponse(ses, username, UserSession[username])
	} else {
		SetState(username, 9)
		return Texts["2+"]
	}
}
