# KAKEIBOKU

Kakeiboku is a personal finance management platform. The `apiportal` service is a Go (Gin) backend with JWT authentication, PostgreSQL integration, and a modular architecture for extensibility.

## Project Structure

- `apps/apiportal/` — Main Go backend service
- `docker/` — Docker Compose and environment files
- `Makefile` — Project-wide build and utility commands
- `README.md` — This documentation

## Quick Start

### Prerequisites
- [Go 1.25+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

### 1. Clone the Repository
```bash
git clone <your-repo-url>
cd kakeiboku
```

### 2. Environment Setup
Edit `apps/apiportal/app.env` as needed. Example:
```
PORT=8080
APP_DOMAIN=127.0.0.1
APP_ENV=local
DB_URI=postgresql://root:password@localhost:5432/mydb?sslmode=disable
SECRET_KEY=... # set your own
ACCESS_TOKEN_DURATION=24h
REFRESH_TOKEN_DURATION=168h
```

### 3. Build and Run with Makefile
Build, test, and run locally:
```bash
make all        # build and test
make run        # run the app
make test       # run all tests
make itest      # run integration tests
make clean      # clean build artifacts
```

### 4. Run with Docker Compose
From the project root:
```bash
docker compose -f docker/docker-compose.dev.yaml up --build
```
Or use the Makefile shortcut:
```bash
make docker-run
```
The API will be available at http://localhost:8080

To stop containers:
```bash
make docker-down
```

## Testing

Unit and integration tests use [testify](https://github.com/stretchr/testify) and GoMock. To run all tests:
```bash
make test
```
For integration tests (DB required):
```bash
make itest
```

## Development Tips

- Use `make watch` for live reload (requires [air](https://github.com/air-verse/air)).
- Ensure `app.env` values are valid Go duration strings (e.g., `24h`, `168h`).
- Docker build context must match the Dockerfile COPY paths.
- PostgreSQL runs in a container by default (see compose file).

## Troubleshooting

- **Build errors:** Check Dockerfile COPY paths and build context.
- **Go duration parse errors:** Use valid duration strings (e.g., `24h`, not `24 hours`).
- **Database connection issues:** Ensure the DB container is running and `DB_URI` matches.
- **Test failures:** Ensure mocks and dependencies are up to date.

## License

MIT