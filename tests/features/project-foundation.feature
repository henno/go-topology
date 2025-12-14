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

