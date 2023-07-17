package main

import (
	"fmt"
	"log"
	"main/api"
	"main/wsocket"
	"time"

	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

func printClients() {
	for {
		fmt.Println(wsocket.Clients)
		time.Sleep(time.Second)
	}
}

func main() {
	go printClients()

	app := fiber.New()
	ws := app.Group("/ws")
	Api := app.Group("/api")
	v1 := Api.Group("/v1")

	ws.Use("/", func(c *fiber.Ctx) error {
		isValid, user := api.CheckIfValidSession(c, []string{"admin"})
		if !isValid {
			return fiber.ErrBadRequest
		}

		session_id := c.Cookies("session_id")
		if _, ok := wsocket.Clients[session_id]; ok {
			return fiber.ErrBadRequest
		}

		wsocket.Clients[session_id] = &wsocket.Client{
			Status:   wsocket.Idle,
			Username: user.Username,
			Rank:     user.Rank,
			Elo:      user.Elo,
		}

		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	wsocket.EventHandlers()
	ws.Get("/liveWaitingList", ikisocket.New(wsocket.CreateClient))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	v1.Post("/login", LoginWithCredentials)
	v1.Get("/getAllUsers", GetAllUsers)
	v1.Post("/createUser", CreateUser)

	v1.Get("/getMyInformation", api.GetMyInformation)

	v1.Post("/createQuestion", api.CreateQuestion)
	v1.Get("/getAllQuestions", api.GetAllQuestions)
	v1.Get("/getRandomQuestion", api.GetRandomQuiz)
	v1.Get("/getQuestionById/:id", api.GetQuestionById)
	v1.Delete("/deleteQuestionById/:id", api.DeleteQuestionById)

	log.Fatal(app.Listen(":5000"))
}
