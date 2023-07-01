package main

import (
	"fmt"

	database "github.com/YuvrajSingh3110/go-fibre_simple_CRM/Database"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
)

func setupRoutes(app *fiber.App) {
	app.Get(GetLeads)
	app.Get(GetOneLead)
	app.Post(NewLead)
	app.Delete(DeleteLead)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "leads.db")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database...")
	database.DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("Database migrated...")
}

func main() {
	app := fiber.New()
	initDatabase()
	defer database.DBConn.Close()
	setupRoutes(app)
	app.Listen(3000)
}
