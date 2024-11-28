package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shipex",
	Short: "A simple cli tool for automating deployment of Go applications to a VPS",
	Long:  "Shipex is a cli tool  for automating deployment of Go applications to a VPS",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
