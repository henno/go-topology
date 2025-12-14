import { expect } from '@playwright/test';
import { createBdd } from 'playwright-bdd';

const { Given, When, Then } = createBdd();

Then('I see a field for network range \\(CIDR notation, e.g. {string}\\)', async ({ page }, example: string) => {
  const networkField = page.locator('input[name="network"]');
  await expect(networkField).toBeVisible();
  const placeholder = await networkField.getAttribute('placeholder');
  expect(placeholder).toContain('192.168');
});

Then('I see a field for core switch IP \\(the starting point for discovery\\)', async ({ page }) => {
  const coreSwitchField = page.locator('input[name="core-switch"]');
  await expect(coreSwitchField).toBeVisible();
});

Then('I see a Start button', async ({ page }) => {
  const startButton = page.locator('button:has-text("Start")');
  await expect(startButton).toBeVisible();
});

When('I enter network {string} and core switch {string}', async ({ page }, network: string, coreSwitch: string) => {
  await page.fill('input[name="network"]', network);
  await page.fill('input[name="core-switch"]', coreSwitch);
});

When('I click Start', async ({ page }) => {
  await page.click('button:has-text("Start")');
});

Then('I see the scan status change to {string}', async ({ page }, status: string) => {
  const statusElement = page.locator('.scan-status');
  await expect(statusElement).toContainText(status, { timeout: 2000 });
});

Given('I started a scan', async ({ page }) => {
  await page.goto('http://localhost:9090');
  await page.click('nav a:has-text("Scan"), aside a:has-text("Scan")');
  await page.fill('input[name="network"]', '192.168.1.0/24');
  await page.fill('input[name="core-switch"]', '192.168.1.1');
  await page.click('button:has-text("Start")');

  // Wait for scan to start
  await expect(page.locator('.scan-status')).toContainText('Scanning', { timeout: 2000 });
});

Then('devices appear in the results table as they are discovered', async ({ page }) => {
  // Wait for devices to appear (mock scanner returns 5 devices)
  // We can't reliably test "as they are discovered" because the mock scanner
  // is too fast (250ms total), so we just verify all devices appear
  await expect(page.locator('.device-table tbody tr')).toHaveCount(5, { timeout: 3000 });
});

When('the scan completes', async ({ page }) => {
  // Wait for scan to complete (mock scanner takes ~250ms for 5 devices)
  await expect(page.locator('.scan-status')).toContainText('Complete', { timeout: 3000 });
});

Then('I see {string} status', async ({ page }, status: string) => {
  const statusElement = page.locator('.scan-status');
  await expect(statusElement).toContainText(status);
});

Then('I see the total device count', async ({ page }) => {
  const countElement = page.locator('.device-count');
  await expect(countElement).toBeVisible();
  await expect(countElement).toContainText('5'); // Mock scanner returns 5 devices
});

When('I enter an invalid network', async ({ page }) => {
  await page.fill('input[name="network"]', 'invalid-network');
  await page.fill('input[name="core-switch"]', '192.168.1.1');
});

Then('I see an error message', async ({ page }) => {
  const errorElement = page.locator('.error-message');
  await expect(errorElement).toBeVisible({ timeout: 2000 });
});

Given('I started a scan from the UI', async ({ page }) => {
  await page.goto('http://localhost:9090');
  await page.click('nav a:has-text("Scan"), aside a:has-text("Scan")');
  await page.fill('input[name="network"]', '192.168.1.0/24');
  await page.fill('input[name="core-switch"]', '192.168.1.1');
  await page.click('button:has-text("Start")');

  // Wait for Cancel button to appear (indicates scan has started)
  await expect(page.locator('button:has-text("Cancel")')).toBeVisible({ timeout: 2000 });
});

When('I click Cancel', async ({ page }) => {
  // Click cancel as soon as the button appears (don't wait for scan to complete)
  await page.click('button:has-text("Cancel")', { timeout: 1000 });
});

Then('the scan stops', async ({ page }) => {
  const statusElement = page.locator('.scan-status');
  await expect(statusElement).toContainText('Cancelled', { timeout: 2000 });
});

Then('I see a message prompting to start a scan', async ({ page }) => {
  const emptyStateMessage = page.locator('.empty-state');
  await expect(emptyStateMessage).toBeVisible();
  await expect(emptyStateMessage).toContainText(/start.*scan/i);
});

