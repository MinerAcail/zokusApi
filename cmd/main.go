package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/mineracail/zokus/database"
	"github.com/mineracail/zokus/graph/generated"
	"github.com/mineracail/zokus/graph/model"
	graph "github.com/mineracail/zokus/graph/resolver"
	"github.com/mineracail/zokus/payment"
	"github.com/rpip/paystack-go"

	"github.com/mineracail/zokus/helper"
	"github.com/mineracail/zokus/middleware"
	"gorm.io/gorm"
)

const defaultPort = "3001"

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
		&model.User{},
		&model.Vendor{},
		&model.FavoriteVendor{},
		&model.FavoriteProduct{},
		&model.SideSwapper{},
		&model.SideSwapperProduct{},
		&model.ProductInventoryMetaData{},
		&model.Pricing{},
		&model.Product{},
		&model.Comment{},
		&model.UserOrder{},
		&model.ProductRequest{},
		&model.Cart{},
		&model.CartItem{},
		&model.Order{},	
		&model.DeliveryMetaData{},
		&model.DeliveryStatus{},
		&model.PurchaseMetaData{},
		&model.CartItemGroup{},
		&model.Group{},
	// &model.{},
	); err != nil {
		return fmt.Errorf("failed to run AutoMigrate: %w", err)
	}
	return nil
}

// InitPaystack initializes a new Paystack client with the given API key
func InitPaystack(apiKey string) *paystack.Client {
	customHttpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	return paystack.NewClient(apiKey, customHttpClient)
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

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	// minioBucket := "your-bucket-name" // Replace with your 	 name

	// Read environment variables for Paystack
	paystackAPIKey := os.Getenv("PAYSTACK_API_KEY")
	// Initialize the Paystack API client
	// Initialize the Paystack API client
	// Initialize Paystack client
	paystackClient := InitPaystack(paystackAPIKey)

	// Configure database connection using database.Config
	config := &database.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  "disable", // Adjust SSLMode as per your PostgreSQL setup
	}
	// Initialize the MoMo API client
	momoAPI := payment.NewMoMoAPI(
		os.Getenv("MOMO_API_KEY"),          // Replace with your MoMo API key
		os.Getenv("MOMO_SUBSCRIPTION_KEY"), // Replace with your MoMo subscription key
		os.Getenv("MOMO_ENVIRONMENT"),      // Either "sandbox" or "production"
	)
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

	// Initialize MinIO client
	minioClient, err := helper.InitMinIO(minioEndpoint, minioAccessKey, minioSecretKey)
	if err != nil {
		log.Fatalf("failed to initialize MinIO client: %v", err)
	}

	resolver := &graph.Resolver{
		DB:          db,
		MinIOClient: minioClient,
		MoMoAPI:     momoAPI,
		PaystackAPI: paystackClient,
	}

	// Initialize the GraphQL server
	router := chi.NewRouter()
	// Configure CORS
	corsConfig := cors.New(cors.Options{
		// { "http://localhost:5173", "http://yourdomain.com", "https://anotherdomain.com"},
		AllowedOrigins:   []string{"http://localhost:5173"}, // You can configure specific origins here
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	router.Use(corsConfig.Handler) // Use the cors handler
	router.Use(middleware.Middleware)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	router.Handle("/query", server)

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
