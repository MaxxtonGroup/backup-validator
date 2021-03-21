# backup-validator
backup-validator is a CLI for validating Restic/Elasticsearch backups by restoring them.

## Usage
Using the binary:
```shell
backup-validator -f test1.yaml -f test2.yaml
```

With docker:
```shell
docker run --rm -v $(pwd):/workdir maxxton/backup-validator --test-file=test1.yaml --test-file=test2.yaml
```

## Test definition
See the [full documentation of the test definition](./docs/definition.md).

```yaml
tests:
- name: grafana
  # Format of the backup (supported: file)
  format: file

  # Use a restic repository
  restic:
    repository: s3:s3.amazonaws.com/my-bucket/grafana
    passwordFile: restic-password-file
    # Set environment variables for restic
    env:
      AWS_ACCESS_KEY_ID: XXX
      AWS_SECRET_ACCESS_KEY: XXX

  # Validate the backup repository
  asserts:
    # Validate the retention of the backup
    - backupRetention:
        snapshots: 4
        olderThan: 96h # 4 days

    # Validate if certain files exists after restoring the backup
    - filesExists:
      - /var/lib/grafana/grafana.db

    # Validate the modification time of a certian file
    - fileModified:
        file: /var/lib/grafana/grafana.db
        newerThan: 48h # 2 days
```

## Installation

Prerequisites: [Docker](https://www.docker.com/) and [Restic](https://restic.net/).

**Linux**
```shell
BACKUP_VALIDATOR_VERSION=<VERSION>
wget https://github.com/MaxxtonGroup/backup-validator/releases/download/v${BACKUP_VALIDATOR_VERSION}/backup-validator_${BACKUP_VALIDATOR_VERSION}_Linux_x86_64.tar.gz
tar -zxf backup-validator_${BACKUP_VALIDATOR_VERSION}_Linux_x86_64.tar.gz
rm backup-validator_${BACKUP_VALIDATOR_VERSION}_Linux_x86_64.tar.gz
mv backup-validator bin/backup-validator
chmod +x bin/backup-validator
```
