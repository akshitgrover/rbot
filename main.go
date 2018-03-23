package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"rbot/controllers"
	"time"
)

type c_json struct {
	Api_token string `json:"api_token"`
	Mongo_Url string `json:"mongo_url"`
}

var UserSessions map[string]time.Time
var Admin map[string]int
var ActiveEvent map[string]string
var RegisterFlag bool

func main() {

	// Reading JSON

	data := ReadJson("./config/config.json")

	// END

	// --------------------------

	// Start Main Script

	session, err := mgo.Dial(data.Mongo_Url)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Register Mongo Connection

	d := funcs.Db{Session: session}

	// END

	// --------------------------

	// API Call: Register

	go http.HandleFunc("/event", d.InsertApi)

	go http.ListenAndServe(":8080", nil)

	// END

	// ---------------------------

	// Slack API

	api := slack.New(data.Api_token)
	rtm := api.NewRTM()
	info, _, _ := api.StartRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore

		case *slack.ConnectedEvent:
			// Ignore

		case *slack.IMCreatedEvent:
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello :)", ev.Channel.ID))

		case *slack.IMCloseEvent:
			// Ignore

		case *slack.MessageEvent:
			username := info.GetUserByID(ev.Msg.User).Name
			message := ev.Msg.Text
			if Admin[username] == 0 && d.CheckAdmin(username) == 1 {
				Admin[username] = 1
				UserSessions[username] = time.Now()
			} else if Admin[username] == 0 {
				Admin[username] = 2
			}
			if Admin[username] != 1 {
				println("\n\nNot An Admin")
			} else if !ValidateUserSession(username) {
				// Ask For EventID
				RegisterFlag = true
				println("\n\nInactive")
				UserSessions[username] = time.Now()
				continue
			}
			if RegisterFlag {
				// Add Active EventID
				RegisterFlag = false
				continue
			}
			fmt.Println("\n-------")
			fmt.Printf("User: %v\n", username)
			fmt.Printf("Message: %v\n", message)
			fmt.Println("-------\n")
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello Akshit", ev.Channel))

		case *slack.PresenceChangeEvent:
			// Ignore

		case *slack.LatencyReport:
			// Ignore

		case *slack.RTMError:
			fmt.Printf("\n\nError: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("\n\nInvalid credentials")
			return

		default:
			fmt.Printf("\n\nUnable To Catch Event, Event Type: %T", ev)
		}
	}

	// END Main Script

}

func ReadJson(file string) c_json {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fi, _ := f.Stat()
	size := fi.Size()
	bs := make([]byte, size)
	_, _ = f.Read(bs)
	var rf c_json
	err = json.Unmarshal(bs, &rf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return rf
}

func ValidateUserSession(username string) bool {
	d := time.Now().Sub(UserSessions[username])
	if d > time.Minute*2 {
		return false
	}
	return true
}
