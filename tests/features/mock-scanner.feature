Feature: Mock Scanner

  Scenario: Mock mode indicator
    Given the server is running in mock mode
    And I open the application
    Then I see a "Mock Mode" indicator

