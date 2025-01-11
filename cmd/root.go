package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fbufler/ssh-tunnel-setup/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ssh-tunnel-setup",
	Short: "ssh-tunnel-setup is a CLI application for setting up an SSH tunnel",
	Long:  `ssh-tunnel-setup is a CLI application for setting up an SSH tunnel`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("ssh-tunnel-setup is a CLI application for setting up an SSH tunnel")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define flags and configuration settings here
	config.LoadConfig()

	// Add subcommands
	rootCmd.AddCommand(ClientCmd())
	rootCmd.AddCommand(ServerCmd())
	rootCmd.AddCommand(TargetCmd())
	rootCmd.AddCommand(RotateCmd())

	// Configure slog
	opts := &slog.HandlerOptions{}
	if config.Debug() {
		opts.Level = slog.LevelDebug
	} else {
		opts.Level = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}
