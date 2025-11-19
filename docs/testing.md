# Testing

## Frontend UI Tests

The frontend app is in `frontend/`. UI tests (Playwright) live under `frontend/e2e/`.

Run UI tests locally:

```bash
cd frontend
npm install
# run Playwright UI tests
npx playwright test
```

## Backend Unit Tests

The backend lives in `backend/` and contains Go unit tests.

Run all backend unit tests:

```bash
cd backend
go test ./...
```

To run a single package or test, use package path or `-run` as usual.

## End-To-End Tests

Full repository end-to-end instructions and test environment setup are in `docs/test-env.md`. Follow that document to prepare the environment and run E2E tests.
