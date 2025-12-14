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

// Scan functionality
let currentScanId = null;
let pollInterval = null;

function initializeScanPage() {
  const startBtn = document.getElementById('start-btn');
  const cancelBtn = document.getElementById('cancel-btn');

  if (!startBtn) return; // Not on scan page

  startBtn.addEventListener('click', startScan);
  cancelBtn.addEventListener('click', cancelScan);
}

async function startScan() {
  const networkInput = document.getElementById('network');
  const coreSwitchInput = document.getElementById('core-switch');
  const errorMessage = document.querySelector('.error-message');
  const startBtn = document.getElementById('start-btn');
  const cancelBtn = document.getElementById('cancel-btn');

  const network = networkInput.value.trim();
  const coreSwitch = coreSwitchInput.value.trim();

  // Hide error message
  errorMessage.style.display = 'none';
  errorMessage.textContent = '';

  // Basic validation
  if (!network || !coreSwitch) {
    showError('Please enter both network range and core switch IP');
    return;
  }

  // Validate CIDR notation (basic check)
  if (!network.includes('/')) {
    showError('Network must be in CIDR notation (e.g., 192.168.1.0/24)');
    return;
  }

  try {
    const response = await fetch('/api/scans', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        network: network,
        core_switch: coreSwitch,
      }),
    });

    if (!response.ok) {
      const error = await response.json();
      showError(error.error || 'Failed to start scan');
      return;
    }

    const scan = await response.json();
    currentScanId = scan.id;

    // Update UI
    startBtn.style.display = 'none';
    cancelBtn.style.display = 'inline-block';
    document.querySelector('.empty-state').style.display = 'none';
    document.querySelector('.scan-status').style.display = 'block';
    document.querySelector('.device-count').style.display = 'block';
    document.querySelector('.device-table-container').style.display = 'block';

    // Clear previous results
    document.getElementById('device-table-body').innerHTML = '';
    document.getElementById('device-count-text').textContent = '0';

    // Start polling for updates
    startPolling();
  } catch (error) {
    showError('Failed to start scan: ' + error.message);
  }
}

async function cancelScan() {
  if (!currentScanId) return;

  try {
    const response = await fetch(`/api/scans/${currentScanId}`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      console.error('Failed to cancel scan');
    }
  } catch (error) {
    console.error('Failed to cancel scan:', error);
  }
}

function startPolling() {
  if (pollInterval) {
    clearInterval(pollInterval);
  }

  pollInterval = setInterval(async () => {
    try {
      const response = await fetch('/api/scans/current');

      if (!response.ok) {
        stopPolling();
        return;
      }

      const scan = await response.json();
      updateScanUI(scan);

      // Stop polling if scan is complete
      if (scan.status !== 'scanning') {
        stopPolling();
      }
    } catch (error) {
      console.error('Failed to poll scan status:', error);
    }
  }, 200); // Poll every 200ms
}

function stopPolling() {
  if (pollInterval) {
    clearInterval(pollInterval);
    pollInterval = null;
  }

  // Reset buttons
  document.getElementById('start-btn').style.display = 'inline-block';
  document.getElementById('cancel-btn').style.display = 'none';
}

function updateScanUI(scan) {
  // Update status
  const statusText = document.getElementById('status-text');
  statusText.textContent = scan.status.charAt(0).toUpperCase() + scan.status.slice(1);

  // Update device count
  const deviceCount = scan.devices ? scan.devices.length : 0;
  document.getElementById('device-count-text').textContent = deviceCount;

  // Update device table
  if (scan.devices && scan.devices.length > 0) {
    const tbody = document.getElementById('device-table-body');
    tbody.innerHTML = '';

    scan.devices.forEach(device => {
      const row = document.createElement('tr');
      row.innerHTML = `
        <td>${device.ip_address || ''}</td>
        <td>${device.hostname || ''}</td>
        <td>${device.type || ''}</td>
        <td>${device.vendor || ''}</td>
      `;
      tbody.appendChild(row);
    });
  }
}

function showError(message) {
  const errorMessage = document.querySelector('.error-message');
  errorMessage.textContent = message;
  errorMessage.style.display = 'block';
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
  checkMockMode();
  initializeScanPage();
});

