package main

import (
	"log"
	"main/api"
	"main/wsocket"

	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()
	ws := app.Group("/ws")
	Api := app.Group("/api")
	v1 := Api.Group("/v1")

	ws.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	wsocket.EventHandlers()
	ws.Get("/liveWaitingList", ikisocket.New(wsocket.WaitingUsers))

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
