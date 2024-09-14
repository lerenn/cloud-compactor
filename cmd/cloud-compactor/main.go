package main

import (
	"fmt"
	"os"

	"github.com/lerenn/cloud-compactor/pkg/cloudcompactor"
	"github.com/spf13/cobra"
)

var (
	configPath string
)

var compactorCmd = &cobra.Command{
	Use:     "cloud-compactor",
	Version: "1.0.0",
	Short:   "cloud compactor is a simple CLI to compact videos on cloud.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if configPath == "" {
			return fmt.Errorf("config is required")
		}

		config, err := cloudcompactor.LoadConfigFromFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		return cloudcompactor.New(config).Run()
	},
}

func init() {
	compactorCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file")
}

func main() {
	if err := compactorCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
