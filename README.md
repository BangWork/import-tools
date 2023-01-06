# Import-Tools

Import-Tools is a web tool for migrating Jira data to ONES system.

## Quick Start

```bash
docker run -d \
  --name import-tools \
  -p 5000:5000 \
  -v /var/atlassian/application-data/jira:/var/atlassian/application-data/jira \
  ghcr.io/bangwork/import-tools:latest
```

Then visit the website: http://localhost:5000

## License
