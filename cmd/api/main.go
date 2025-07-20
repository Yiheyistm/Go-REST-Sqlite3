package main

import (
	"database/sql"
	"log"

	_ "github.com/Yiheyistm/go-restful-api/docs"
	"github.com/Yiheyistm/go-restful-api/internal/database"
	"github.com/Yiheyistm/go-restful-api/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// @title           Go RESTful API
// @version         1.0
// @description     This is a sample RESTful API server in Go, using Gin and Sqlite3.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

type application struct {
	Port      int
	JwtSecret string
	Model     database.Models
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		Port:      env.GetEnvInt("PORT", 8080),
		JwtSecret: env.GetEnvString("JWT_SECRET", "some_secret_123"),
		Model:     models,
	}

	if err := app.server(); err != nil {
		log.Fatal(err)
	}
}
