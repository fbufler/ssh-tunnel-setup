package cmd

import (
	"github.com/fbufler/ssh-tunnel-setup/config"
	"github.com/fbufler/ssh-tunnel-setup/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Setup server",
		Long:  "Setting up the server side of the tunnel",
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.ServerSetup(config.Server())
		},
	}

	cmd.Flags().StringP("name", "n", "", "Server name")
	cmd.Flags().StringP("tunnel-user", "u", "", "Tunnel user")
	cmd.Flags().StringP("tunnel-pass", "p", "", "Tunnel password")
	cmd.Flags().StringP("sshd-config-path", "c", "", "Path to sshd config")
	cmd.Flags().StringP("sshd-config-backup-path", "b", "", "Path to sshd config backup")
	cmd.Flags().Bool("debug", false, "Debug")

	viper.BindPFlag("server.name", cmd.Flags().Lookup("name"))
	viper.BindPFlag("server.tunnel_user", cmd.Flags().Lookup("tunnel-user"))
	viper.BindPFlag("server.tunnel_pass", cmd.Flags().Lookup("tunnel-pass"))
	viper.BindPFlag("server.sshd_config_path", cmd.Flags().Lookup("sshd-config-path"))
	viper.BindPFlag("server.sshd_config_backup_path", cmd.Flags().Lookup("sshd-config-backup-path"))
	viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))

	return cmd
}
