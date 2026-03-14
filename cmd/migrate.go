package cmd

import "github.com/spf13/cobra"

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations (stub)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
