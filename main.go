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
	"rbot/models"
)

type c_json struct {
	Api_token string `json:"api_token"`
}

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
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

	// --------------------------

	// Reading JSON

	data := ReadJson("./config/config.json")

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
			fmt.Println("-------")
			fmt.Printf("User: %v\n", info.GetUserByID(ev.Msg.User).Name)
			fmt.Printf("Message: %v\n", ev.Msg.Text)
			fmt.Println("-------")
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello Akshit", ev.Channel))

		case *slack.PresenceChangeEvent:
			// Ignore

		case *slack.LatencyReport:
			// Ignore

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore
		}
	}

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
