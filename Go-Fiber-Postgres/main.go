package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/YuvrajSingh3110/goFiberPostgres/models"
	"github.com/YuvrajSingh3110/goFiberPostgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

type Repository struct {
	DB *gorm.DB
}

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Request failed",
		})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Could not create the book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Book has been created",
	})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Could not get books",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModels := &models.Books{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("The id is ", id)
	err := r.DB.Where("id = ?", id).First(bookModels).Error
	if err != nil{
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "could not get the book with the given id",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "book id fetched successfully",
			"data": bookModels,
		})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModels := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModels, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Could not delete book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book is deleted",
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookById)
	api.Get("/get_books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
