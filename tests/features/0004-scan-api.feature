@serial
Feature: Scan API

  Scenario: Start scan via API
    When I POST to /api/scans with network "192.168.1.0/24" and core switch "192.168.1.1"
    Then I receive a scan ID
    And the response status is "scanning"

  Scenario: Get scan status
    Given a scan is running
    When I GET /api/scans/current
    Then I see the scan status and discovered count

  Scenario: Only one scan at a time
    Given a scan is running
    When I POST to /api/scans
    Then I receive a 409 Conflict error

  Scenario: Cancel scan
    Given a scan is running
    When I DELETE the current scan
    Then the scan is cancelled

