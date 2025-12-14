# NetMap

Network topology discovery tool. Web UI to scan networks and display discovered devices in real-time.

## Features (MVP)

- **Web UI** — Start scans, view results in browser
- **Live Updates** — Devices appear as they're discovered via WebSocket
- **Mock Mode** — Test and develop without real network
- **Device Display** — Type icons, hostname, vendor info

## Quick Start

```bash
# Start server in mock mode (for testing)
bun docker:up

# Open http://localhost:9090
```

## Development

See [CLAUDE.md](CLAUDE.md) for development workflow.

```bash
# Read an issue
gh issue view 4

# Create branch
git checkout -b 4-scan-api

# Full dev cycle: rebuild + test
bun dev

# Or step by step
bun docker:rebuild
bun pw
```

## Project Structure

```
src/                     Application source code
web/                     Frontend (HTML, CSS, JS)
tests/
  features/              Gherkin specs
  steps/                 Playwright step definitions
docs/                    Documentation
config.json              Server configuration
Dockerfile               Container build
docker-compose.yml       Production compose
docker-compose.dev.yml   Development/mock mode overlay
```

## Documentation

- [CLAUDE.md](CLAUDE.md) — Development guide
- [docs/api.md](docs/api.md) — API reference
- [docs/testing.md](docs/testing.md) — Testing guide
- [docs/workflow.md](docs/workflow.md) — Git workflow
- [docs/architecture.md](docs/architecture.md) — Concurrency, data flow
- [docs/drivers.md](docs/drivers.md) — Device drivers
- [docs/topology.md](docs/topology.md) — Topology tree
- [docs/discovery.md](docs/discovery.md) — Subnet discovery
- [docs/cache.md](docs/cache.md) — Caching strategy

## Tech Stack

- **Runtime:** Docker (backend language is your choice)
- **Frontend:** Vanilla JS
- **Tests:** Playwright + playwright-bdd

## License

MIT
