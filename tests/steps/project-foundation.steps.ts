import { expect } from '@playwright/test';
import { createBdd } from 'playwright-bdd';

const { Given, When, Then } = createBdd();

Given('I open the application', async ({ page }) => {
  await page.goto('/');
});

Given('I am on the home page', async ({ page }) => {
  await page.goto('/');
});

Then('I see the page title {string}', async ({ page }, title: string) => {
  await expect(page).toHaveTitle(title);
});

Then('I see a sidebar with navigation', async ({ page }) => {
  const sidebar = page.locator('aside.sidebar');
  await expect(sidebar).toBeVisible();
  const nav = sidebar.locator('nav');
  await expect(nav).toBeVisible();
});

Then('I see {string} and {string} links', async ({ page }, link1: string, link2: string) => {
  await expect(page.locator(`a:has-text("${link1}")`)).toBeVisible();
  await expect(page.locator(`a:has-text("${link2}")`)).toBeVisible();
});

When('I click {string} in the navigation', async ({ page }, linkText: string) => {
  await page.locator(`nav a:has-text("${linkText}"), aside a:has-text("${linkText}")`).click();
});

Then('I am on the scan page', async ({ page }) => {
  await expect(page).toHaveURL(/\/scan/);
});

