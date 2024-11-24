package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func createBook(db *gorm.DB, book *Book) {
	result := db.Create(book)

	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
	}

	fmt.Println("Create Book successful!")
}

func getBook(db *gorm.DB, id uint) *Book {
	var book Book
	result := db.First(&book, id)

	if result.Error != nil {
		log.Fatalf("Error get book: %v", result.Error)
	}

	return &book
}

func getBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)

	if result.Error != nil {
		log.Fatalf("Error get book: %v", result.Error)
	}

	return books
}

func updateBook(db *gorm.DB, book *Book) {
	result := db.Save(&book)

	if result.Error != nil {
		log.Fatalf("Error update book:%v", result.Error)
	}

	fmt.Println("Updated Book successful!")
}

func deleteBook(db *gorm.DB, id uint) {
	var book Book

	result := db.Delete(&book, id)
	if result.Error != nil {
		log.Fatalf("Error Delete book:%v", result.Error)
	}

	fmt.Println("Delete Book successful!")
}

func searchBook(db *gorm.DB, bookName string) *Book {
	var book Book

	result := db.Where("name = ?", bookName).First(&book)

	if result.Error != nil {
		log.Fatalf("Error Not Found book:%v", result.Error)
	}

	return &book
}

func searchBooksbyAuthor(db *gorm.DB, authorName string) []Book {
	var books []Book

	result := db.Where("author = ?", authorName).Order("price desc").Find(&books)

	if result.Error != nil {
		log.Fatalf("Error Not Found book:%v", result.Error)
	}

	return books
}
