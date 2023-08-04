package wsocket

import (
	"encoding/json"
	"fmt"
	"main/database"
	"math"
	"math/rand"
	"time"

	"github.com/antoniodipinto/ikisocket"
)

type Status int

const (
	Idle    Status = 0
	Waiting Status = 1
	InGame  Status = 2
)

type EventsCtoS int

const (
	JoinWaiting  EventsCtoS = 0
	LeaveWaiting EventsCtoS = 1
	Answers      EventsCtoS = 2
)

type EventsStoC int

const (
	AllClients EventsStoC = 0
	UserData   EventsStoC = 1
	NewGame    EventsStoC = 2
	Results    EventsStoC = 3
)

type Client struct {
	SocketID string `json:"-"`
	Status   Status
	Username string
	Rank     string
	Elo      float64
}

type QuestionOptions struct {
	Question string
	Options  []string
}

type Game struct {
	PlayerASessionID string `json:"-"`
	PlayerBSessionID string `json:"-"`
	PlayerAUsername  string
	PlayerBUsername  string
	PlayerACorrect   int `json:"-"`
	PlayerBCorrect   int `json:"-"`
	Questions        [5]QuestionOptions
	TimeUp           int64
	Answers          [5]int `json:"-"`
}

type ResultsData struct {
	Answers    [5]int
	PlayerAElo float64
	PlayerBElo float64
}

type MessageCtoS struct {
	Event EventsCtoS
	Data  interface{}
}

type MessageStoC struct {
	Event EventsStoC
	Data  interface{}
}

var Clients = make(map[string]*Client)
var Games = []*Game{}

func BroadcastClients(kws *ikisocket.Websocket, except bool) {
	clientsSlice := []*Client{}
	for _, v := range Clients {
		clientsSlice = append(clientsSlice, v)
	}

	jn, _ := json.Marshal(MessageStoC{
		Event: AllClients,
		Data:  clientsSlice,
	})
	kws.Broadcast(jn, except)
}

func CloseConnection(kws *ikisocket.Websocket, userid string) {
	delete(Clients, userid)
	BroadcastClients(kws, true)
	kws.Close()
}

func Matching(kws *ikisocket.Websocket) {
	// instant matching
	// TODO: implement queue
	var (
		playerASessionID string = kws.GetStringAttribute("session_id")
		playerBSessionID string
	)
	for k, client := range Clients {
		if client.Status == Waiting && k != playerASessionID {
			playerBSessionID = k
			break
		}
	}
	if playerBSessionID == "" {
		return
	}

	type CleanQuestions []struct {
		Question string
		Answer   string
		Wrong0   string
		Wrong1   string
		Wrong2   string
	}

	var result CleanQuestions
	database.Db.Debug().Raw("SELECT question, answer, wrong0, wrong1, wrong2 FROM questions ORDER BY RANDOM() LIMIT 5").Scan(&result)

	options := [5]QuestionOptions{}
	answers := [5]int{}
	for i := range result {
		options[i].Question = result[i].Question
		options[i].Options = append(options[i].Options, result[i].Answer)
		options[i].Options = append(options[i].Options, result[i].Wrong0)
		options[i].Options = append(options[i].Options, result[i].Wrong1)
		options[i].Options = append(options[i].Options, result[i].Wrong2)
		correct := 0
		rand.Shuffle(4, func(n, j int) {
			if n == correct {
				correct = j
			} else if j == correct {
				correct = n
			}
			options[i].Options[n], options[i].Options[j] = options[i].Options[j], options[i].Options[n]
		})
		answers[i] = correct
	}

	rand.Seed(time.Now().UnixNano())
	game := &Game{
		PlayerASessionID: playerASessionID,
		PlayerBSessionID: playerBSessionID,
		PlayerAUsername:  Clients[playerASessionID].Username,
		PlayerBUsername:  Clients[playerBSessionID].Username,
		PlayerACorrect:   -1,
		PlayerBCorrect:   -1,
		Questions:        options,
		TimeUp:           time.Now().Add(time.Minute).Unix(),
		Answers:          answers,
	}
	Games = append(Games, game)
	jn, _ := json.Marshal(MessageStoC{
		Event: NewGame,
		Data:  game,
	})
	Clients[playerASessionID].Status = InGame
	Clients[playerBSessionID].Status = InGame
	kws.Emit(jn)
	ikisocket.EmitTo(Clients[playerBSessionID].SocketID, jn)
	BroadcastClients(kws, false)
}

