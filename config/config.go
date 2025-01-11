package config

import (
	"fmt"
	"log/slog"

	"github.com/fbufler/ssh-tunnel-setup/package/system"
	"github.com/spf13/viper"
)

const TunnelUser = "tunneluser"

type ClientConfig struct {
	Name          string `mapstructure:"name"`
	KeyName       string `mapstructure:"key-name"`
	KeyDirectory  string `mapstructure:"key-directory"`
	KeyUser       string `mapstructure:"key-user"`
	ServerName    string `mapstructure:"server-name"`
	ServerPort    int    `mapstructure:"server-port"`
	ServerUser    string `mapstructure:"server-user"`
	ServerPass    string `mapstructure:"server-pass"`
	ServerKeyName string `mapstructure:"server-key-name"`
}

func (c *ClientConfig) required() error {
	// KeyName KeyDirectory KeyUser ServerName ServerPort ServerUser ServerPass or ServerKeyName
	missingFields := []string{}
	if c.KeyName == "" {
		missingFields = append(missingFields, "KeyName")
	}
	if c.KeyDirectory == "" {
		missingFields = append(missingFields, "KeyDirectory")
	}
	if c.KeyUser == "" {
		missingFields = append(missingFields, "KeyUser")
	}
	if c.ServerName == "" {
		missingFields = append(missingFields, "ServerName")
	}
	if c.ServerPort == 0 {
		missingFields = append(missingFields, "ServerPort")
	}
	if c.ServerUser == "" {
		missingFields = append(missingFields, "ServerUser")
	}
	if c.ServerPass == "" && c.ServerKeyName == "" {
		missingFields = append(missingFields, "ServerPass or ServerKeyName")
	}
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required parameters: %v", missingFields)
	}
	return nil
}

type RotateConfig struct {
	KeyName      string `mapstructure:"key-name"`
	KeyDirectory string `mapstructure:"key-directory"`
	KeyUser      string `mapstructure:"key-user"`
	ServerName   string `mapstructure:"server-name"`
	ServerPort   int    `mapstructure:"server-port"`
	ServerUser   string `mapstructure:"server-user"`
}

func (c *RotateConfig) required() error {
	// KeyName KeyDirectory KeyUser ServerName ServerPort ServerUser
	missingFields := []string{}
	if c.KeyName == "" {
		missingFields = append(missingFields, "KeyName")
	}
	if c.KeyDirectory == "" {
		missingFields = append(missingFields, "KeyDirectory")
	}
	if c.KeyUser == "" {
		missingFields = append(missingFields, "KeyUser")
	}
	if c.ServerName == "" {
		missingFields = append(missingFields, "ServerName")
	}
	if c.ServerPort == 0 {
		missingFields = append(missingFields, "ServerPort")
	}
	if c.ServerUser == "" {
		missingFields = append(missingFields, "ServerUser")
	}
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required parameters: %v", missingFields)
	}
	return nil
}

type ServerConfig struct {
	Name                 string `mapstructure:"name"`
	TunnelUser           string `mapstructure:"tunnel-user"`
	TunnelPass           string `mapstructure:"tunnel-pass"`
	SSHDConfigPath       string `mapstructure:"sshd-config-path"`
	SSHDConfigBackupPath string `mapstructure:"sshd-config-backup-path"`
}

func (c *ServerConfig) required() error {
	// TunnelUser SSHDConfigPath SSHDConfigBackupPath
	missingFields := []string{}
	if c.TunnelUser == "" {
		missingFields = append(missingFields, "TunnelUser")
	}
	if c.SSHDConfigPath == "" {
		missingFields = append(missingFields, "SSHDConfigPath")
	}
	if c.SSHDConfigBackupPath == "" {
		missingFields = append(missingFields, "SSHDConfigBackupPath")
	}
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required parameters: %v", missingFields)
	}
	return nil
}

type TunnelConfig struct {
	HostIdentifier string `mapstructure:"host-identifier"`
	SSHConfigPath  string `mapstructure:"ssh-config-path"`
	KeyDirectory   string `mapstructure:"key-directory"`
	ServerKeyName  string `mapstructure:"server-key-name"`
	LocalUser      string `mapstructure:"local-user"`
	LocalHost      string `mapstructure:"local-host"`
	LocalPort      int    `mapstructure:"local-port"`
	ServerUser     string `mapstructure:"server-user"`
	ServerName     string `mapstructure:"server-name"`
	ServerPort     int    `mapstructure:"server-port"`
}

