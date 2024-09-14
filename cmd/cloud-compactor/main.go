package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lerenn/cloud-compactor/pkg/cloudcompactor"
	"github.com/spf13/cobra"
)

var (
	path     string
	address  string
	user     string
	password string
)

var compactorCmd = &cobra.Command{
	Use:     "cloud-compactor",
	Version: "0.0.1",
	Short:   "cloud compactor is a simple CLI to compact videos on cloud.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if address == "" {
			return fmt.Errorf("address is required")
		}

		var config cloudcompactor.Config
		if strings.HasPrefix(address, "ftps://") {
			config.FTP.Address = strings.TrimPrefix(address, "ftps://")
			config.FTP.User = user
			config.FTP.Password = password
		} else {
			return fmt.Errorf("address must start with ftps://")
		}

		config.Path = path

		return cloudcompactor.New(config).Run()
	},
}

func init() {
	compactorCmd.Flags().StringVarP(&path, "location", "l", "", "Path to compact")
	compactorCmd.Flags().StringVarP(&address, "address", "a", "", "Server address")
	compactorCmd.Flags().StringVarP(&user, "user", "u", "", "Username")
	compactorCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
}

func main() {
	if err := compactorCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
