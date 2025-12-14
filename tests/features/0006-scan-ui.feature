@serial
Feature: Scan UI

  Scenario: Scan form
    Given I am on the home page
    When I click "Scan" in the navigation
    Then I am on the scan page
    And I see a field for network range (CIDR notation, e.g. "192.168.1.0/24")
    And I see a field for core switch IP (the starting point for discovery)
    And I see a Start button

  Scenario: Start scan from UI
    Given I am on the home page
    When I click "Scan" in the navigation
    And I enter network "192.168.1.0/24" and core switch "192.168.1.1"
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
    Given I am on the home page
    When I click "Scan" in the navigation
    And I enter an invalid network
    And I click Start
    Then I see an error message

  Scenario: Cancel from UI
    Given I started a scan from the UI
    When I click Cancel
    Then the scan stops

  Scenario: Empty state
    Given I am on the home page
    When I click "Scan" in the navigation
    Then I see a message prompting to start a scan

