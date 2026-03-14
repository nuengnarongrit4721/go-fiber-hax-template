package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "gofiber-hax",
	Short: "Hexagonal GoFiber starter with manual DI",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
}
