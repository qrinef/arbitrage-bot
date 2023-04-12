package cmd

import (
	"github.com/qrinef/arbitrage-bot/services/container"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "arbitrage-bot-cli",
		Run: func(cmd *cobra.Command, args []string) {
			containerService := container.NewService()
			containerService.Start()

			containerService.PoolsService.Start()
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}
