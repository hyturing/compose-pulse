# Contributing to cpulse

Thank you for taking the time to contribute. This guide covers everything you need to get a change from idea to merged PR.

---

## Development Setup

**Prerequisites:** Go 1.22+, Docker Desktop (or Docker daemon) running locally, `golangci-lint`.

```sh
git clone https://github.com/hyturing/compose-pulse.git
cd compose-pulse

# Install dependencies
go mod download

# Build
make build          # → bin/cpulse

# Run tests
make test

# Lint
make lint
```

To install `golangci-lint`:
```sh
brew install golangci-lint          # macOS
# or
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

---

## Branching Model

| Branch | Purpose |
|---|---|
| `main` | Always releasable — no direct commits |
| `feat/<name>` | New features |
| `fix/<name>` | Bug fixes |
| `docs/<name>` | Documentation-only changes |

Open a PR against `main`. Keep PRs focused — one concern per PR.

---

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add regex filter to log modal
fix: handle services with no healthcheck block
docs: add Homebrew install instructions
refactor: extract state indicator into styles.go
test: add cycle detection unit test
```

Scope is optional but appreciated for larger packages:

```
feat(dag): add Kahn's topological sort
fix(docker): reconnect on socket timeout
```

---

## Pull Request Checklist

Before submitting, confirm:

- [ ] `make build` passes
- [ ] `make test` passes
- [ ] `make lint` passes (no new lint violations)
- [ ] New behavior has test coverage where practical
- [ ] `README.md` updated if user-facing behavior changed
- [ ] The PR description explains **why**, not just what changed

---

## Project Structure

```
internal/compose/   — docker-compose.yml parsing
internal/dag/       — dependency graph (Kahn's algorithm)
internal/docker/    — Docker SDK wrapper + 500ms monitor
internal/ui/        — Bubble Tea model, tree renderer, log modal
cmd/cpulse/         — entry point
```

See [`plans/implementation-plan.md`](plans/) for the full architectural breakdown (this file is gitignored; clone the repo to read it).

---

## Reporting Bugs

Open a [GitHub Issue](https://github.com/hyturing/compose-pulse/issues/new) with:

1. OS, Go version (`go version`), Docker version
2. The `docker-compose.yml` that triggers the issue — redact any secrets
3. Exact error output or a screen recording

---

## Proposing Features

Open an issue with the `enhancement` label **before** writing code. Describe the problem you're solving and why the existing behavior falls short. This avoids duplicated effort and keeps the scope manageable.

Features currently out of scope for v1 (already planned for later):
- `gojq`-style filtering (v0.3)
- Graph canvas layout (v0.4)
- Container lifecycle actions — start, stop, restart (v0.5)
- Kubernetes / Docker Swarm support

---

## Code of Conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). Be kind and constructive.
