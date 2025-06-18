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
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/netscrawler/Restaurant_is/gate/internal/clients"
	"github.com/netscrawler/Restaurant_is/gate/internal/config"
	"github.com/netscrawler/Restaurant_is/gate/internal/handlers"
	"github.com/netscrawler/Restaurant_is/gate/internal/middleware"
	"github.com/netscrawler/Restaurant_is/gate/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// @Summary Health check
// @Description Проверка состояния сервиса
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthCheck(c *gin.Context) {
	// Пример использования трейсинга
	telemetryObj, ok := c.MustGet("telemetry").(*telemetry.Telemetry)
	if ok {
		ctx, span := telemetryObj.StartSpan(c.Request.Context(), "healthCheck")
		defer span.End()
		telemetryObj.RecordMetric("health_check_called", 1)
		c.Request = c.Request.WithContext(ctx)
	}
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "gate",
		"version": "1.0.0",
	})
}

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()

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

	// Инициализация телеметрии
	telemetryObj, err := telemetry.New(&cfg.Telemetry, slog.Default())
	if err != nil {
		log.Fatalf("failed to init telemetry: %v", err)
	}
	defer telemetryObj.Shutdown(context.Background())

	// Инициализация gin
	r := gin.Default()

	// Включаем автоматический трейсинг HTTP-запросов через otelgin
	r.Use(otelgin.Middleware(cfg.Telemetry.ServiceName))

	// Endpoint для Prometheus метрик
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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

	// // Swagger UI для каждого сервиса
	// r.GET(
	// 	"/swagger/auth/*any",
	// 	ginSwagger.WrapHandler(
	// 		swaggerFiles.Handler,
	// 		ginSwagger.URL("http://localhost:8080/swagger-static/auth/auth.swagger.json"),
	// 	),
	// )
	// r.GET(
	// 	"/swagger/user/*any",
	// 	ginSwagger.WrapHandler(
	// 		swaggerFiles.Handler,
	// 		ginSwagger.URL("http://localhost:8080/swagger-static/user/user.swagger.json"),
	// 	),
	// )
	// r.GET(
	// 	"/swagger/menu/*any",
	// 	ginSwagger.WrapHandler(
	// 		swaggerFiles.Handler,
	// 		ginSwagger.URL("http://localhost:8080/swagger-static/menu/menu.swagger.json"),
	// 	),
	// )
	// r.GET(
	// 	"/swagger/order/*any",
	// 	ginSwagger.WrapHandler(
	// 		swaggerFiles.Handler,
	// 		ginSwagger.URL("http://localhost:8080/swagger-static/order/order.swagger.json"),
	// 	),
	// )
	//
	// // Объединенный Swagger UI
	// r.GET(
	// 	"/swagger/combined/*any",
	// 	ginSwagger.WrapHandler(
	// 		swaggerFiles.Handler,
	// 		ginSwagger.URL("http://localhost:8080/swagger-static/combined-api.json"),
	// 	),
	// )

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
		telemetryObj.Shutdown(context.Background())
		os.Exit(0)
	}()

	log.Printf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
