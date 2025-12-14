# NetMap Tests

This directory contains end-to-end tests using Playwright and playwright-bdd.

## Setup

```bash
# Install dependencies
bun install

# Install Playwright browsers
bunx playwright install chromium
```

## Running Tests

```bash
# Run all tests (starts server automatically)
bun test

# Run with visible browser
bun test --headed

# Run in debug mode
bun test --debug

# Run specific test
bun test --grep "Test Infrastructure"
```

## Writing Tests

1. Create a `.feature` file in `tests/features/`
2. Write Gherkin scenarios
3. Create step definitions in `tests/steps/`
4. Run `bunx bddgen` to generate test files (or just run `bun test`)

## Structure

```
tests/
  features/          # Gherkin .feature files
  steps/             # TypeScript step definitions
  README.md          # This file
```

## Notes

- Tests run against the server in mock mode
- The server is started automatically via `webServer` config in `playwright.config.ts`
- Generated test files are in `.features-gen/` (gitignored)

