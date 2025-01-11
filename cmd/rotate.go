package cmd

import (
	"github.com/fbufler/ssh-tunnel-setup/config"
	"github.com/fbufler/ssh-tunnel-setup/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RotateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate key pair on client or target",
		Long:  "Rotating the key pair on the client or target side",
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.Rotate(config.Rotate())
		},
	}

	cmd.Flags().StringP("key-name", "k", "", "Key name")
	cmd.Flags().StringP("key-directory", "d", "", "Key directory")
	cmd.Flags().StringP("key-user", "u", "", "Key user")
	cmd.Flags().StringP("server-name", "s", "", "Server name")
	cmd.Flags().IntP("server-port", "p", 22, "Server port")
	cmd.Flags().StringP("server-user", "U", "", "Server user")
	cmd.Flags().Bool("debug", false, "Debug")

	viper.BindPFlag("rotate.key_name", cmd.Flags().Lookup("key-name"))
	viper.BindPFlag("rotate.key_directory", cmd.Flags().Lookup("key-directory"))
	viper.BindPFlag("rotate.key_user", cmd.Flags().Lookup("key-user"))
	viper.BindPFlag("rotate.server_name", cmd.Flags().Lookup("server-name"))
	viper.BindPFlag("rotate.server_port", cmd.Flags().Lookup("server-port"))
	viper.BindPFlag("rotate.server_user", cmd.Flags().Lookup("server-user"))
	viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))

	return cmd
}
