// Swagger / OpenAPI docs for TravelKu Backend.
//
// @title           TravelKu Backend API
// @version         1.0
// @description     REST API for TravelKu backend (Auth & Users).
//
// @contact.name    TravelKu
//
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Use: Bearer <token>
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authfeature "backend-travelku/internal/auth"
	bookingfeature "backend-travelku/internal/booking"
	"backend-travelku/internal/config"
	"backend-travelku/internal/middleware"
	paketfeature "backend-travelku/internal/packs"
	userfeature "backend-travelku/internal/user"
	jwt "backend-travelku/pkg/auth"
	"backend-travelku/pkg/logger"	

	"github.com/gin-gonic/gin"

	_ "backend-travelku/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize JWT with secret from env
	jwt.InitJWT(cfg.JWTSecret)

	// Initialize logger
	if err := logger.InitLogger(cfg.AppEnv); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
		os.Exit(1)
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&userfeature.User{},
		&paketfeature.Paket{},
		&bookingfeature.Booking{},
	)
	if err != nil {
		logger.Fatal("Failed to run Gorm AutoMigrate", zap.Error(err))
		os.Exit(1)
	}

	// Initialize repositories
	userRepo := userfeature.NewRepository(db)
	authRepo := authfeature.NewRepository(userRepo) // ← pass userRepo (composition)
	paketRepo := paketfeature.NewRepository(db)
	bookingRepo := bookingfeature.NewRepository(db)

	// Initialize services
	authService := authfeature.NewService(authRepo)
	userService := userfeature.NewService(userRepo)
	paketService := paketfeature.NewService(paketRepo)
	bookingService := bookingfeature.NewService(bookingRepo, paketRepo)

	// Initialize controllers
	authController := authfeature.NewController(authService)
	userController := userfeature.NewController(userService)
	paketController := paketfeature.NewController(paketService)
	bookingController := bookingfeature.NewController(bookingService)

	// Setup router
	router := setupRouter(cfg, authController, userController, paketController, bookingController)

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ServerTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerTimeout) * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server starting on port " + cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed")
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown")
		os.Exit(1)
	}

	logger.Info("Server shutdown successfully")
}

// setupRouter mengatur route aplikasi
func setupRouter(cfg *config.Config, authController *authfeature.Controller, userController *userfeature.Controller, paketController *paketfeature.Controller, bookingController *bookingfeature.Controller) *gin.Engine {
	// Set environment
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware global
	router.Use(middleware.CORSMiddleware(cfg.CorsAllowOrigins))
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    cfg.AppName,
		})
	})

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ==================== PUBLIC ROUTES ====================

	// ==================== AUTH ROUTES ====================
	authfeature.RegisterRoutes(router, authController)

	// ==================== USER ROUTES ====================
	userfeature.RegisterRoutes(router, userController)

	// ==================== PAKET ROUTES ====================
	paketfeature.RegisterRoutes(router, paketController)

	// ==================== BOOKING ROUTES ====================
	bookingfeature.RegisterRoutes(router, bookingController)

	return router
}

