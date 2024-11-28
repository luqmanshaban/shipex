/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/luqmanshaban/shipex/functions"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Used to print a specific app (server) logs",
	Long: `Used to print a specific app (server) logs by searching for a specific port in the Documents directory.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		LogUsingPort(args[0])
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}




func LogUsingPort(port string) {
	config := SSHConfig{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	commands := []string{
		fmt.Sprintf(`
		#!/bin/bash

# Define the port to search for (e.g., 3333)
PORT=%v

# Search for the port inside files in the ~/Documents directory and capture the folder path
FOLDER_PATH=$(grep -r "PORT: $PORT" ~/Documents/* | grep -oP '.*(?=/app.yaml)' | head -n 1)

# Check if the folder path was found
if [[ -z "$FOLDER_PATH" ]]; then
    echo "No folder found with port $PORT in ~/Documents."
    exit 1
fi

# Check if the log.log file exists in that folder
LOG_FILE="$FOLDER_PATH/log.log"
if [[ ! -f "$LOG_FILE" ]]; then
    echo "log.log file not found in $FOLDER_PATH."
    exit 1
fi

# Output the path of the folder and the log file
echo "Found folder for port $PORT: $FOLDER_PATH"
echo "Tailing log file: $LOG_FILE"

# Tail the log.log file for the specified port
tail -f "$LOG_FILE" 

		`,port),
	}

	err := ConnectToVPS(config, commands)
	if err != nil {
		log.Printf("Error executing command on the VPS: %v", err)
		return
	}

	// Log the outputs
	functions.Success("Commands executed successfully!")
}