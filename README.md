# Import-Tools

Import-Tools is a web tool for migrating Jira data to ONES system.

## Quick Start

### Binary

Download the latest release from [GitHub Releases](https://github.com/BangWork/import-tools/releases)

```bash
cd scripts
./start.sh
```

### Docker

```bash
docker run -d \
  --name import-tools \
  -p 5000:5000 \
  -v /var/atlassian/application-data/jira:/var/atlassian/application-data/jira \
  ghcr.io/bangwork/import-tools:latest
```

If you are using a shared disk, you can bind it via Docker Volume:

```bash
docker run -d \
  --name import-tools \
  -p 5000:5000 \
  -v /var/atlassian/application-data/jira:/var/atlassian/application-data/jira \
  -v /path/to/shared_disk:/var/import_tools/shared_disk \
  ghcr.io/bangwork/import-tools:latest
```

Then visit the website: http://localhost:5000

## License


