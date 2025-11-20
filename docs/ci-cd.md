# Continuous Integration / Continuous Deployment

## Continous Integration

This workflow runs on every pull request to main. It consists of four jobs:

- Commitlint
- Frontend CI
- Backend CI
- E2E Tests

These checks ensure that changes maintain clean code and correctness.

### Commitlint

Commitlint checks the commits to ensure that they conform to the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.

### Frontend CI

This job does formatting, linting, and runs the UI tests on the frontend. Formatting uses [prettier](https://prettier.io/). Linting is done by [eslint](https://eslint.org/).

### Backend CI

This job does formatting, linting, and runs the unit tests on the backend. Formatting uses [gofmt](https://pkg.go.dev/cmd/gofmt). Linting is done by [golangci](https://golangci-lint.run/).

### E2E Tests

This job runs the end-to-end tests. It creates a temporary PostgreSQL database and runs the database migrations. It populates the test database with data for the tests. It starts the frontend and backend servers. Finally, it runs the tests with [playwright](https://playwright.dev/).

## Continuous Deployment

Continuous deployment is the process of updating the deployed application so that it reflects the `main` branch. This occurs when a push is made to `main`, such as after a pull request is merged.

### Frontend

Frontend deployments are handled automatically by Vercel when main is updated.

### Backend

The `CI/CD Backend` Github Actions workflow handles backend deployment. The backend is deployed on a DigitalOcean VPS.

#### Workflow

1. Install Ansible. This is used for deployment process.
2. Checkout the repository with the latest code.
3. Set up Go.
4. Run Go tests.
5. Build a production Docker image of the backend using `backend/Dockerfile`.
6. Log in to Github Container Registry.
7. Push the backend image to the GHCR repository.
8. Set up SSH access on the VPS.
9. Run the Ansible playbook for deployment.

#### Ansible Playbook

The Ansible playbook in `deployment/ansible/backend.yml` is responsible for deploying the Go backend on the VPS.

1. Install Docker
2. Install Docker Compose
3. Ensure app directory exists
4. Copy Docker Compose file backend-compose.yml to app directory
5. Log in to GitHub Container Registry
6. Ensure Let's Encrypt directory exists
7. Ensure acme.json exists. This is used for certificate storage
8. Run database migrations
9. Update the production database using the script in `backend/cmd/update_db/main.go"
10. Deploy the backend with Docker Compose
