package cmd

import (
	"github.com/qrinef/arbitrage-bot/services/container"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(discoveryCmd)
}

var discoveryCmd = &cobra.Command{
	Use: "discovery",
	Run: func(cmd *cobra.Command, args []string) {
		containerService := container.NewService()
		containerService.Start()

		containerService.DiscoveryService.Start()
	},
}
