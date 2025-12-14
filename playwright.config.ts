import { defineConfig, devices } from '@playwright/test';
import { defineBddConfig } from 'playwright-bdd';

const testDir = defineBddConfig({
  features: 'tests/features/**/*.feature',
  steps: 'tests/steps/**/*.ts',
});

export default defineConfig({
  testDir,
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:9090',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'go-topology',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: 'docker compose -f docker-compose.yml -f docker-compose.dev.yml up',
    url: 'http://localhost:9090',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
});

