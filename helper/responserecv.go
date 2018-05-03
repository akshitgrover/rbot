package helper

import (
	"strconv"
)

var TimerState = make(map[string]int)

func RecvResponse(ses MSession, username string, message string) string {
	var state = GetState(username)
	if !Instants[message] && state == 0 {
		return Texts["hello"]
	}
	if !Instants[message] && state == -1 {
		return Texts["default"]
	}
	if Instants[message] && !CheckSession(username) && ReqActiveEidInstants[message] {
		SetState(username, 1)
		return Texts["1+"]
	}
	if Instants[message] && !ReqActiveEidInstants[message] {
		SetState(username, StateInstants[message])
		return Texts[strconv.Itoa(GetState(username))+"+"]
	}
	if state == 1 && !ses.CheckEventValid(message) {
		return Texts["1--"]
	} else if state == 1 {
		AddSession(username, message)
		AddTimer(username)
		SetState(username, -1)
		return Texts["1-+"]
	}
	if state == -1 {
		if !CheckTimer(username) && ReqActiveEidInstants[message] {
			SetState(username, StateInstants[message])
			TimerState[username] = GetState(username)
			SetState(username, 7)
			return Texts["timerexp"]
		} else if CheckTimer(username) && ReqActiveEidInstants[message] {
			SetState(username, StateInstants[message])
			return RecvResponse(ses, username, UserSession[username])
		}
	}
	if state == 2 {
		SetState(username, -1)
		return StateTwo(ses, message)
	}
	if state == 5 {
		SetState(username, -1)
		AddAdmin(message)
		return Texts["5-"]
	}
	if state == 7 {
		return StateSeven(ses, username, message)
	}
	if state == 9 && !ses.CheckEventValid(message) {
		return Texts["1--"]
	} else if state == 9 {
		AddSession(username, message)
		SetState(username, TimerState[username])
		TimerState[username] = 0
		AddTimer(username)
		return RecvResponse(ses, username, UserSession[username])
	}
	return ""
}
