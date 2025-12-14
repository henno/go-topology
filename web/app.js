// NetMap frontend application
// Minimal setup for now - will be expanded in future issues

console.log('NetMap loaded');

// Check if server is running in mock mode
async function checkMockMode() {
  try {
    const response = await fetch('/api/status');
    const data = await response.json();

    if (data.mock_mode) {
      displayMockModeIndicator();
    }
  } catch (error) {
    console.error('Failed to check mock mode:', error);
  }
}

function displayMockModeIndicator() {
  const indicator = document.createElement('div');
  indicator.className = 'mock-mode-indicator';
  indicator.textContent = 'Mock Mode';
  document.body.appendChild(indicator);
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
  checkMockMode();
});

