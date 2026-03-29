package cmd

import (
	"gofiber-hax/internal/app"
	"gofiber-hax/internal/infra/config"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations and ensure indexes",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		return app.Migrate(cfg)
	},
}
