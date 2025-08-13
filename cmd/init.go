package cmd

import (
	"fmt"
	"os"

	"vssh/internal/config"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize vssh configuration",
	Long: `Initialize vssh by creating a default configuration file.

This command creates a default configuration file at ~/.config/vssh/config.yaml
with example settings that you can customize for your environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.GetConfigPath()

		// Check if config file already exists
		if _, err := os.Stat(configPath); err == nil {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("Configuration file already exists at %s\n", configPath)
				fmt.Println("Use --force to overwrite the existing configuration")
				return
			}
		}

		// Create the default configuration
		if err := config.CreateDefaultConfig(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating configuration file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration file created at %s\n", configPath)
		fmt.Println("\nPlease edit the configuration file to match your Vault setup:")
		fmt.Printf("  - Set vault.address to your Vault server URL\n")
		fmt.Printf("  - Configure your preferred authentication method\n")
		fmt.Printf("  - Add user configurations as needed\n")
		fmt.Printf("\nFor more information, see: https://github.com/ncecere/vssh\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Add force flag to overwrite existing config
	initCmd.Flags().BoolP("force", "f", false, "overwrite existing configuration file")
}
