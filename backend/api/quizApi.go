package api

import (
	"fmt"

	"main/database"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type QuestionValidation struct {
	Question string `validate:"required,min=3,max=32"`
	Answer   string `validate:"required,min=3,max=32"`
	Wrong0   string `validate:"required,min=3,max=32"`
	Wrong1   string `validate:"required,min=3,max=32"`
	Wrong2   string `validate:"required,min=3,max=32"`
}

func validateStruct(q interface{}) bool {
	validate := validator.New()

	err := validate.Struct(q)

	return err == nil
}

func CreateQuestion(c *fiber.Ctx) error {
	isValid, _ := CheckIfValidSession(c, []string{"admin"})

	if !isValid {
		return c.SendStatus(401)
	}

	c.Accepts("json")

	q := new(QuestionValidation)
	if err := c.BodyParser(q); err != nil {
		return c.SendStatus(500)
	}

	if !validateStruct(q) {
		return c.SendStatus(400)
	}

	var question database.Question
	question.Question = q.Question
	question.Answer = q.Answer
	question.Wrong0 = q.Wrong0
	question.Wrong1 = q.Wrong1
	question.Wrong2 = q.Wrong2

	database.Db.Debug().Create(&question)

	return c.SendStatus(200)
}

func GetAllQuestions(c *fiber.Ctx) error {
	isValid, _ := CheckIfValidSession(c, []string{"admin"})

	if !isValid {
		return c.SendStatus(401)
	}

	var questions []database.Question
	database.Db.Debug().Find(&questions)

	return c.JSON(questions)
}

func GetQuestionById(c *fiber.Ctx) error {
	var id string = c.Params("id")
	var questions []database.Question
	database.Db.Debug().Find(&questions, &id)

	return c.JSON(questions)
}

func GetRandomQuiz(c *fiber.Ctx) error {
	isValid, _ := CheckIfValidSession(c, []string{"member"})
	if !isValid {
		c.SendStatus(401)
	}
	var result database.Question
	database.Db.Debug().Raw("SELECT * FROM questions ORDER BY RANDOM() LIMIT 1").Scan(&result)

	return c.JSON(result)
}

func DeleteQuestionById(c *fiber.Ctx) error {
	isValid, _ := CheckIfValidSession(c, []string{"admin"})
	if !isValid {
		return c.SendStatus(401)
	}

	var id string = c.Params("id")
	var questions []database.Question
	database.Db.Debug().Delete(&questions, &id)

	return c.SendStatus(200)
}

type WaitingListValidation struct {
	name string `validate:"required,min=3,max=32"`
}

type Folder struct {
	Name         string `json:"name"`
	FoldersIDs   []int  `json:"folders"`
	QuestionsIDs []int  `json:"questions"`
}

func EnterWaitingList(c *fiber.Ctx) error {
	isValid, _ := CheckIfValidSession(c, []string{"member", "admin"})
	if !isValid {
		return c.SendStatus(401)
	}

	c.Accepts("json", "text")

	w := new(WaitingListValidation)
	if err := c.BodyParser(w); err != nil {
		fmt.Println(err)
		return c.SendStatus(500)
	}

	if !validateStruct(w) {
		return c.SendStatus(400)
	}

	// var question database.Question

	// database.Db.Debug().Create(&Folder)

	fmt.Println(w.name)

	return c.SendStatus(200)
}
