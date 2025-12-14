# Testing Guide

## Philosophy

**Playwright tests everything.** All acceptance criteria are tested through the real browser UI. This catches what unit tests miss: JS bugs, CSS issues, timing problems, actual user experience.

Unit tests are for internal pure functions (parsers, algorithms) that have no UI.

## Setup

Use `playwright-bdd` to run Gherkin feature files with Playwright.

```bash
bun install
bunx playwright install chromium
bun test
```

The test runner starts the server automatically in mock mode before running tests.

## Project Structure

```
tests/
  features/     # Gherkin .feature files
  steps/        # Step definitions (TypeScript)
```

## Writing Tests

### Feature Files

**Naming convention:** Feature files must be named `{issue-number}-{description}.feature` where the issue number is zero-padded to 4 digits (e.g., `0004-scan-api.feature`).

Copy acceptance criteria from GitHub issues into feature files:

```bash
# View the issue to get acceptance criteria
gh issue view 4
```

```gherkin
# tests/features/0004-scan-api.feature
Feature: Scan API

  Scenario: Start a scan
    Given I am on the scan page
    When I enter network "192.168.1.0/24" and core switch "192.168.1.1"
    And I click Start
    Then a scan begins
```

### Step Definitions

Use `playwright-bdd` pattern with Playwright's `page` fixture:

```typescript
// tests/steps/scan.steps.ts
import { expect } from '@playwright/test';
import { createBdd } from 'playwright-bdd';

const { Given, When, Then } = createBdd();

Given('I am on the scan page', async ({ page }) => {
  await page.goto('/scan');
});

When('I enter network {string} and core switch {string}', async ({ page }, network, coreSwitch) => {
  await page.fill('[name="network"]', network);
  await page.fill('[name="core-switch"]', coreSwitch);
});

When('I click Start', async ({ page }) => {
  await page.click('button:has-text("Start")');
});

Then('a scan begins', async ({ page }) => {
  await expect(page.locator('.scan-status')).toContainText(/scanning/i);
});
```

### Selector Strategy

Prefer semantic selectors that survive refactoring:

```typescript
// Good - semantic
page.locator('button:has-text("Start")')
page.locator('[name="network"]')
page.locator('.device-list tr')

// Avoid - brittle
page.locator('#btn-start-scan-v2')
page.locator('div > div > button:nth-child(2)')
```

Use `data-testid` only when semantic selection isn't possible.

### No Arbitrary Waits

**Never use `waitForTimeout()` or arbitrary delays.** They make tests flaky or slow.

```typescript
// WRONG - arbitrary wait
await page.waitForTimeout(2000);
await page.click('button');

// CORRECT - wait for condition
await page.waitForSelector('.scan-status');
await expect(page.locator('.device-row')).toHaveCount(5);
await page.locator('button:has-text("Start")').click(); // auto-waits
```

Playwright auto-waits for elements before actions. For dynamic content:

```typescript
// Wait for element to appear
await expect(page.locator('.scan-complete')).toBeVisible();

// Wait for text content
await expect(page.locator('.status')).toHaveText('Complete');

// Wait for network idle (use sparingly)
await page.waitForLoadState('networkidle');

// Poll for condition
await expect(async () => {
  const count = await page.locator('.device-row').count();
  expect(count).toBeGreaterThan(0);
}).toPass({ timeout: 10000 });
```

## Running Tests

```bash
bun test                    # All tests
bun test --headed           # See the browser
bun test --debug            # Step through
bun test --grep "Scan"      # Filter by name
```

## Mock Mode

Tests run against server in mock mode (`NETMAP_MOCK=true`). The mock scanner:
- Returns fixture devices from embedded JSON
- Emits devices with configurable delay
- Provides deterministic results

## TDD Workflow

1. Read issue: `gh issue view {id}`
2. Write feature file from acceptance criteria
3. Write step definitions that fail
4. Commit failing tests
5. Implement until green
6. Verify tests are real:
   ```bash
   git stash          # Remove implementation
   bun test           # Must fail
   git stash pop      # Restore
   bun test           # Must pass
   ```
7. After merge: `gh issue close {id}`

## CI

```yaml
- run: bun install && bunx playwright install chromium
- run: docker compose up -d
- run: bun test
- run: docker compose down
```
