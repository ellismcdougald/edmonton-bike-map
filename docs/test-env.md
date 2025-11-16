# Testing Environment

This page explains how to set up and run a test environment for local end-to-end testing.

## Development Environment

It is assumed that you have followed the `Development Environment` instructions and are able to run the app locally in that enviroment.

The testing environment is separate from the development enviroment. It uses its own local database. The testing environment can be run alongside the development environment. The services use different ports. Because we are managing two different Docker Compose environments, there are many scripts defined in `package.json` to assist with Docker commands. Scripts defined here can build, start, stop, and restart containers as well as display logs. Some of these scripts are shown below.

## Steps

### 1. Environment Variables

Create a `.env.test` file in the root. You may use `.env.test.example` as a template.

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
