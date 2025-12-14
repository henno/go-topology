Feature: Test Infrastructure

  @manual
  Scenario: Test runner works
    When I run "bun test"
    Then Playwright executes feature files
    And step definitions are used

