# Testing Environment

This page explains how to set up and run a test environment for local end-to-end testing.

## Development Environment

It is assumed that you have followed the `Development Environment` instructions and are able to run the app locally in that enviroment.

The testing environment is separate from the development enviroment. It uses its own local database. The testing environment can be run alongside the development environment. The services use different ports. Because we are managing two different Docker Compose environments, there are many scripts defined in `package.json` to assist with Docker commands. Scripts defined here can build, start, stop, and restart containers as well as display logs. Some of these scripts are shown below.

## Steps

### 1. Environment Variables

Create a `.env.test` file in the root. You may use `.env.test.example` as a template.

These are the necessary environment variables.

```
POSTGRES_TEST_USER=test_user
POSTGRES_TEST_PASSWORD=test_password
POSTGRES_TEST_DB=test_db
POSTGRES_TEST_PORT=5434
TEST_DATABASE_URL=postgres://test_user:test_password@db:5432/test_db?sslmode=disable
FRONTEND_PORT=3001
BACKEND_PORT=4000
JWT_SECRET=your_jwt_secret
API_URL=http://localhost:4000
```

You can choose a user, password, database name, and port for your database. You do not need to create a database yourself. One will be created in the `test-db` container when you start the environment. These variables simply configure that database.

`TEST_DATABASE_URL` is used by the `test-migrate` and `test-backend` containers to connect to the database. Use the information from the `TEST_POSTGRES` variables for the "your_user", "your-password", and "your_db" items. The database can be accessed within the Docker network at `test-db:5432` as you see above.

You can specify the ports to start the frontend and backend on. If you change the backend port, make sure you update `API_URL` accordingly.

`JWT_SECRET` is the key used to sign and verify JWTs. This should be a long, random string. It must be kept secret.

`API_URL` is used by the frontend to connect to the backend. By default, the backend is on port 8080.

### 2. Start the Testing Environment

Build and start the containers.

```
npm run test:build
npm run test:start
```

This starts:

- PostgreSQL test database on port `5434`
- Backend server on port `4000`
- Frontend server on port `3000`

The first time, database migrations will run automatically.

### 3. Populate database with data

If this is the first time starting the development environment, the database migrations will be automatically applied but the database tables will be empty. To populate the database, use the script `backend/cmd/update_db/main.go`. This will update the `nodes`, `ways`, and `way_nodes` tables with the most recent OpenStreetMap data while keeping existing reviews. It works with an empty database, so it can be used to populate the database for the first time.

```
// Ensure DATABASE_URL environment variable is set using localhost:[LOCAL_PORT]
// Do not use db:[DOCKER_PORT]. This only works to address the database from within one of the Docker containers. It will not work on localhost, which is where this script will run
// You can see what port the database maps to on localhost in docker-compose.test.yml
export DATABASE_URL=postgres://test_user:test_password@localhost:5434/test_db?sslmode=disable
cd backend/cmd/update_db
go run main.go
```

You should then restart the backend so that it uses the new data.

```
npm run test:br
```

### 4. Test that services are starting correctly

View the logs for the frontend and backend containers.

```
npm run test:fl
npm run test:bl
```

### 5. Run tests

End-to-end tests are in the `e2e` directory. They can be run with `npx playwright test`.

```
cd e2e
npx playwright test
```

### 6. Stop the environment

To stop all containers:

```
npm run test:stop
```
