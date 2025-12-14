# CLAUDE.md — NetMap

Network topology discovery tool. MVP: web UI to run network scans with live device discovery via WebSocket.

See @README.md for project overview.

## Project Map

```text
src/                     Application source code
web/                     Frontend (HTML, CSS, JS)
tests/
  features/              Gherkin .feature files
  steps/                 Playwright step definitions
docs/                    Documentation
config.json              Server configuration
Dockerfile               Container build
docker-compose.yml       Production compose
docker-compose.dev.yml   Development/mock mode overlay
```

## Commands

```bash
# Docker
docker compose up -d              # Start server
docker compose down               # Stop server
docker build -t netmap .          # Build image

# Mock mode
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Tests (server must be running)
bun test
bun test --headed                 # With visible browser

# CI locally (requires act + Docker)
act -j test
```

## GitHub CLI

```bash
gh issue list                     # List open issues
gh issue view 4                   # View issue #4 details
gh issue view 4 --comments        # Include comments
gh issue create --title "..." --body "..."
gh issue close 4                  # Close issue #4
gh issue comment 4 --body "..."   # Add comment
```

Use `gh issue view <N>` to read the full acceptance criteria before starting work.

## TDD Workflow

**Before coding:**
```bash
gh issue view {id}
git checkout main && git pull && git checkout -b {issue-id}-keywords
```

**Workflow:**

1. Copy acceptance criteria from issue to `tests/features/`
2. Write step definitions that fail
3. Commit failing tests
4. Implement until green
5. Refactor

**Verify tests are real before final commit:**

```bash
git stash          # Remove implementation
bun test           # Must FAIL
git stash pop      # Restore
bun test           # Must pass
```

**After completing:** Run tests and linter, then ask:
> "Anything else or should I squash-merge to main?"

**After merge:** Close issue and delete branches:

```bash
git checkout main && git pull
git branch -d {branch}
git push origin --delete {branch}
gh issue close {id}
```

## Key Patterns

- **Dependency injection** — Scanner interface allows mock/real implementations
- **Streaming results** — Scanner emits devices to stream/queue, WebSocket forwards to UI
- **Single scan at a time** — Server rejects concurrent scans with 409

## Do Not

- Store passwords in code or logs
- Create circular dependencies between modules
- Over-specify implementation details in tests
- Use `waitForTimeout()` in Playwright — use auto-wait or `expect().toPass()`

## Documentation

**Before implementing, read the relevant doc:**

- API endpoints or HTTP handlers → read @docs/api.md first
- Writing tests or step definitions → read @docs/testing.md first  
- Git workflow, commits, branches → read @docs/workflow.md first
- Concurrency patterns → read @docs/architecture.md first
- Adding device vendor support → read @docs/drivers.md first
- Topology tree changes → read @docs/topology.md first
- Subnet auto-discovery → read @docs/discovery.md first
- Cache implementation → read @docs/cache.md first