func (c *TunnelConfig) required() error {
	// SSHConfigPath HostIdentifier ServerName LocalUser KeyDirectory ServerKeyName LocalHost LocalPort RemotePort
	missingFields := []string{}
	if c.SSHConfigPath == "" {
		missingFields = append(missingFields, "SSHConfigPath")
	}
	if c.HostIdentifier == "" {
		missingFields = append(missingFields, "HostIdentifier")
	}
	if c.ServerName == "" {
		missingFields = append(missingFields, "ServerName")
	}
	if c.LocalUser == "" {
		missingFields = append(missingFields, "LocalUser")
	}
	if c.KeyDirectory == "" {
		missingFields = append(missingFields, "KeyDirectory")
	}
	if c.ServerKeyName == "" {
		missingFields = append(missingFields, "ServerKeyName")
	}
	if c.LocalHost == "" {
		missingFields = append(missingFields, "LocalHost")
	}
	if c.LocalPort == 0 {
		missingFields = append(missingFields, "LocalPort")
	}
	if c.ServerUser == "" {
		missingFields = append(missingFields, "ServerUser")
	}
	if c.ServerPort == 0 {
		missingFields = append(missingFields, "ServerPort")
	}
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required parameters: %v", missingFields)
	}
	return nil
}

type Config struct {
	Client         ClientConfig `mapstructure:"client"`
	Server         ServerConfig `mapstructure:"server"`
	Rotate         RotateConfig `mapstructure:"rotate"`
	Tunnel         TunnelConfig `mapstructure:"tunnel"`
	Debug          bool         `mapstructure:"debug"`
	TrustedHostKey string       `mapstructure:"trusted-host-key"`
}

var AppConfig Config

func setDefaults() {
	homeDir, _ := system.HomeDir()
	if homeDir == "" {
		slog.Warn("Home directory not found, using current directory")
		homeDir = "."
	}
	viper.SetDefault("client.name", "default-client")
	viper.SetDefault("client.key_name", "default-key")
	viper.SetDefault("client.key_directory", fmt.Sprint(homeDir, "/.ssh"))
	viper.SetDefault("client.key_user", "default-client-user")
	viper.SetDefault("client.user", "default-user")
	viper.SetDefault("client.server_name", "localhost")
	viper.SetDefault("client.server_port", 8080)

	viper.SetDefault("server.name", "default-server")
	viper.SetDefault("server.tunnel_user", TunnelUser)
	viper.SetDefault("server.sshd_config_path", "/etc/ssh/sshd_config")
	viper.SetDefault("server.sshd_config_backup_path", "/etc/ssh/sshd_config.bak")

	viper.SetDefault("tunnel.ssh_config_path", fmt.Sprint(homeDir, "/.ssh/config"))
	viper.SetDefault("tunnel.host_identifier", "default-host")
	viper.SetDefault("tunnel.key_directory", fmt.Sprint(homeDir, "/.ssh"))
	viper.SetDefault("tunnel.server_key_name", "default-key")
	viper.SetDefault("tunnel.local_host", "localhost")
	viper.SetDefault("tunnel.local_port", 3306)
	viper.SetDefault("tunnel.remote_port", 3306)
}

func LoadConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // or viper.SetConfigType("json")
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error(fmt.Sprintf("Error reading config file, %s", err))

	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		slog.Error(fmt.Sprintf("Error unmarshalling config, %s", err))
	}
}

func Server() *ServerConfig {
	serverCfg := &AppConfig.Server
	if serverCfg.required() != nil {
		slog.Error(fmt.Sprintf("Error reading server config, %s", serverCfg.required()))
		panic("Required server config missing")
	}
	return serverCfg
}

func Client() *ClientConfig {
	clientCfg := &AppConfig.Client
	if clientCfg.required() != nil {
		slog.Error(fmt.Sprintf("Error reading client config, %s", clientCfg.required()))
		panic("Required client config missing")
	}
	return clientCfg
}

func Rotate() *RotateConfig {
	rotateCfg := &AppConfig.Rotate
	if rotateCfg.required() != nil {
		slog.Error(fmt.Sprintf("Error reading rotate config, %s", rotateCfg.required()))
		panic("Required rotate config missing")
	}
	return rotateCfg
}

func Tunnel() *TunnelConfig {
	tunnelCfg := &AppConfig.Tunnel
	if tunnelCfg.required() != nil {
		slog.Error(fmt.Sprintf("Error reading tunnel config, %s", tunnelCfg.required()))
		panic("Required tunnel config missing")
	}
	return tunnelCfg
}

func UnsafeTunnel() *TunnelConfig {
	return &AppConfig.Tunnel
}

func StoreRotationConfig(cfg *RotateConfig) error {
	AppConfig.Rotate = *cfg
	viper.Set("rotate", AppConfig.Rotate)
	return viper.WriteConfig()
}

func StoreTunnelConfig(cfg *TunnelConfig) error {
	AppConfig.Tunnel = *cfg
	viper.Set("tunnel", AppConfig.Tunnel)
	return viper.WriteConfig()
}

func Debug() bool {
	return AppConfig.Debug
}

func TrustedHostKey() string {
	return AppConfig.TrustedHostKey
}
