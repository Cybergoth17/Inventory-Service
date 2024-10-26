package api

import (
	"context"
	"errors"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"inventory-service/internal/config"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	err := godotenv.Load(".env")
	cfg, err := config.ReadEnv()
	if err != nil {
		return err
	}

	//logger
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("zap new: %w", err)
	}
	logger = logger.With(
		zap.String("service_name", "tgw-car-insurance-service"),
	)

	//vault
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = cfg.Vault.Addr
	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		return err
	}

	vaultClient.SetToken(cfg.Vault.Token)

	errCh := make(chan error, 1)

	httpServer := http.Server{
		Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
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
