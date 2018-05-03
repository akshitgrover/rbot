package helper

import (
//"strconv"
)

var TimerState = make(map[string]int)

func RecvResponse(ses MSession, username string, message string) string {
	var state = GetState(username)
	println(state)
	if !Instants[message] && state == 0 {
		return Texts["hello"]
	}
	if !Instants[message] && state <= 0 {
		return Texts["default"]
	}
	if Instants[message] && !CheckSession(username) && state != 1 {
		SetState(username, 1)
		return Texts["1+"]
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
		SetState(username, StateInstants[message])
		if !CheckTimer(username) {
			TimerState[username] = GetState(username)
			SetState(username, 7)
			return Texts["timerexp"]
		}
		SetState(username, -1)
		return StateTwo(ses, UserSession[username])
	}
	if state == 2 {
		SetState(username, -1)
		return StateTwo(ses, message)
	}
	if state == 7 {
		return StateSeven(ses, username, message)
	}
	if state == 9 && !ses.CheckEventValid(message) {
		return Texts["1--"]
	} else if state == 9 {
		AddSession(username, message)
		flag_state := TimerState[username]
		SetState(username, -1)
		TimerState[username] = 0
		AddTimer(username)
		return RecvResponse(ses, username, InverStateInstants[flag_state])
	}
	return ""
}
