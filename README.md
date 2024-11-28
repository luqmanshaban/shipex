# shipex CLI Documentation

`shipex` is a command-line tool designed to streamline the deployment process for Go projects using Git and SSH. It automates tasks such as pushing local changes to a remote Git repository, SSH-ing into a VPS, pulling the latest changes, installing dependencies, running tests, and re-running the app.

## Installation

To install `shipex`, download the precompiled binary from the release page or build it from source:

### From source:

```bash
git clone https://github.com/luqmanshaban/shipex.git
cd shipex
go install
```

## Usage

```bash
shipex <command> [flags]
```

## Available Commands

### `deploy`

The `deploy` command deploys your local Git repository to a remote VPS, automating the process of pushing local changes, SSH-ing into the VPS, pulling the latest code, and performing various operations such as installing dependencies, running tests, and restarting the app.

#### Syntax:

```bash
shipex deploy <directory> [flags]
```

#### Arguments:

* `<directory>`: The path to your local Git repository to be deployed.

#### Flags:

* `-m`, `--message`: The commit message for the changes (default: "Committing latest changes").
* `-f`, `--file`: A list of files to commit. If not specified, all changes will be staged.

#### Example:

```bash
shipex deploy /path/to/local/repo -m "Deploy latest changes" -f "file1.txt","file2.txt"
## or 
go run main.go deploy /path/to/local/repo - m "Deploy latest changes" - f "file1.text"
## or 
go run . deploy /path/to/local/repo # this will commit all push all changes
```

This will commit the changes, push them to the remote repository, and deploy the changes to the VPS.


### `log`

The `log` command retrieves and displays logs for a specific application running on a specified port. It searches for the application's directory within the `~/Documents` folder, verifies the existence of a `log.log` file, and tails the log output.

#### Syntax:

```bash
shipex log <port>
```

#### Arguments:

* `<port>`: The port number of the application whose logs you want to view.

#### Example:

```bash
shipex log 3333
```

This command:

1. Searches for the folder containing the application associated with the port `3333` in the `~/Documents` directory.
2. Checks if a `log.log` file exists in the identified folder.
3. Tails the `log.log` file, streaming the logs in real time.

---

### Detailed Flow of `log`

1. **Search Application Directory** :

* The tool scans the `~/Documents` directory for files containing the specified port.
* It identifies the folder path where the application (`app.yaml`) resides.

1. **Verify Log File** :

* Once the folder is located, it checks for the existence of a `log.log` file within the folder.

1. **Stream Logs** :

* If the `log.log` file is found, the tool tails the log file and streams its output to the terminal in real time.

---

## Notes:

* Ensure that the `log.log` file is present in the application's directory.
* If the port or logs are not found, appropriate error messages will be displayed to help debug the issue.

#### Common Errors:

* **No folder found for the port** : Ensure the specified port is configured in an application's file within `~/Documents`.
* **Log file not found** : Verify that a `log.log` file exists in the located application directory.

---

---

## Detailed Flow of `deploy`

1. **Commit Local Changes** :

   If files are specified with the `-f` flag, those files will be staged and committed with the provided commit message. If no files are specified, all changes in the working directory will be staged and committed.
2. **Push Changes to Git Remote** :

   After committing the changes, the tool pushes the commit to the remote Git repository.
3. **SSH into the VPS** :

   Using the provided SSH credentials (from environment variables `HOST`, `PORT`, `USERNAME`, `PASSWORD`), the tool establishes an SSH connection to the VPS.
4. **Pull Latest Changes on VPS** :

   The tool checks if the repository exists on the VPS. If it exists, it performs a `git pull` to fetch the latest changes. If the repository doesn't exist, it clones it from the remote Git repository.
5. **Install Dependencies** :

   The tool runs the installation command (as specified in the configuration `app.yaml` file on the VPS).
6. **Run Tests** :

   After installing dependencies, the tool will run tests on the VPS to verify that everything is working as expected.
7. **Re-run the Application** :

   Finally, the application is restarted using the specified command.

## NOTE -/

- This tool only works on a local repo that has a remote origin.
- The github authentication is done via SSH.

---

## Environment Variables

The following environment variables must be set in your environment for `shipex` to connect to the VPS:

* `HOST`: The IP address or domain name of the VPS.
* `PORT`: The SSH port of the VPS (typically `22`).
* `USERNAME`: The SSH username to log into the VPS.
* `PASSWORD`: The SSH password to log into the VPS.

For example, you can set these variables in your `.bashrc` or `.zshrc` file:

```bash
export HOST="your-vps-ip"
export PORT="22"
export USERNAME="your-ssh-username"
export PASSWORD="your-ssh-password"
```

---

## Examples

### Deploy with default commit message:

```bash
shipex deploy /path/to/local/repo
```

This command will commit all changes in the specified repository, push them to the remote Git repository, and deploy the latest changes to the VPS.

### Deploy with a custom commit message:

```bash
shipex deploy /path/to/local/repo -m "Bugfix: Corrected issue with login"
```

This will use "Bugfix: Corrected issue with login" as the commit message.

### Deploy with specific files to commit:

```bash
shipex deploy /path/to/local/repo -f "file1.txt" "file2.txt"
```

This will only commit the changes made to `file1.txt` and `file2.txt`.

---

## Error Handling

If any errors occur during the execution of the `deploy` command, they will be displayed with an error message indicating where the failure happened (e.g., during SSH connection, `git pull`, etc.). Common error scenarios include:

* Incorrect SSH credentials.
* Missing or invalid `app.yaml` file.
* Issues with Git configuration or remote repository.

---

## License

`shipex` is released under the MIT License. See [LICENSE](https://chatgpt.com/LICENSE) for more information.
