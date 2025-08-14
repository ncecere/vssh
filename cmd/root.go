package cmd

import (
	"fmt"
	"os"

	"vssh/internal/auth"
	"vssh/internal/config"
	"vssh/internal/ssh"
	"vssh/internal/utils"
	"vssh/internal/vault"
	"vssh/pkg/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     *types.Config

	// Version information
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// SetVersionInfo sets the version information for the CLI
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vssh [user@]hostname [command]",
	Short: "SSH with Vault-signed certificates",
	Long: `vssh is a CLI tool that signs SSH keys with HashiCorp Vault and uses 
the signed certificate for SSH authentication. It acts as a wrapper around SSH,
providing seamless certificate-based authentication through Vault.

Examples:
  vssh user@server.com
  vssh user@server.com ls -la
  vssh -p 2222 user@server.com`,
	DisableFlagParsing: false,
	Args:               cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		// Initialize logger
		debug, _ := cmd.Flags().GetBool("debug")
		verbose, _ := cmd.Flags().GetBool("verbose")
		if debug || verbose {
			utils.InitLogger(true)
		} else {
			utils.InitLogger(false)
		}

		logger := utils.GetLogger()
		logger.Debug("Starting vssh")

		// Load configuration
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			logger.Fatalf("Failed to load configuration: %v", err)
		}

		logger.Debugf("Configuration loaded successfully")
		logger.Debugf("Vault address: %s", cfg.Vault.Address)
		logger.Debugf("Auth method: %s", cfg.Vault.AuthMethod)

		// Create Vault client
		vaultClient, err := vault.NewClient(&cfg.Vault)
		if err != nil {
			logger.Fatalf("Failed to create Vault client: %v", err)
		}

		// Create authenticator and ensure we have a valid token
		authenticator := auth.NewAuthenticator(vaultClient, &cfg.Vault, logger)
		if err := authenticator.EnsureAuthenticated(); err != nil {
			logger.Fatalf("Authentication failed: %v", err)
		}

		// Parse SSH target
		target, err := ssh.ParseSSHTarget(args[0])
		if err != nil {
			logger.Fatalf("Invalid SSH target: %v", err)
		}

		logger.Debugf("Parsed SSH target - Username: %s, Hostname: %s", target.Username, target.Hostname)

		// Create SSH signer and ensure certificate
		signer := ssh.NewSigner(vaultClient, cfg, logger)
		certPath, err := signer.EnsureSSHCertificate(target.Username)
		if err != nil {
			logger.Fatalf("Failed to ensure SSH certificate: %v", err)
		}

		logger.Debugf("About to parse SSH arguments: %v", args)

		// Parse SSH arguments
		sshOptions, command, err := ssh.ParseSSHArgs(args)
		if err != nil {
			logger.Fatalf("Failed to parse SSH arguments: %v", err)
		}

		logger.Debugf("SSH options parsed successfully")

		// Get private key path for identity
		privateKeyPath, err := signer.GetPrivateKeyPath(target.Username)
		if err != nil {
			logger.Fatalf("Failed to get private key path: %v", err)
		}
		sshOptions.IdentityFile = privateKeyPath

		logger.Debugf("Private key path: %s", privateKeyPath)

		// Create SSH client and connect
		sshClient := ssh.NewClient(cfg, logger)

		// Validate SSH binary is available
		if err := sshClient.ValidateSSHBinary(); err != nil {
			logger.Fatalf("SSH validation failed: %v", err)
		}

		logger.Debugf("SSH binary validation passed")

		fmt.Printf("Connecting to %s with Vault-signed certificate...\n", args[0])
		logger.Infof("Using certificate: %s", certPath)
		logger.Infof("Using private key: %s", privateKeyPath)

		// Execute SSH connection
		logger.Debugf("About to execute SSH connection")
		if err := sshClient.Connect(target, certPath, sshOptions, command); err != nil {
			logger.Fatalf("SSH connection failed: %v", err)
		}

		logger.Debugf("SSH connection completed successfully")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("vssh %s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Built: %s\n", date)
		},
	})

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/vssh/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug output")

	// SSH-compatible flags
	rootCmd.Flags().StringP("port", "p", "", "port to connect to on the remote host")
	rootCmd.Flags().StringP("identity", "i", "", "selects a file from which the identity (private key) is read")
	rootCmd.Flags().BoolP("force-protocol-version1", "1", false, "forces ssh to try protocol version 1 only")
	rootCmd.Flags().BoolP("force-protocol-version2", "2", false, "forces ssh to try protocol version 2 only")
	rootCmd.Flags().BoolP("ipv4", "4", false, "forces ssh to use IPv4 addresses only")
	rootCmd.Flags().BoolP("ipv6", "6", false, "forces ssh to use IPv6 addresses only")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
			os.Exit(1)
		}

		// Search config in XDG config directory
		configDir := fmt.Sprintf("%s/.config/vssh", home)
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("debug") {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
		}
	}
}
