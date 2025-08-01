// @title Restaurant API Gateway
// @version 1.0
// @description API Gateway для ресторанной системы
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите "Bearer" за которым следует пробел и JWT токен.

package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/netscrawler/Restaurant_is/gate/internal/clients"
	"github.com/netscrawler/Restaurant_is/gate/internal/config"
	"github.com/netscrawler/Restaurant_is/gate/internal/handlers"
	metricsapp "github.com/netscrawler/Restaurant_is/gate/internal/metrics"
	"github.com/netscrawler/Restaurant_is/gate/internal/middleware"
	"github.com/netscrawler/Restaurant_is/gate/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

// @Router /health [get].
func healthCheck(c *gin.Context) {
	// Пример использования трейсинга
	telemetryObj, ok := c.MustGet("telemetry").(*telemetry.Telemetry)
	if ok {
		ctx, span := telemetryObj.StartSpan(c.Request.Context(), "healthCheck")
		defer span.End()
		telemetryObj.RecordMetric("health_check_called", 1)

		c.Request = c.Request.WithContext(ctx)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "gate",
		"version": "1.0.0",
	})
}

func setupTelemetry(cfg *config.TelemertyConfig, logger *slog.Logger) func(ctx context.Context) {
	// Инициализация телеметрии
	telemetryObj, err := telemetry.New(cfg, logger)
	if err != nil {
		panic(err)
	}

	metr := metricsapp.New(logger, telemetryObj, cfg.MetricsPort)

	go func() {
		metr.MustRun()
	}()

	return func(ctx context.Context) {
		err := telemetryObj.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()
	logger := setupLogger(prod)

	// Инициализация gRPC клиентов
	clientsMap := map[string]string{
		"auth":  cfg.GetServiceAddress("auth"),
		"user":  cfg.GetServiceAddress("user"),
		"menu":  cfg.GetServiceAddress("menu"),
		"order": cfg.GetServiceAddress("order"),
	}

	grpcClients, err := clients.NewGRPCClients(clientsMap)
	if err != nil {
		log.Fatalf("gRPC clients error: %v", err)
	}

	defer grpcClients.Close()

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(grpcClients.AuthClient)

	// Инициализация handlers
	authHandler := handlers.NewAuthHandler(grpcClients.AuthClient)
	userHandler := handlers.NewUserHandler(grpcClients.UserClient)
	menuHandler := handlers.NewMenuHandler(grpcClients.MenuClient)
	orderHandler := handlers.NewOrderHandler(grpcClients.OrderClient)

	teleShutdown := setupTelemetry(&cfg.Telemetry, logger)
	defer teleShutdown(context.Background())

	// Инициализация gin
	r := gin.Default()
	// Включаем автоматический трейсинг HTTP-запросов через otelgin
	r.Use(otelgin.Middleware(cfg.Telemetry.ServiceName))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(middleware.LoggingMiddleware(logger))
	// Кастомный CORS: разрешаем Authorization, любые методы, любые заголовки, любые origin
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", healthCheck)

	// Отдача swagger статических файлов
	r.Static("/swagger-static", "./static/swagger")

	r.GET(
		"/swagger/combined/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL("http://localhost:8080/swagger-static/combined-api.json"),
		),
	)

	// Публичные роуты (не требуют авторизации)
	public := r.Group("/api/v1")
	{
		// Auth роуты
		auth := public.Group("/auth")
		{
			auth.POST("/client/login/init", authHandler.LoginInit)
			auth.POST("/client/login/confirm", authHandler.LoginConfirm)
			auth.POST("/staff/login", authHandler.LoginStaff)
			auth.POST("/staff/register", authHandler.RegisterStaff)
			auth.POST("/yandex/login", authHandler.LoginYandex)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/validate", authHandler.Validate)
		}

		// Публичные роуты для меню
		menu := public.Group("/menu")
		{
			menu.GET("/dishes", menuHandler.ListDishes)
			menu.GET("/dishes/:id", menuHandler.GetDish)
		}
	}

	// Защищенные роуты (требуют авторизации)
	protected := r.Group("/api/v1")
	protected.Use(authMiddleware.ValidateToken())
	{
		// User роуты
		user := protected.Group("/users")
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
		}

		// Order роуты
		order := protected.Group("/orders")
		{
			order.POST("/", orderHandler.CreateOrder)
			order.GET("/", orderHandler.ListOrders)
			order.GET("/:id", orderHandler.GetOrder)
		}

		// Admin роуты (требуют роль admin)
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			// User management
			admin.GET("/users", userHandler.ListUsers)
			admin.POST("/image", menuHandler.GetUploadURL)
			admin.POST("/users", userHandler.CreateUser)
			admin.PUT("/users/:id", userHandler.UpdateUser)
			admin.DELETE("/users/:id", userHandler.DeleteUser)

			// Staff management
			admin.GET("/staff", userHandler.ListStaff)
			admin.PUT("/staff/:id", userHandler.UpdateStaff)
			admin.POST("/staff/:id/roles", userHandler.AssignRole)
			admin.DELETE("/staff/:id/roles", userHandler.RevokeRole)

			// Menu management
			admin.POST("/menu/dishes", menuHandler.CreateDish)
			admin.PUT("/menu/dishes/:id", menuHandler.UpdateDish)
			admin.DELETE("/menu/dishes/:id", menuHandler.DeleteDish)

			// Order management
			admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
		}
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		teleShutdown(context.Background())
		os.Exit(0)
	}()

	log.Printf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)

	if err := r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		panic(err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case dev, prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
