package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"rbot/config"
	"rbot/helper"
	"strings"
)

type c_json struct {
	Api_token string `json:"api_token"`
	Mongo_Url string `json:"mongo_url"`
}

// var UserSessions = make(map[string]time.Time)
// var Admin = make(map[string]int)

func main() {

	// Start Main Script

	data := helper.ReadConfigJson()
	msession, err := mgo.Dial(config.Mongo_Url)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fses := helper.MSession{Ses: msession}

	// --------------------------

	// API Call: Register

	go http.ListenAndServe(":8080", nil)

	// END

	// ---------------------------

	// Slack API

	api := slack.New(data.Token)
	rtm := api.NewRTM()
	info, _, _ := api.StartRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			//Ignore

		case *slack.ConnectedEvent:
			// Ignore

		case *slack.IMCreatedEvent:
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello :)", ev.Channel.ID))

		case *slack.IMCloseEvent:
			// Ignore

		case *slack.MessageEvent:
			username := info.GetUserByID(ev.Msg.User).Name
			message := strings.ToLower(ev.Msg.Text)

			go LogMessage(username, message)
			msg := helper.RecvResponse(fses, username, message)
			log.Println(msg + "~")
			rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))

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

func LogMessage(username string, message string) {
	fmt.Println("\n-------")
	fmt.Printf("User: %v\n", username)
	fmt.Printf("Message: %v\n", message)
	fmt.Println("-------\n")
}
