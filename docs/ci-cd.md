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
