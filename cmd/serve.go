package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gofiber-hax/internal/app"
	"gofiber-hax/internal/infra/config"
	"gofiber-hax/internal/infra/logs"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "start",
	Short: "Start HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		logger := logs.New(cfg.Log)
		appInstance, err := app.Build(cfg, logger)
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		errCh := make(chan error, 1)
		go func() {
			errCh <- appInstance.HTTP.Start()
		}()

		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = appInstance.HTTP.Shutdown(shutdownCtx)
			if appInstance.Close != nil {
				_ = appInstance.Close(shutdownCtx)
			}
			return nil
		case err := <-errCh:
			return err
		}
	},
}
