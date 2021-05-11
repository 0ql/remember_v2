package wsocket

import (
	"encoding/json"
	"fmt"
	"main/api"

	"github.com/antoniodipinto/ikisocket"
)

type WaitingUserInfo struct {
	Username string
	Rank     string
	Elo      int
}

var clients = make(map[string]WaitingUserInfo)

func EventHandlers() {
	ikisocket.On(ikisocket.EventDisconnect, func(e *ikisocket.EventPayload) {
		delete(clients, e.Kws.GetStringAttribute("user_id"))

		// Broadcast clients when new User enters
		jn, _ := json.Marshal(clients)
		e.Kws.Broadcast(jn, true)
	})
}

func WaitingUsers(kws *ikisocket.Websocket) {
	isValid, user := api.CheckIfValidSocket(kws, []string{"member", "admin"})
	if !isValid {
		kws.Close()
	}

	kws.SetAttribute("user_id", fmt.Sprint(user.ID))

	var wUser WaitingUserInfo
	wUser.Username = user.Username
	wUser.Rank = user.Rank
	wUser.Elo = user.Elo

	clients[fmt.Sprint(user.ID)] = wUser

	// Broadcast clients when new User enters
	jn, _ := json.Marshal(clients)
	kws.Broadcast(jn, true)
}
