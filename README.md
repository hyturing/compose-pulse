# cpulse

[![CI](https://github.com/hyturing/compose-pulse/actions/workflows/ci.yml/badge.svg)](https://github.com/hyturing/compose-pulse/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hyturing/compose-pulse)](https://goreportcard.com/report/github.com/hyturing/compose-pulse)
[![GitHub release](https://img.shields.io/github/v/release/hyturing/compose-pulse)](https://github.com/hyturing/compose-pulse/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A real-time terminal UI that visualizes your Docker Compose startup sequence as an interactive dependency tree — so you can see exactly which service is blocking everything else.

```
  ● api-gateway
  ├─ ● auth-service
  │  └─ ● postgres
  ├─ ◐ payment-service        ← still starting
  │  └─ ● redis
  └─ ✕ worker                 ← exited with error
     └─ ● rabbitmq
```

Press `Enter` on any node to inspect its logs. Press `/` to filter by regex.

---

## The Problem

`docker compose up` dumps interleaved logs from every container into one stream. When a service fails, you stop everything, scroll through thousands of lines, open multiple terminal tabs, and lose 10 minutes hunting the root cause.

**cpulse** makes the dependency graph visual and live — healthy services stay green, failing ones turn red immediately, and their logs are one keypress away.

---

## Installation

### Homebrew (macOS / Linux)

```sh
brew tap hyturing/cpulse
brew install cpulse
```

### Download a binary

Grab the latest release for your platform from [GitHub Releases](https://github.com/hyturing/compose-pulse/releases).

```sh
# macOS (Apple Silicon)
curl -L https://github.com/hyturing/compose-pulse/releases/latest/download/cpulse_darwin_arm64.tar.gz | tar xz
sudo mv cpulse /usr/local/bin/
```

### Build from source

Requires Go 1.22+.

```sh
git clone https://github.com/hyturing/compose-pulse.git
cd compose-pulse
make build
# binary at ./bin/cpulse
```

---

## Usage

Run `cpulse` from the same directory as your `docker-compose.yml`:

```sh
cpulse
```

Or point to a specific file:

```sh
cpulse --file path/to/docker-compose.yml
```

`cpulse` reads your compose file, builds the dependency graph, and starts watching the Docker daemon. Bring up your stack in another terminal as usual:

```sh
docker compose up
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `--file` | auto-detect | Path to the compose file |
| `--version` | — | Print version and exit |

---

## Keyboard Shortcuts

| Key | Action |
|---|---|
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `Enter` | Open / close log modal for selected service |
| `/` | Filter logs by regex |
| `Esc` | Close search bar / close modal |
| `q` / `Ctrl+C` | Quit |

---

## Status Indicators

| Symbol | Color | Meaning |
|---|---|---|
| `●` | Green | Running and healthy |
| `◐` | Yellow | Starting / health check running |
| `●` | Red | Health check failed / unhealthy |
| `○` | Gray | Waiting on a dependency |
| `✕` | Red | Exited with error |

---

## Requirements

- Docker Desktop or Docker Engine running locally
- A `docker-compose.yml` using v3+ syntax with `depends_on` and/or `healthcheck` blocks
- macOS or Linux

---

## Roadmap

- **v0.2** — Switch from polling to `docker events` streaming (lower latency)
- **v0.3** — Full `jq`-style log filtering via `gojq`
- **v0.4** — Real graph canvas with box-drawing edges
- **v0.5** — Container lifecycle actions (restart, stop)

---

## Contributing

Pull requests are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for setup instructions and guidelines.

## License

[MIT](LICENSE) — © 2026 hyturing
