package main

import (
	"fmt"

	"github.com/YuvrajSingh3110/go-fibre_simple_CRM/database"
	"github.com/YuvrajSingh3110/go-fibre_simple_CRM/lead"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func setupRoutes(app *fiber.App) {
	app.Get("api/v1/lead", lead.GetLeads)
	app.Get("api/v1/lead/:id", lead.GetOneLead)
	app.Post("api/v1/lead", lead.NewLead)
	app.Delete("api/v1/lead/:id", lead.DeleteLead)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "leads.db")
	if err != nil {
		panic("Failed to connect database")
	}
	fmt.Println("Connected to database...")
	database.DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("Database migrated...")
}

func main() {
	fmt.Println("CRM basic api")
	app := fiber.New()
	initDatabase()
	defer database.DBConn.Close()
	setupRoutes(app)
	app.Listen(3000)
}
