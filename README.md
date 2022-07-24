## Snippets API
---
Try out Snippets hosted by Fitant here: https://ioctl.ml

Or, if you prefer, use the API directly: https://api.ioctl.ml

```
> echo "Hello, World!" > hello-world.txt 

> # just use curl --upload-file

> curl --upload-file hello-world.txt https://api.ioctl.ml

{"URL":"https://api.ioctl.ml/r/EmbossChemicals"}

> curl https://api.ioctl.ml/r/EmbossChemicals
Hello, World!

> echo "i/o/ctl is awesome!"
i/o/ctl is awesome!
```
---

### Features:
- Store and fetch encrypted and end-to-end encrypted snippets
- Snippets get saved against an easy to remember, id like `HedgingSmitten` 
- Snippets are compressed using [zlib](https://en.wikipedia.org/wiki/Zlib)
- Use S3 as storage backend along with global CDN
- Send snippet as formatted JSON, E2E Encrypted JSON or raw body
- Use `POST` and `PUT` at endpoint `/` to save snippet
- Use `POST` and `PUT` at endpoint `/e2e` to end-to-end encypted snippet
- Get snippet using `GET` at endpoint `/<ID>` or `/r/<ID>` or directly from S3 / CDN
- Snippets are by default ephemeral and stored for 7 days

---

### Quick Start:
***Prerequisites: AWS S3, Docker, Docker Compose and curl***
- `git clone --depth=1 https://github.com/fitant/snippets-api`
- `cd snippets-api`
- edit quickstart.env and add your AWS S3 details (currently tied to AWS)
- `docker compose up -d`
- Upload a snippet
    - `curl --upload-file dev.env http://localhost:8080/`
    - Copy the URL field returned in JSON
- Fetch snippet
    - `curl <URL you copied>`
- `docker compose down`

### Setup for Development:
- `git clone --depth=1 https://github.com/fitant/snippets-api`
- `cd snippets-api`
- `go mod download`
- edit dev.env and add your AWS S3 details (currently tied to AWS)
- Start application server
  - `env $(cat dev.env | xargs -L 1) go run src/main.go`

### Deployment / Self-Hosting:
The easiest way to self-host is to simply run an instance of [`realsidsun/snippets-api`](https://hub.docker.com/repository/docker/realsidsun/snippets-api) on a server, and reverse-proxy it after creating a S3 Bucket on AWS and a programatic access key and secret with the following permissions:

```json
"Action": [
    "s3:PutObject",
    "s3:GetObject",
    "s3:GetObjectAttributes",
    "s3:ListBucket",
    "s3:PutObjectAcl"
]
```

***NOTE:*** Look at the Config section down below before you deploy

---

### Cryptographic Specification
- ***The Cryptographic specification is defined [here](https://github.com/Sid-Sun/cryptography/blob/c3234ab9e4b71fcee2a927a7a4f19a663b8d3e8c/specifications/ioctl.md)***

---

### Config
Configuration is done through environment variables

#### General Config:

| Name      | Type / Options                     | Description                         | Required | Default |
|-----------|------------------------------------|-------------------------------------|----------|---------|
| ENV       | string                             | Application Environment             | no       | dev     |
| LOG_LEVEL | debug / info / warn / error        | Log Level to print                  | no       | debug   |
| OVERRIDES | comma and colon seperated mappings | override certain IDs for About, etc | no       |         |

***Example Overrides:*** About:BackwashLicorice,PrivacyPolicy:TranceUnsterile

#### Cryptographic Config:

| Name                   | Type / Options | Description                           | Required | Default |
|------------------------|----------------|---------------------------------------|----------|---------|
| SALT                   | string         | Common SALT used for ID Derivation    | yes      |         |
| ARGON2_ID_MEMORY       | number         | ARGON2 ID Memory / space param in MB  | no       | 32      |
| ARGON2_ID_ROUNDS       | number         | ARGON2 ID rounds / iterations param   | no       | 32      |
| ARGON2_ID_PARALLELISM  | number         | ARGON2 ID parallelism param           | no       | 12      |
| ARGON2_KEY_MEMORY      | number         | ARGON2 KEY Memory / space param in MB | no       | 64      |
| ARGON2_KEY_ROUNDS      | number         | ARGON2 KEY rounds / iterations param  | no       | 12      |
| ARGON2_KEY_PARALLELISM | number         | ARGON2 KEY parallelism param          | no       | 16      |

#### AWS S3 Config:

| Name           | Type   | Description                        | Required |
|----------------|--------|------------------------------------|----------|
| AWS_ACCESS_KEY | string | AWS Programmatic Access Key / ID   | yes      |
| AWS_SECRET_KEY | string | Associated Programmatic Secret Key | yes      |
| AWS_REGION     | string | AWS Hosting Region                 | yes      |
| AWS_S3_BUCKET  | string | S3 Bucket Name                     | yes      |

#### HTTP Server Config:

| Name               | Type / Options           | Description                                   | Required | Default               |
|--------------------|--------------------------|-----------------------------------------------|----------|-----------------------|
| HTTP_LISTEN_HOST   | string                   | HTTP Server listen host                       | no       | 127.0.0.1             |
| HTTP_LISTEN_PORT   | number                   | Replica Set name if using replicaset instance | no       | 8080                  |
| HTTP_CORS_LIST     | comma seperated strings  | Allowed HTTP cross origins list               | no       | http://localhost:*    |
| HTTP_BASE_URL      | string                   | HTTP/S frontend URL to use for formatting     | no       | http://localhost:8080 |
| HTTP_API_ENDPOINT  | string                   | API mount Endpoint from base                  | no       | /snippets             |
| HTTP_RETURN_FORMAT | json / raw               | Default URI for URL to created snippet        | no       | raw                   |
