# Development Environment

This page explains how to set up and run the Edmonton Bike Map web app locally using Docker.

## Prerequisites

Ensure that Docker (including Docker Compose), Go, Node, and Git are installed.

## Steps

### 1. Clone the repository

```
git clone git@github.com:ellismcdougald/edmonton-bike-map.git
cd edmonton-bike-map
```

### 2. Environment Variables

Create a `.env` file in the root. You may use `.env.example` as a template. These variables allow the backend to connect to the database.

### 3. Start the Development Environment

Build and start the containers.

```
docker compose up --build
```

This starts:

- PostgreSQL database on port `5433`
- Backend server on port `8080`
- Frontend server on port `5173`

The first time, database migrations will run automatically.

### 4. Populate database with data

If this is the first time starting the development environment, the database migrations will be automatically applied but the database tables will be empty. To populate the database, use the script `backend/cmd/update_db/main.go`. This will update the `nodes`, `ways`, and `way_nodes` tables with the most recent OpenStreetMap data while keeping existing reviews. It works with an empty database, so it can be used to populate the database for the first time.

```
// Ensure DATABASE_URL environment variable is set using localhost:[LOCAL_PORT]
// Do not use db:[DOCKER_PORT]. This only works to address the database from within one of the Docker containers. It will not work on localhost, which is where this script will run
// You can see what port the database maps to on localhost in docker-compose.yml
cd backend/cmd/update_db
go run main.go
```

You should then restart the backend so that it uses the new data.

```
docker compose restart backend
```

### 5. Access the app

- Frontend: [http://localhost:5173](http://localhost:5173)
- Backend: [http://localhost:8080](http://localhost:8080)

### 6. Stop the environment

To stop all containers:

```
docker compose down
```

To remove volumes and start fresh:

```
docker compose down -v
```

## Notes

- Frontend live reload is handled by _Vite_.
- Backend live reload is handled by _Air_.
- Database data is persisted in the `postgres_data` volume, so data is retained between sessions.
