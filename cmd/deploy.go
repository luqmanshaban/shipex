/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	// "time"

	"github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing/object"
	gitSSH "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/luqmanshaban/shipex/functions"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// holds the configuration for ssh connetion
type SSHConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

var commitMsg string
var filesToCommit []string

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "`deploy` will deploy the go repository to the VPS",
	Long: `The deploy command performs several steps to deploy the project:
	1. Stage all local changes (or specified files).
	2. Commit the changes to the local repository.
	3. Push the changes to the remote repository.
	4. SSH into the VPS and pull the latest changes.
	5. Install necessary packages.
	6. Run tests and re-run the app.
	7. Broadcast logs for debugging.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		functions.Info("Deploying")
		if commitMsg == "" {
			commitMsg = "Commiting latest changes"
		}

		done := make(chan bool) // channel to signal when gh is done

		go func() {
			Github(args[0], commitMsg, filesToCommit)
			done <- true
		}()
		<-done

		runCommand(args[0])
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")
	deployCmd.Flags().StringVarP(&commitMsg, "message", "m", "", "Commit message for the changes (default: 'Commiting lastest changes')")
	deployCmd.Flags().StringSliceVarP(&filesToCommit, "file", "f", []string{}, "Files to commit (must be one or more)")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Github(dir, cmsg string, files []string) {

	r, err := git.PlainOpen(dir)
	functions.CheckIfError(err)

	w, err := r.Worktree()
	functions.CheckIfError(err)

	// stage all changes
	functions.Info("staging changes")
	if len(files) > 1 {
		for _, file := range files {
			functions.Info("Staging: %s", file)
			_, err = w.Add(file)
			functions.CheckIfError(err)
		}
	} else {
		err = w.AddGlob(".")
		functions.CheckIfError(err)
	}
	functions.Success("Changes staged!")

	// commit
	functions.Info("Commiting Changes")
	_, err = w.Commit(cmsg, &git.CommitOptions{})
	functions.CheckIfError(err)
	functions.Success("Changes Committed!")

	// push
	functions.Info("Pushing to remote repository")
	path := os.Getenv("SSH_KEY_PATH")
	sshKey, err := os.ReadFile(path)
	functions.CheckIfError(err)
	auth, err := gitSSH.NewPublicKeys("git", sshKey, "")
	if err != nil {
		fmt.Printf("Failed to authenticate with github via ssh %v", err)
	}
	functions.CheckIfError(err)

	err = r.Push(&git.PushOptions{
		Auth:     auth,
		Progress: os.Stdout,
	})
	functions.CheckIfError(err)
	functions.Success("Changes pushed to the remote repository!")
}

// Establishes a connection and executes a command
func ConnectToVPS(config SSHConfig, commands []string) error {
	// Defines SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the server
	functions.Attempting("Connecting to the VPS")
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		functions.CheckIfError(err)
		return err
	}
	defer client.Close()
	functions.Success("Connection established!")

	// Execute commands
	for _, command := range commands {

		// Create a session for each command
		session, err := client.NewSession()
		if err != nil {
			fmt.Printf("\x1b[31;1mError creating session: %v", err)
			return err
		}

		// Set up output streams
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		// Execute the command
		if err := session.Run(command); err != nil {
			fmt.Printf("\x1b[31;1mError executing command: %v\nCommand: %s", err, command)
			session.Close() // Ensure session is closed if there's an error
			return err
		}

		// Close the session immediately after the command is executed
		session.Close()
	}

	functions.Success("All commands executed successfully!")
	return nil
}

func runCommand(dir string) {
	config := SSHConfig{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	yamlFile := fmt.Sprintf("%v/app.yaml", dir)
	v, err := getYamlData(yamlFile)
	if err != nil {
		fmt.Printf("Error reading yaml file: %v", err.Error())
	}

	commands := []string{
		// Change directory to the correct location
		fmt.Sprintf(`
		source ~/.bashrc
		cd %[6]s &&
		 if [ -d "%[1]s" ]; then 
            echo "Repository exists. Pulling latest changes..." &&
            cd %[1]s && git pull && %[3]s && %[4]s;
        else 
            echo "Cloning repository...";
            git clone %[2]s && cd %[1]s && %[3]s && %[4]s;
        fi
		
		# Find and kill the process running on port %[7]s (if any)
        echo "Checking for running process on port %[7]s..."
        PID=$(lsof -t -i:%[7]s)
        if [ -n "$PID" ]; then
		echo "Killing process $PID running on port %[7]s..."
		kill -9 $PID
        else
		echo "No process found running on port %[7]s."
        fi
		
		
        # Start the server in the background with nohup
        echo "Server started in the background."
        nohup %[5]s > log.log 2>&1 &
        exit  # Exit the shell session immediately after starting the server
`, v.Name, v.Github, v.Install, v.Build, v.Run, v.Location, v.Port),
	}

	// Execute the command on the VPS
	err = ConnectToVPS(config, commands)
	if err != nil {
		log.Printf("Error executing command on the VPS: %v", err)
		return
	}

	// Log the outputs
	functions.Success("Commands executed successfully!")
}

func getYamlData(path string) (functions.Yaml, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return functions.Yaml{}, err
	}

	c, err := functions.ReadYaml(f)
	if err != nil {
		return functions.Yaml{}, err
	}

	return c, nil
}
