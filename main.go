package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	PublisherID uint
	Publisher   Publisher
	Authors     []Author `gorm:"many2many:author_books;"`
}

type Publisher struct {
	gorm.Model
	Details string
	Name    string
}

type Author struct {
	gorm.Model
	Name  string
	Books []Book `gorm:"many2many:author_books;"`
}

type AuthorBook struct {
	AuthorID uint
	Author   Author
	BookID   uint
	Book     Book
}

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
	dbname   = "mydatabase"
)

func main() {
	// create connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
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

	// connect to postgres database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	// kill process when connection lost
	if err != nil {
		panic("failed to connect database !")
	}

	fmt.Println("Connect to Datasbase Successful !")

	//Start Auto Migration
	db.AutoMigrate(&Book{}, &Publisher{}, &Author{}, &AuthorBook{})

	// // ขาสร้าง
	// publisher := Publisher{
	// 	Details: "JaoSalmon",
	// 	Name:    "JaoJao",
	// }
	// // _ = createPublisher(db, &publisher)

	// // Example data for a new author
	// author1 := Author{
	// 	Name: "JanJao",
	// }

	// _ = createAuthor(db, &author1)
	// author2 := Author{
	// 	Name: "JanJao",
	// }

	// // _ = createAuthor(db, &author2)

	// // Example data for a new book with an author
	// book := Book{
	// 	Name:        "Book JAO1",
	// 	Author:      "test",
	// 	Description: "Book Description",
	// 	PublisherID: publisher.ID,               // Use the ID of the publisher created above
	// 	Authors:     []Author{author1, author2}, // Add the created author
	// }
	// // _ = createBookWithAuthor(db, &book)

	// ขาเรียก

	// Example: Get a book with its publisher
	bookWithPublisher, err := getBookWithPublisher(db, 2) // assuming a book with ID 1
	if err != nil {
		// Handle error
	}

	// Example: Get a book with its authors
	bookWithAuthors, err := getBookWithAuthors(db, 2) // assuming a book with ID 1
	if err != nil {
		// Handle error
	}

	// Example: List books of a specific author
	authorBooks, err := listBooksOfAuthor(db, 2) // assuming an author with ID 1
	if err != nil {
		// Handle error
	}

	fmt.Println(bookWithPublisher)
	fmt.Println(bookWithAuthors)
	fmt.Println(authorBooks)

}

func createPublisher(db *gorm.DB, publisher *Publisher) error {
	result := db.Create(publisher)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func createAuthor(db *gorm.DB, author *Author) error {
	result := db.Create(author)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func createBookWithAuthor(db *gorm.DB, book *Book) error {
	// First, create the book
	if err := db.Create(book).Error; err != nil {
		return err
	}

	return nil
}
func getBookWithPublisher(db *gorm.DB, bookID uint) (*Book, error) {
	var book Book
	result := db.Preload("Publisher").First(&book, bookID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func getBookWithAuthors(db *gorm.DB, bookID uint) (*Book, error) {
	var book Book
	result := db.Preload("Authors").First(&book, bookID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func listBooksOfAuthor(db *gorm.DB, authorID uint) ([]Book, error) {
	var books []Book
	result := db.Joins("JOIN author_books on author_books.book_id = books.id").
		Where("author_books.author_id = ?", authorID).
		Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}
