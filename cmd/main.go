package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/mineracail/zokus/database"
	"github.com/mineracail/zokus/graph"
	"github.com/mineracail/zokus/graph/model"
	"github.com/mineracail/zokus/helper"
	"gorm.io/gorm"
)

const defaultPort = "3001"

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// PerformMigrations performs database migrations
func PerformMigrations(db *gorm.DB) error {
	// Read the init.sql file
	sqlFile, err := helper.ReadFile("database/init.sql")
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	// Execute the SQL commands in the file
	if err := db.Exec(string(sqlFile)).Error; err != nil {
		return fmt.Errorf("failed to execute SQL commands from init.sql: %w", err)
	}

	// Run GORM AutoMigrate for your models
	if err := db.AutoMigrate(
		&Product{}, 
		&model.User{},
		&model.Vender{},
		&model.NewProduct{},
		&model.UserOrderProducts{},
		// &model.ProductRequest{},
		// &model.{},
		); err != nil {
		return fmt.Errorf("failed to run AutoMigrate: %w", err)
	}
	return nil
}
func main() {
	port := os.Getenv("Port")
	if port == "" {
		port = defaultPort
	}
	//// Read environment variables
	dbUser := os.Getenv("UserName")
	dbPassword := os.Getenv("Password")
	dbName := os.Getenv("DBName")
	dbHost := os.Getenv("Host") // Use service name in Docker Compose
	dbPort := "5432"            // Default PostgreSQL port

	// Configure database connection using database.Config
	config := &database.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  "disable", // Adjust SSLMode as per your PostgreSQL setup
	}

	// Establish database connection
	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}
	// Perform database migrations
	err = PerformMigrations(db)
	if err != nil {
		log.Fatal("failed to perform migrations:", err)
	}

	resolver := &graph.Resolver{
		DB: db,
	}

	// Initialize the GraphQL server
	router := chi.NewRouter()
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	router.Handle("/query", server)

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
