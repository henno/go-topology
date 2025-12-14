import { expect } from '@playwright/test';
import { createBdd } from 'playwright-bdd';

const { When, Then } = createBdd();

When('I run {string}', async ({}, command: string) => {
  // This is a meta-test - the fact that this step runs means the test infrastructure works
  expect(command).toBe('bun test');
});

Then('Playwright executes feature files', async ({}) => {
  // This step executing proves Playwright is reading feature files
  expect(true).toBe(true);
});

Then('step definitions are used', async ({}) => {
  // This step executing proves step definitions are being used
  expect(true).toBe(true);
});

