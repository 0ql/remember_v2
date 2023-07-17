import { ParentProps, createContext } from "solid-js"
import { createStore, SetStoreFunction } from "solid-js/store"
import { Navigator } from "@solidjs/router"

export enum EventsCtoS {
	JoinWaiting,
	LeaveWaiting,
	Answers
}

export enum EventsStoC {
	AllClients,
	UserData,
	NewGame,
	Results
}

export enum Status {
	Idle,
	Waiting,
	InGame
}

export type MessageCtoS = {
	Event: EventsCtoS
	Data: any
}

export type MessageStoC = {
	Event: EventsStoC
	Data: any
}

export type Client = {
	Status: Status;
	Username: string;
	Rank: string;
	Elo: number;
}

export type Game = {
	PlayerAUsername: string
	PlayerBUsername: string
	Questions: {
		Question: string
		Options: string[]
	}[]
	TimeUp: number
}

type GlobalContextType = {
	clients: Client[],
	userdata: Client | null,
	game: Game | null
}

const [state, setState] = createStore<GlobalContextType>({
	clients: [],
	userdata: null,
	game: null
})

export const GlobalContext = createContext<[get: GlobalContextType, set: SetStoreFunction<GlobalContextType>]>([state, setState])

let socket: WebSocket

export const joinWaitingRoom = () => {
	// state.Status = Status.Waiting
	socket.send(JSON.stringify({
		Event: EventsCtoS.JoinWaiting
	}))
}

export const leaveWaitingRoom = () => {
	socket.send(JSON.stringify({
		Event: EventsCtoS.LeaveWaiting
	}))
}

export const endGame = (answers: number[]) => {
	console.log(answers)
	socket.send(JSON.stringify({
		Event: EventsCtoS.Answers,
		Data: answers
	}))
}

export const createWebSocket = (navigate: Navigator) => {
	socket = new WebSocket("ws://localhost:5000/ws/liveWaitingList")

	socket.onopen = (e) => {
		console.log("Websocket open")
	};

	socket.onmessage = (e: MessageEvent) => {
		const msg: MessageStoC = JSON.parse(e.data)
		switch (msg.Event) {
			case EventsStoC.UserData:
				console.log("UserData Receved: ", msg.Data)
				createContext()
				setState(state => {
					return {
						clients: state.clients,
						userdata: msg.Data
					}
				})
				break
			case EventsStoC.AllClients:
				console.log("AllClients Receved:", msg.Data)
				setState(state => {
					return {
						clients: msg.Data,
						userdata: state.userdata
					}
				})
				break
			case EventsStoC.NewGame:
				console.log("NewGame Receved:", msg.Data)
				navigate("/game")
				setState(state => {
					return {
						state,
						game: msg.Data
					}
				})
				break
			case EventsStoC.Results:
				console.log("Results Receved:", msg.Data)
				break
		}
	};

	socket.onclose = (e) => {
		if (e.wasClean) console.log("Ferbindung wurde geschlossen")
		else console.log("Verbindung wurde unterbrochen")
	};

	socket.onerror = (e) => {
		console.error(e)
	};
}

export const GlobalContextProvider = (props: ParentProps) => {
	return (
		<GlobalContext.Provider value={[state, setState]}>
			{props.children}
		</GlobalContext.Provider>
	);
}
