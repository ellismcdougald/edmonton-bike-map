# Development Environment

This page explains how to set up and run the Edmonton Bike Map web app locally using Docker.

## Prerequisites

Ensure that Docker (including Docker Compose), Go, Node, and Git are installed.

## Testing Environment

The testing environment is separate from the development enviroment. It uses its own local database. The testing environment can be run alongside the development environment. The services use different ports. Because we are managing two different Docker Compose environments, there are many scripts defined in `package.json` to assist with Docker commands. Scripts defined here can build, start, stop, and restart containers as well as display logs. Some of these scripts are shown below.

Set up the development environment and ensure that it works. Then try the testing environment.

## Steps

### 1. Clone the repository

```
git clone git@github.com:ellismcdougald/edmonton-bike-map.git
cd edmonton-bike-map
```

### 2. View 'docker-compose.yml'

Take a look at `docker-compose.yml`. There are four containers: `dev-db`, `dev-migrate`, `dev-frontend`, and `dev-backend`. `dev-migrate` is where the database migrations will run.

### 2. Environment Variables

Create a `.env` file in the root. You may use `.env.example` as a template. These variables allow the backend to connect to the database.

These are the necessary environment variables.

```
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_db
POSTGRES_PORT=5433
DATABASE_URL=postgres://your_user:your_password@dev-db:5432/your_db?sslmode=disable
JWT_SECRET=your_jwt_secret
API_URL=http://localhost:8080
```

You can choose a user, password, database name, and port for your database. You do not need to create a database yourself. One will be created in the `dev-db` container when you start the environment. These variables simply configure that database.

`DATABASE_URL` is used by the `migrate` and `dev-backend` containers to connect to the database. Use the information from the `POSTGRES` variables for the "your_user", "your-password", and "your_db" items. The database can be accessed within the Docker network at `dev-db:5432` as you see above.

`JWT_SECRET` is the key used to sign and verify JWTs. This should be a long, random string. It must be kept secret.

`API_URL` is used by the frontend to connect to the backend. By default, the backend is on port 8080.

### 3. Start the Development Environment

Build and start the containers.

```
npm run dev:build
npm run dev:start
```

This starts:

- PostgreSQL database on port `5433`
- Backend server on port `8080`
- Frontend server on port `5173`

The first time, database migrations will run automatically.

### 4. Populate database with data

If this is the first time starting the development environment, the database migrations will be automatically applied but the database tables will be empty. To populate the database, use the script `backend/cmd/update_db/main.go`. This will update the `nodes`, `ways`, and `way_nodes` tables with the most recent OpenStreetMap data while keeping existing reviews. It works with an empty database, so it can be used to populate the database for the first time.

This script uses the `DATABASE_URL` environment variable. Since this script is not being run within the Docker network, we must address the database using "localhost" rather than "dev-db". You must provide a different `DATABASE_URL` to this script than you use when starting the development environment.

```
// Ensure DATABASE_URL environment variable is set using localhost:[LOCAL_PORT]
// Do not use db:[DOCKER_PORT]. This only works to address the database from within one of the Docker containers. It will not work on localhost, which is where this script will run
// You can see what port the database maps to on localhost in docker-compose.yml
cd backend/cmd/update_db
go run main.go
```

You should then restart the backend so that it uses the new data.

```
npm run dev:br
```

### 5. Access the app

- Frontend: [http://localhost:5173](http://localhost:5173)
- Backend: [http://localhost:8080](http://localhost:8080)

### 6. Stop the environment

To stop all containers:

```
npm run dev:stop
```

## Notes

- Frontend live reload is handled by _Vite_.
- Backend live reload is handled by _Air_.
- Database data is persisted in the `postgres_data` volume, so data is retained between sessions.
