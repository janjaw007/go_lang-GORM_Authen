package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

func main() {
	// Create Connection String
	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold time that start query
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	if err != nil {
		// panic = kill process
		panic("failed to connect database")
	}

	fmt.Println(db)
	fmt.Println("Connect Successful!")

	// AutoMigrate will compare struct and table on sql
	// AutoMigrate Recive 2 argument (struct, ) that strcut is structure for tables will be create
	db.AutoMigrate(&Book{}, &User{})
	fmt.Println("Migrate successful!")

	app := fiber.New()

	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(getBooks(db))
	})

	app.Get("/book/:id", func(c *fiber.Ctx) error {
		var id = c.Params("id")
		strId, err := strconv.Atoi(id)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(getBook(db, uint(strId)))
	})

	app.Post("/book/createBook", func(c *fiber.Ctx) error {
		var book Book

		if err := c.BodyParser(&book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := createBook(db, &book); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create book",
			})
		}
		// Return the created book as a response
		return c.JSON(fiber.Map{
			"message": "Create Book Successful",
		})
	})

	app.Put("/book/updateBook/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		strId, err := strconv.Atoi(id)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		book := new(Book)

		book.ID = uint(strId)

		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := updateBook(db, book); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update book",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Update Book Successful",
		})
	})

	app.Delete("/book/deleteBook/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		strId, err := strconv.Atoi(id)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := deleteBook(db, uint(strId)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete book",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Delete Book Successful",
		})

	})

	// User API

	app.Post("/user/register", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = createUser(db, user)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Register Success",
		})
	})

	app.Listen(":8080")

	// CRUD
	// Create

	// newBook := &Book{
	// 	Name:        "JanJaoExcited",
	// 	Author:      "JanJao",
	// 	Description: "JanJao WebDev Explore",
	// 	Price:       700,
	// }

	// createBook(db, newBook)

	// Get
	// currentBook := getBook(db, 1)
	// // fmt.Println("current Book is ", currentBook)

	//  Update
	// currentBook.Name = "New Jao"
	// currentBook.Price = 1000

	// updateBook(db, currentBook)

	// Delete

	// deleteBook(db, 1)

	// Seacrh Function
	// currentBook := searchBook(db, "JanJaoExcited")

	// fmt.Println(currentBook)

	// currentBook := searchBooksbyAuthor(db, "JanJao")

	// for _, book := range currentBook {
	// 	fmt.Println(book.ID, book.Name, book.Author, book.Price)
	// }
}
