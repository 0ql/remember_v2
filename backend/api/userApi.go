package api

import (
	"main/database"
	"time"

	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
)

func CheckIfValidSession(c *fiber.Ctx, allowedRanks []string) (bool, database.User) {
	var user database.User
	res := database.Db.Where("sessiontoken = ?", c.Cookies("session_id")).First(&user)

	if res.Error != nil {
		return false, user
	}

	if user.Expires < time.Now().Unix() {
		return false, user
	}

	for i := range allowedRanks {
		if allowedRanks[i] == user.Rank {
			return true, user
		}
	}

	return false, user
}

func CheckIfValidSocket(kws *ikisocket.Websocket, allowedRanks []string) (bool, database.User) {
	var user database.User
	res := database.Db.Where("sessiontoken = ?", kws.Cookies("session_id")).First(&user)

	if res.Error != nil {
		return false, user
	}

	if user.Expires < time.Now().Unix() {
		return false, user
	}

	for i := range allowedRanks {
		if allowedRanks[i] == user.Rank {
			return true, user
		}
	}

	return false, user
}

func GetMyInformation(c *fiber.Ctx) error {
	c.Accepts("json")

	isValid, user := CheckIfValidSession(c, []string{"member", "admin"})

	if !isValid {
		return c.SendStatus(400)
	}

	return c.JSON(fiber.Map{
		"username": user.Username,
		"email":    user.Email,
		"rank":     user.Rank,
		"elo":      user.Elo,
	})
}
