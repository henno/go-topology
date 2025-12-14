import { expect } from '@playwright/test';
import { createBdd } from 'playwright-bdd';

const { Given, When, Then } = createBdd();

let scanResponse: any;
let scanId: string;
let statusResponse: any;

When('I POST to \\/api\\/scans with network {string} and core switch {string}', async ({ request }, network: string, coreSwitch: string) => {
  scanResponse = await request.post('http://localhost:9090/api/scans', {
    data: {
      network: network,
      core_switch: coreSwitch
    }
  });

  if (!scanResponse.ok()) {
    const text = await scanResponse.text();
    console.log(`POST /api/scans failed with status ${scanResponse.status()}: ${text}`);
  } else {
    const data = await scanResponse.json();
    scanId = data.id;
  }
});

Then('I receive a scan ID', async () => {
  expect(scanResponse.ok()).toBeTruthy();
  const data = await scanResponse.json();
  expect(data.id).toBeTruthy();
  expect(typeof data.id).toBe('string');
});

Then('the response status is {string}', async ({}, expectedStatus: string) => {
  const data = await scanResponse.json();
  expect(data.status).toBe(expectedStatus);
});

Given('a scan is running', async ({ request }) => {
  // Wait for any existing scan to complete (max 2 seconds)
  for (let i = 0; i < 20; i++) {
    const currentScanResponse = await request.get('http://localhost:9090/api/scans/current');
    if (currentScanResponse.status() === 404) {
      // No scan running, we can start a new one
      break;
    }
    const currentScan = await currentScanResponse.json();
    if (currentScan.status !== 'scanning') {
      // Scan completed, we can start a new one
      break;
    }
    // Wait 100ms before checking again
    await new Promise(resolve => setTimeout(resolve, 100));
  }

  // Start a scan
  scanResponse = await request.post('http://localhost:9090/api/scans', {
    data: {
      network: '192.168.1.0/24',
      core_switch: '192.168.1.1'
    }
  });

  if (!scanResponse.ok()) {
    const text = await scanResponse.text();
    throw new Error(`Failed to start scan: ${scanResponse.status()} - ${text}`);
  }

  const data = await scanResponse.json();
  scanId = data.id;
});

When('I GET \\/api\\/scans\\/current', async ({ request }) => {
  statusResponse = await request.get('http://localhost:9090/api/scans/current');
});

Then('I see the scan status and discovered count', async () => {
  expect(statusResponse.ok()).toBeTruthy();
  const data = await statusResponse.json();
  expect(data.status).toBeTruthy();
  expect(typeof data.discovered_count).toBe('number');
});

When('I POST to \\/api\\/scans', async ({ request }) => {
  scanResponse = await request.post('http://localhost:9090/api/scans', {
    data: {
      network: '192.168.1.0/24',
      core_switch: '192.168.1.1'
    }
  });
});

Then('I receive a 409 Conflict error', async () => {
  expect(scanResponse.status()).toBe(409);
});

When('I DELETE the current scan', async ({ request }) => {
  statusResponse = await request.delete(`http://localhost:9090/api/scans/${scanId}`);
});

Then('the scan is cancelled', async ({ request }) => {
  expect(statusResponse.ok()).toBeTruthy();

  // Wait for the scan status to change to cancelled (max 500ms)
  await expect(async () => {
    const currentScanResponse = await request.get('http://localhost:9090/api/scans/current');
    expect(currentScanResponse.ok()).toBeTruthy();
    const data = await currentScanResponse.json();
    expect(data.status).toBe('cancelled');
  }).toPass({ timeout: 500 });
});

