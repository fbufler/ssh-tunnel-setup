package cmd

import (
	"log/slog"
	"os/user"

	"github.com/fbufler/ssh-tunnel-setup/config"
	"github.com/fbufler/ssh-tunnel-setup/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TargetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Setup client",
		Long:  "Setting up the client side of the tunnel",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := internal.ClientSetup(config.Client())
			if err != nil {
				return err
			}
			storeTunnelConfig(config.Client())
			internal.SetupTunnel(config.Tunnel())
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "client name")
	cmd.Flags().StringP("key-name", "k", "", "Key name")
	cmd.Flags().StringP("key-directory", "d", "", "Key directory")
	cmd.Flags().StringP("key-user", "u", "", "Key user")
	cmd.Flags().StringP("server-name", "s", "", "Server name")
	cmd.Flags().IntP("server-port", "p", 22, "Server port")
	cmd.Flags().StringP("server-user", "U", "", "Server user")
	cmd.Flags().StringP("server-pass", "P", "", "Server password")
	cmd.Flags().StringP("server-key-name", "K", "", "Server key name")
	cmd.Flags().Bool("debug", false, "Debug")

	viper.BindPFlag("client.name", cmd.Flags().Lookup("name"))
	viper.BindPFlag("client.key_name", cmd.Flags().Lookup("key-name"))
	viper.BindPFlag("client.key_directory", cmd.Flags().Lookup("key-directory"))
	viper.BindPFlag("client.key_user", cmd.Flags().Lookup("key-user"))
	viper.BindPFlag("client.server_name", cmd.Flags().Lookup("server-name"))
	viper.BindPFlag("client.server_port", cmd.Flags().Lookup("server-port"))
	viper.BindPFlag("client.server_user", cmd.Flags().Lookup("server-user"))
	viper.BindPFlag("client.server_pass", cmd.Flags().Lookup("server-pass"))
	viper.BindPFlag("client.server_key_name", cmd.Flags().Lookup("server-key-name"))
	viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))

	return cmd
}

func storeTunnelConfig(cfg *config.ClientConfig) {
	currentTunnel := config.UnsafeTunnel()
	if currentTunnel == nil {
		currentTunnel = &config.TunnelConfig{}
	}
	if currentTunnel.SSHConfigPath == "" {
		currentTunnel.SSHConfigPath = cfg.KeyDirectory + "/config"
	}
	if currentTunnel.KeyDirectory == "" {
		currentTunnel.KeyDirectory = cfg.KeyDirectory
	}
	if currentTunnel.ServerKeyName == "" {
		currentTunnel.ServerKeyName = cfg.KeyName
	}
	if currentTunnel.ServerUser == "" {
		currentTunnel.ServerUser = cfg.ServerUser
	}
	if currentTunnel.ServerName == "" {
		currentTunnel.ServerName = cfg.ServerName
	}
	if currentTunnel.LocalUser == "" {
		user, err := user.Current()
		if err != nil {
			slog.Warn("Could not get current user: %v", err)
		} else {
			currentTunnel.LocalUser = user.Username
		}
	}
	config.StoreTunnelConfig(currentTunnel)
}
