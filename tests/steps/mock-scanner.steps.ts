import { expect } from '@playwright/test';
import { createBdd } from 'playwright-bdd';

const { Given, Then } = createBdd();

Given('the server is running in mock mode', async () => {
  // Server is started in mock mode via docker-compose.dev.yml
  // which sets NETMAP_MOCK=true environment variable
});

Then('I see a {string} indicator', async ({ page }, text: string) => {
  await expect(page.getByText(text)).toBeVisible();
});

