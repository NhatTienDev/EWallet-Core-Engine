package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/nhattiendev/ewallet/docs"
	userHandler "github.com/nhattiendev/ewallet/internal/user/delivery"
	userRepository "github.com/nhattiendev/ewallet/internal/user/infrastructure"
	userUseCase "github.com/nhattiendev/ewallet/internal/user/usecase"
	walletHandler "github.com/nhattiendev/ewallet/internal/wallet/delivery"
	walletRepository "github.com/nhattiendev/ewallet/internal/wallet/infrastructure"
	walletUseCase "github.com/nhattiendev/ewallet/internal/wallet/usecase"
	"github.com/nhattiendev/ewallet/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	bePort := os.Getenv("BE_PORT")
	dbURL := os.Getenv("DB_URL")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	mailpitSMTPHost := os.Getenv("MAILPIT_SMTP_HOST")
	mailpitSMTPPort := os.Getenv("MAILPIT_SMTP_PORT")
	if bePort == "" || dbURL == "" || jwtSecretKey == "" || mailpitSMTPHost == "" || mailpitSMTPPort == "" {
		log.Fatal("Error: Missing required environment variables (BE_PORT, DB_URL, JWT_SECRET_KEY, SMTP_HOST, SMTP_PORT)")
	}

	// Initialize PostgreSQL connection
	db, err := sql.Open("postgres", dbURL) // Initialize configuration, not connected yet
	if err != nil {
		log.Fatalf("Error: Failed to connect to DB configuration: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error: Failed to ping DB, please check password and DB status: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	// Initialize shared middleware
	authMiddleware := middleware.AuthMiddleware([]byte(jwtSecretKey))

	userCreatedChan := make(chan int64, 100) // Queue holds up to 100 events

	mailpitSenderRepository := userRepository.NewMailpitSenderRepository(mailpitSMTPHost, mailpitSMTPPort)
	
	uRepository := userRepository.NewUserRepository(db)
	uUseCase := userUseCase.NewUserUseCase(uRepository, mailpitSenderRepository, jwtSecretKey, userCreatedChan)
	uHandler := userHandler.NewUserHandler(uUseCase)

	wRepository := walletRepository.NewWalletRepository(db)
	wUseCase := walletUseCase.NewWalletUseCase(wRepository)
	wHandler := walletHandler.NewWalletHandler(wUseCase)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		
		walletUseCase.WalletWorker(ctx, userCreatedChan, wUseCase)
	}()

	// go func() {
	// 	<-ctx.Done()
	// 	log.Println("Shutting down background worker...")

	// 	close(userCreatedChan)
	// }()

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
		httpSwagger.URL("http://localhost:"+bePort+"/swagger/doc.json"),
	))

	uHandler.RegisterUserRoutes(r, authMiddleware)
	wHandler.RegisterWalletRoutes(r, authMiddleware)

	server := &http.Server{
		Addr:    ":" + bePort,
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port %s...", bePort)
		log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", bePort)
		log.Printf("Mailpit UI available at http://localhost:%s", os.Getenv("MAILPIT_UI_PORT"))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error: Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Graceful shutdown initiated...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error: Failed to shutdown server gracefully: %v", err)
	}

	log.Println("Waiting for all background workers...")

	wg.Wait()

	log.Println("All background workers stopped")

	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Application shutdown completed")
}