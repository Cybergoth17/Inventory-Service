package api

import (
	"context"
	"errors"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"inventory-service/internal/api"
	"inventory-service/internal/api/inventory/handlers"
	"inventory-service/internal/config"
	"inventory-service/internal/repository/inventory"
	inventoryclient "inventory-service/internal/service/inventory"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	cfg, err := config.ReadEnv()
	if err != nil {
		return err
	}

	// Logger
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("zap new: %w", err)
	}
	logger = logger.With(zap.String("service_name", "tgw-car-insurance-service"))

	// Context
	ctx := context.Background()

	// Vault
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = cfg.Vault.Addr
	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		return err
	}
	vaultClient.SetToken(cfg.Vault.Token)

	// Redis
	redisClient := config.InitRedis(net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port), cfg.Redis.Password, 0, ctx)

	// Mongo
	mongoClient := config.InitMongoDB(fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port), cfg.Mongo.DbName, cfg.Mongo.CollectionName)

	// Repository
	mongoRepo := inventory.NewMongoRepository(mongoClient.Database(cfg.Mongo.DbName).Collection(cfg.Mongo.CollectionName))
	redisRepo := inventory.NewRedisRepository(redisClient)

	// Inventory Client
	inventoryClient, err := inventoryclient.NewClient(mongoRepo, redisRepo)
	if err != nil {
		return err
	}

	// Inventory Handler
	inventoryHandler := handlers.NewInventoryHandler(inventoryClient)

	// Create API and Router
	api := api.NewAPI(inventoryHandler)
	router := api.NewRouter()

	// Swagger
	http.Handle("/swagger/", httpSwagger.WrapHandler) // Serve Swagger UI
	http.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "../../docs/swagger.json")
	})

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("Server running on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case sig := <-stop:
		logger.Info("graceful shutdown os.Signal", zap.Any("signal", sig))
		if err := httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("http server shutdown: %w", err)
		}
	case err = <-errCh:
		return err
	}

	return nil
}
