package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/lib/pq"
	_ "github.com/nhattiendev/ewallet/docs"
	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/internal/user/usecase"
	"github.com/nhattiendev/ewallet/internal/user/delivery"
	"github.com/nhattiendev/ewallet/internal/user/infrastructure"
)

// @title E-Wallet Core Engine API
// @version 1.0
// @description API Documentation for the E-Wallet Core Application.

// @host localhost:8080
// @BasePath /

// Configure "Authorize" button in Swagger UI
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <token>" to authenticate.
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, reading from system env")
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if port == "" || dbURL == "" || jwtSecretKey == "" {
		log.Fatal("Error: Missing required environment variables (PORT, DB_URL, JWT_SECRET_KEY)")
	}

	// Initialize PostgreSQL connection
	db, err := sql.Open("postgres", dbURL) // Initialize configuration, not connected yet
	if err != nil {
		log.Fatalf("Error: Failed to connect to DB configuration: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error: Failed to ping DB, please check password and DB status: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	// Initialize shared middleware
	authMiddleware := middleware.AuthMiddleware([]byte(jwtSecretKey))

	userRepo := infrastructure.NewUserRepository(db)
	userUC := usecase.NewUserUseCase(userRepo, jwtSecretKey)
	userHandler := delivery.NewUserHandler(userUC)

	// General router configuration
	r := chi.NewRouter()

	// CORS middleware configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	
	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+port+"/swagger/doc.json"),
	))

	userHandler.RegisterUserRoutes(r, authMiddleware)

	log.Printf("Starting server on port %s...", port)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", port)
	serverAddr := ":" + port
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Error: Failed to start server: %v", err)
	}
}