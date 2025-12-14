# Development Workflow

## Overview

1. `gh issue view {id}` — read acceptance criteria
2. `git checkout -b {id}-keywords` — create branch
3. Copy Gherkin to `tests/features/`
4. Write Playwright step definitions that fail
5. Implement until tests pass
6. Run tests and linter
7. Ask: "Anything else or squash-merge?"
8. Squash merge to main, push
9. Delete branches: `git branch -d {branch}` and `git push origin --delete {branch}`
10. `gh issue close {id}`

## GitHub CLI

```bash
gh issue list                 # List open issues
gh issue view 4               # Read issue details
gh issue view 4 --web         # Open in browser
gh issue close 4              # Close after merge
gh issue create --title "..." --body "..."
```

## TDD with Playwright

1. **Read issue** — `gh issue view {id}`
2. **Gherkin first** — Copy scenarios to `.feature` file
3. **Step definitions** — Write Playwright steps
4. **Watch tests fail** — Red phase
5. **Implement** — Backend + frontend until green
6. **Refactor** — Clean up, keep tests green

## Git

**Branch naming:** `{issue-id}-keywords`

```bash
git checkout main && git pull
git checkout -b 4-scan-api
```

**Commit format:**
```
<type>(<scope>): <subject>

Closes #{issue-id}
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`

**Before merge:**
```bash
# Start server and run tests
docker compose up -d
bun test
docker compose down
```

## Issue Templates

### Feature

**Title:** `As a [role] I want [action] so that [benefit]`

**Body:**
```markdown
[1–3 sentences why needed]

**Acceptance Criteria**

\`\`\`gherkin
Feature: [name]

  Scenario: [name]
    Given [context]
    When [action]
    Then [result]
\`\`\`
```

### Bug

**Title:** `Bug: [description]`

**Body:**
```markdown
**Steps to reproduce**
1. ...

**Expected:** ...
**Actual:** ...
```

## CI (GitHub Actions)

On push to `main`:

1. **Test** — lint, unit tests, Playwright tests
2. **Build** (if tests pass) — builds Docker image

Distribution via Docker image — runs anywhere containers run.

### Run CI Locally with `act`

Save GitHub Actions minutes by running workflows locally first:

```bash
# Install (requires Docker)
brew install act          # macOS
choco install act-cli     # Windows

# Run all workflows
act

# Run specific job
act -j test

# List available jobs
act -l
```

On first run, choose "Medium" image (~500MB) — sufficient for most workflows.

### If CI Fails

1. `git checkout -b hotfix-{description}`
2. Fix and test locally
3. Squash merge to main