func EloFormula(a float64, b float64, resa float64, resb float64) (float64, float64) {
	return a + resa - (1 / (1 + math.Pow(10, (b - a) / 400))), b + resb - (1 / (1 + math.Pow(10, (a - b) / 400)))
}

func EndGameSendResults(kws *ikisocket.Websocket, game *Game) {
	fmt.Println("sending results")
	var s, sb = 0.5, 0.5
	if game.PlayerACorrect > game.PlayerBCorrect {
		s = 1
		sb = 0
	} else if game.PlayerACorrect < game.PlayerBCorrect {
		s = 0
		sb = 1
	}
	Clients[game.PlayerASessionID].Elo, Clients[game.PlayerBSessionID].Elo = EloFormula(Clients[game.PlayerASessionID].Elo, Clients[game.PlayerBSessionID].Elo, s, sb)
	jn, _ := json.Marshal(MessageStoC{
		Event: Results,
		Data: ResultsData{
			Answers:    game.Answers,
			PlayerAElo: Clients[game.PlayerASessionID].Elo,
			PlayerBElo: Clients[game.PlayerBSessionID].Elo,
		},
	})
	ikisocket.EmitTo(Clients[game.PlayerASessionID].SocketID, jn)
	ikisocket.EmitTo(Clients[game.PlayerBSessionID].SocketID, jn)
	game.PlayerACorrect = -1
	game.PlayerBCorrect = -1
	Clients[game.PlayerASessionID].Status = Idle
	Clients[game.PlayerBSessionID].Status = Idle
	BroadcastClients(kws, false)
}

func EventHandlers() {
	ikisocket.On(ikisocket.EventDisconnect, func(e *ikisocket.EventPayload) {
		delete(Clients, e.Kws.GetStringAttribute("session_id"))
		BroadcastClients(e.Kws, true)
	})
	ikisocket.On(ikisocket.EventMessage, func(e *ikisocket.EventPayload) {
		message := MessageCtoS{}
		err := json.Unmarshal(e.Data, &message)
		if err != nil {
			e.Kws.Close()
			return
		}

		switch message.Event {
		case JoinWaiting:
			Clients[e.Kws.GetStringAttribute("session_id")].Status = Waiting
			BroadcastClients(e.Kws, false)
			Matching(e.Kws)
		case LeaveWaiting:
			Clients[e.Kws.GetStringAttribute("session_id")].Status = Idle
			BroadcastClients(e.Kws, false)
		case Answers:
			for _, game := range Games {
				isPlayerA := false
				fmt.Println("a")
				if game.PlayerASessionID == e.Kws.GetStringAttribute("session_id") {
					isPlayerA = true
				} else if game.PlayerBSessionID != e.Kws.GetStringAttribute("session_id") {
					continue
				}
				fmt.Println("b")
				correct := 0
				for i, answer := range message.Data.([]interface{}) {
					if answer != float64(game.Answers[i]) {
						correct += 1
					}
				}
				fmt.Println(correct)
				fmt.Println(game.PlayerACorrect)
				fmt.Println(game.PlayerBCorrect)
				if isPlayerA {
					fmt.Println("setting PlayerACorrect")
					game.PlayerACorrect = correct
					if game.PlayerBCorrect != -1 {
						fmt.Println("sending resultsA...")
						EndGameSendResults(e.Kws, game)
					}
				} else {
					fmt.Println("setting PlayerBCorrect")
					game.PlayerBCorrect = correct
					if game.PlayerACorrect != -1 {
						fmt.Println("sending resultsB...")
						EndGameSendResults(e.Kws, game)
					}
				}
				break
			}
		}
	})
}

func CreateClient(kws *ikisocket.Websocket) {
	session_id := kws.Cookies("session_id")
	kws.SetAttribute("session_id", session_id)
	Clients[session_id].SocketID = kws.GetUUID()

	jn, _ := json.Marshal(MessageStoC{
		Event: UserData,
		Data:  Clients[session_id],
	})
	kws.Emit(jn)

	BroadcastClients(kws, false)
}
