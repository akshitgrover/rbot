package helper

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"rbot/controllers"
	"rbot/models"
	"strconv"
	"strings"
	"time"
)

var activeSessions = make(map[string]*mgo.Session)

var Instants = map[string]bool{
	"get event":    true,
	"add event":    true,
	"remove event": true,
	"make admin":   true,
	"add dburl":    true,
	"get apiurl":   true,
}

var StateInstants = map[string]int{
	"get event":    2,
	"add event":    3,
	"remove event": 4,
	"make admin":   5,
	"add dburl":    6,
	"get apiurl":   8,
}

var ReqActiveEidInstants = map[string]bool{
	"get event":    true,
	"remove event": true,
	"get apiurl":   true,
}

var InverStateInstants = map[int]string{
	2: "get event",
	4: "remove event",
	8: "get apiurl",
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
	"5+":       "Please Mention The Admin Username!",
	"5-":       "Hooray! Admin Added Successfully",
	"3+":       "Please Provide Event Structures",
	"3-":       "Error Occured While Entering An Event! Try Again!",
	"3--":      "Incomplete Details Provided",
	"3++":      "Event Added Successfully! Event Id: ",
	"6+":       "Please Provide Db Url!",
	"6-":       "Cannot Find A Event Add State! Please Help By Providing An EventId!",
	"6-+":      "Event Db Url Added Successfully!",
}

func StateTwo(ses MSession, eventid string) string {
	if !ses.CheckEventValid(eventid) {
		return Texts["1--"]
	}
	data := ReadConfigJson()
	url := data.Urls[eventid]
	if url == "" {
		return Texts["dberror"]
	}
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

func StateThree(ses MSession, username, message string) string {
	var ev models.Event
	ev.Fields = []string{}
	var flag = strings.Split(message, "\n")
	var count int
	for _, s := range flag {
		str_slice := strings.Split(s, ":")
		log.Println(str_slice[0])
		if strings.Trim(strings.ToLower(str_slice[0]), " ") == "name" {
			str_slice[1] = strings.ToLower(str_slice[1])
			f := strings.Join(strings.Split(str_slice[1], " "), "")
			ev.Id = strings.Trim(f, " ") + time.Now().Month().String()
			ev.Name = strings.Trim(str_slice[1], " ")
			count++
		} else if strings.Trim(strings.ToLower(str_slice[0]), " ") == "fields" {
			f := strings.Split(str_slice[1], ",")
			for _, sf := range f {
				ev.Fields = append(ev.Fields, sf)
			}
			count++
		}

	}
	if count != 2 {
		SetState(username, 3)
		return Texts["3--"]
	}
	var err error
	funcs.Insert(ses.Ses, "event", ev, &err)
	if err != nil {
		SetState(username, 3)
		return Texts["3-"]
	}
	AddEventDbState(username, ev.Id)
	return Texts["3++"] + ev.Id
}

func StateSix(username string, message string) string {
	eventid := GetEventDbState(username)
	if eventid == "" {
		SetState(username, 9)
		TimerState[username] = 10
		return Texts["6-"]
	}
	AddDbURL(eventid, message)
	return Texts["6-+"]
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

func StateEight(ses MSession, username string, res string) string {
	return ""
}
