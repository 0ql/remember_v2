package main

import (
	"crypto/sha256"
	"encoding/base64"
	"main/database"
	"time"

	"github.com/dchest/uniuri"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// Validates the given struct
func validateStruct(q interface{}) bool {
	validate := validator.New()
	err := validate.Struct(q)
	return err == nil
}

type credentialValidation struct {
	Email    string `validate:"required,min=3,max=32"`
	Password string `validate:"required,min=3,max=32"`
}

// Login with Email and Password
// Sends back SessionCookie
func LoginWithCredentials(c *fiber.Ctx) error {
	c.Accepts("json")

	credentials := new(credentialValidation)
	if err := c.BodyParser(credentials); err != nil {
		return c.SendStatus(500)
	}

	if !validateStruct(credentials) {
		return c.SendStatus(400)
	}

	// hash the Password
	h := sha256.New()
	h.Write([]byte(credentials.Password))

	// look up User
	var user database.User
	res := database.Db.Debug().Where("email = ? AND pwdhash = ?", credentials.Email, base64.URLEncoding.EncodeToString(h.Sum(nil))).First(&user)

	if res.Error != nil {
		return c.SendStatus(400)
	}

	// Update DB with new Session Token
	user.Sessiontoken = uniuri.NewLen(24)
	expires := time.Now().Add(10 * 24 * time.Hour)
	user.Expires = expires.Unix()
	database.Db.Save(&user)

	// generate Cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = user.Sessiontoken
	cookie.Expires = expires

	c.Cookie(cookie)

	return c.SendStatus(200)
}

// Sends a JSON array off all Users in the Database
func GetAllUsers(c *fiber.Ctx) error {
	var users []database.User

	database.Db.Find(&users)

	return c.JSON(users)
}

type userValidation struct {
	Email    string `validate:"required,min=3,max=32"`
	UserName string `validate:"required,min=3,max=32"`
	Password string `validate:"required,min=3,max=32"`
}

// Creates a new User based on JSON input
func CreateUser(c *fiber.Ctx) error {
	c.Accepts("json")

	// parse JSON input
	inputUser := new(userValidation)
	if err := c.BodyParser(inputUser); err != nil {
		return c.SendStatus(500)
	}

	// Validate JSON request
	if !validateStruct(inputUser) {
		return c.SendStatus(400)
	}

	// Check if username or email exists
	res := database.Db.Debug().Where("username = ? OR email = ?", inputUser.UserName, inputUser.Email).First(database.User{})
	if res.Error == nil {
		// username exists
		return c.SendString("Username Taken")
	}

	var newUser database.User
	newUser.Email = inputUser.Email
	newUser.Username = inputUser.UserName
	newUser.Rank = "admin"
	newUser.Elo = 100

	// Create Hash of Password
	h := sha256.New()
	h.Write([]byte(inputUser.Password))
	newUser.Pwdhash = base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Generate Session Token
	newUser.Sessiontoken = uniuri.NewLen(24)
	expires := time.Now().Add(10 * 24 * time.Hour)
	newUser.Expires = expires.Unix()
	database.Db.Create(&newUser)

	// generate Cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = newUser.Sessiontoken
	cookie.Expires = expires

	c.Cookie(cookie)

	return c.SendStatus(200)
}
