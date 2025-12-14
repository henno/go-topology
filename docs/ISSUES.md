# NetMap Issues

Specifications for the NetMap network topology discovery tool. 

**This file is a seed for GitHub issues.** Create issues in GitHub using `gh issue create`, then work from GitHub as the source of truth.

**Testing approach:** All scenarios run through Playwright against the real browser UI.

**Issue flow:** Issues must be implemented in order — each builds on previous ones.

**REST conventions:** Endpoints use plural nouns (`/scans`, `/devices`).

---

## Issue #1: As a developer I want test infrastructure so that I can write browser tests

Before implementing features, we need Playwright configured to run Gherkin scenarios against the browser.

**Acceptance Criteria**

```gherkin
Feature: Test Infrastructure

  @manual
  Scenario: Test runner works
    When I run "bun test"
    Then Playwright executes feature files from tests/features/
    And step definitions from tests/steps/ are used
```

**Technical decisions:**
- Use `playwright-bdd` to connect Gherkin to Playwright
- Use `bun` as package manager

---

## Issue #2: As a developer I want project scaffolding so that I can start implementing features

The project needs a web server that serves static web files with basic navigation layout.

**Acceptance Criteria**

```gherkin
Feature: Project Foundation

  Scenario: Web UI loads
    Given I open the application
    Then I see the page title "NetMap"
    And I see a sidebar with navigation
    And I see "Scan" and "Devices" links

  Scenario: Navigation works
    Given I am on the home page
    When I click "Scan" in the navigation
    Then I am on the scan page
```

**Technical decisions:**
- Minimal dependencies for HTTP server
- Web assets served from `web/` directory
- Config from JSON file

---

## Issue #3: As a developer I want mock scanner so that I can test without real network

Mock mode lets us develop and test the UI without needing actual network devices.

**Acceptance Criteria**

```gherkin
Feature: Mock Scanner

  Scenario: Mock mode indicator
    Given the server is running in mock mode
    When I open the application
    Then I see a "Mock Mode" indicator
```

**Technical decisions:**
- Scanner interface for dependency injection
- MockScanner returns 5 fixture devices (router, switch, computer, printer, unknown)
- Mock mode activated via env var or CLI flag

---

## Issue #4: As a network admin I want to start scans via API so that the backend can discover devices

HTTP endpoints for scan lifecycle. UI integration comes in Issue #6.

**Acceptance Criteria**

```gherkin
Feature: Scan API

  Scenario: Start scan via API
    When I POST to /scans with network "192.168.1.0/24" and core switch "192.168.1.1"
    Then I receive a scan ID
    And the response status is "scanning"

  Scenario: Get scan status
    Given a scan is running
    When I GET /scans/{id}
    Then I see the scan status and discovered count

  Scenario: Only one scan at a time
    Given a scan is running
    When I POST to /scans
    Then I receive a 409 Conflict error

  Scenario: Cancel scan
    Given a scan is running
    When I DELETE /scans/{id}
    Then the scan is cancelled
```

**Endpoints:**
- `POST /scans` — Start scan
- `GET /scans/current` — Current scan (for reconnection)
- `GET /scans/{id}` — Scan by ID  
- `DELETE /scans/{id}` — Cancel scan

---

## Issue #5: As a network admin I want live updates so that I see devices as they are discovered

WebSocket streaming for real-time progress.

**Acceptance Criteria**

```gherkin
Feature: Live Scan Updates

  Scenario: Devices stream via WebSocket
    Given I connect to /scans/{id}/stream
    When the scanner discovers a device
    Then I receive a WebSocket message with the device

  Scenario: Completion message
    Given I am connected to a scan stream
    When the scan completes
    Then I receive a completion message with total count
```

**Technical decisions:**
- Use WebSocket library appropriate for your language

---

## Issue #6: As a network admin I want a scan page so that I can run scans from the browser

The scan page with form, status display, and live-updating results table.

**Acceptance Criteria**

```gherkin
Feature: Scan UI

  Scenario: Scan form
    Given I am on the scan page
    Then I see a field for network range (CIDR notation, e.g. "192.168.1.0/24")
    And I see a field for core switch IP (the starting point for discovery)
    And I see a Start button

  Scenario: Start scan from UI
    Given I am on the scan page
    When I enter network "192.168.1.0/24" and core switch "192.168.1.1"
    And I click Start
    Then I see the scan status change to "Scanning"

  Scenario: Devices appear as discovered
    Given I started a scan
    Then devices appear in the results table as they are discovered

  Scenario: Scan completion
    Given I started a scan
    When the scan completes
    Then I see "Complete" status
    And I see the total device count

  Scenario: Invalid network shows error
    Given I am on the scan page
    When I enter an invalid network
    And I click Start
    Then I see an error message

  Scenario: Cancel from UI
    Given a scan is running
    When I click Cancel
    Then the scan stops

  Scenario: Empty state
    Given no scan has been run
    Then I see a message prompting to start a scan
```

---

## Issue #7: As a network admin I want to see device details so that I can identify what's on my network

Device table shows type icons, hostname, and vendor with appropriate fallbacks for missing data.

**Acceptance Criteria**

```gherkin
Feature: Device Display

  Scenario: Device table columns
    Given a scan completed
    Then each device row shows: type icon, IP, hostname, vendor

  Scenario: Device type icons
    Given a scan completed
    Then routers show a router icon
    And switches show a switch icon
    And computers show a computer icon
    And unknown devices show a question mark icon

  Scenario: Missing hostname
    Given a device has no hostname
    Then the hostname column shows "—"

  Scenario: Missing vendor  
    Given a device has no vendor
    Then the vendor column shows "Unknown"

  Scenario: All mock devices appear
    Given I run a scan in mock mode
    When the scan completes
    Then I see exactly 5 devices
```

**Note:** Mock fixtures include: 1 router (Routerboard.com), 1 switch (Zyxel), 1 computer (Dell), 1 printer (HP), 1 unknown device.

---

## Summary

| Issue | Title | Key Deliverable |
|-------|-------|-----------------|
| #1 | Test Infrastructure | `bun test` runs Playwright |
| #2 | Project Foundation | Web server + web shell |
| #3 | Mock Scanner | Mock mode + Scanner interface |
| #4 | Scan API | HTTP endpoints |
| #5 | Live Updates | WebSocket streaming |
| #6 | Scan UI | Complete scan page |
| #7 | Device Display | Type icons, hostname/vendor display |

**After Issue #7:** Working web UI that runs mock scans with live updates and displays device info with icons.

---

## Technical Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Runtime | Docker | Cross-platform, consistent environment |
| HTTP server | Minimal deps | No heavy frameworks |
| WebSocket | Any library | Real-time updates |
| Config | JSON | Simple |
| Assets | `web/` directory | Easy to edit |
| Frontend | Vanilla JS | No build step |
| Tests | Playwright + playwright-bdd | Real browser |

---

## Creating Issues in GitHub

Use `gh issue create` to seed the issues. Example for Issue #1:

```bash
gh issue create \
  --title "As a developer I want test infrastructure so that I can write browser tests" \
  --body "Before implementing features, we need Playwright configured to run Gherkin scenarios against the browser.

**Acceptance Criteria**

\`\`\`gherkin
Feature: Test Infrastructure

  @manual
  Scenario: Test runner works
    When I run \"bun test\"
    Then Playwright executes feature files from tests/features/
    And step definitions from tests/steps/ are used
\`\`\`

**Technical decisions:**
- Use playwright-bdd to connect Gherkin to Playwright
- Use bun as package manager"
```

After creating all issues, work from GitHub using `gh issue view {id}` and `gh issue close {id}`.
