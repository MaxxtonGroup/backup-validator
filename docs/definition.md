# Test Defintion

```yaml
tests:
- name: <string>                  # Name of the test (required)
  format: <string>                # Format of the backup, possible options: file, mongo, postgresql (required)

  restic:                         # Restore the backup using Restic (required)
    repository: <string>          # Location of the Restic respoistory (required)
    password: <string>            # Use a password to open the Restic repository (note: this is an insecure option, use 'passwordFile' instead)
    passwordFile: <string>        # Use a password file to open the Restic repository
    env: <map>                    # Key-value pair to pass environment variables to the Restic CLI.

  importOptions: <string[]>       # Additional arguments to pass to the restore command of the 'format' provider.

  docker:                         # Use a Docker container to import the backup into a database servier
    image: <string>               # Docker image to use (required, only optional for the 'file' format)
    environment:                  # Pass environment variables to the Docker container
    - "<key>=<value>"
    readyCheck: <string[]>        # Add a command to check when the Docker container is fully started up and ready to import data

  asserts:                        # List of asserts that validate if the backup is valid
    - backupRetention:            # Uses the Retic provider to validate the backup retention
        snapshots: <number>       # Least amount of available snapshots
        olderThan: <duration>     # Least age of the oldest snapshot

    - filesExists: <string[]>     # Glob patterns to validate that certain files exists

    - fileModified:               # Check the modified date of a certain file
        file: <string>            # Glob pattern to find a file to check
        newerThan: <duration>     # Max age of the last modification to the file

    - databasesExists: <string[]> # List of databases that should exists

    - databaseSize:
        database: <string>        # Name of the database
        size: <string>            # Minimal size in bytes of the database (eg. 120mB)

    - tablesExists:
        database: <string>        # Name of the database
        tables: <string[]>        # List of table names that should exists

```