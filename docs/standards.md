# Developer Standards

This page defines the commit and pull request standards for contributions. These standards are enforced by CI.

## Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.

### Format

```
<type>(<scope>): <subject>


<body>


<footer>
```

- The subject is required and must be in imperative mood.
- The scope is optional.
- The body is optional except when explaining non-obvious changes.
- The footer is required when introducing breaking changes.

### Types

- feat
- fix
- docs
- style
- refactor
- test
- chore
- CI

### Breaking Changes

Breaking changes must include a footer using this format:
`BREAKING CHANGE: <description>`

## Pull Requests

PRs should be small and focused.

### PR Titles

Must follow the same Conventional Commits format as commit messages.

### PR Descriptions

Should include:

1. What the change does
2. Why the change is needed
3. How it is implemented
4. Testing steps

### Branch Naming

Use the same prefix as the commit type e.g. `feat/<short-description>`

## Formatting and Linting

Contributions must follow formatting and linting rules.

### Frontend

- **Prettier** handles code formatting for frontend files.
- **ESLint** is used for linting.

### Backend

- **gofmt** handles code formatting for Go files.
- **golangci-lint** is used for linting.

## Pre-Commit Hooks

- **Husky** is used to run formatting and linting checks before each commit.
- Commit messages are validated with **Commitlint**.

## Continuous Integration

- CI enforces formatting and linting for frontend and backend.
- Commit messages are verified with **Commitlint**.
