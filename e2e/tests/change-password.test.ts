import { expect, test } from "@playwright/test";
import { Client } from "pg";
import * as dotenv from "dotenv";
import bcrypt from "bcrypt";

dotenv.config({ path: "../.env.test" });

const dbConfig = {
  user: process.env.POSTGRES_TEST_USER,
  password: process.env.POSTGRES_TEST_PASSWORD,
  database: process.env.POSTGRES_TEST_DB,
  host: "localhost",
  port: parseInt(process.env.POSTGRES_TEST_PORT || "5434"),
};

const FRONTEND_URL = `http://localhost:${process.env.FRONTEND_PORT || 3001}`;

test.describe("change password", () => {
  const testUsername = "test-user";
  const testPassword = "test-password";
  const testNewPassword = "new-test-password";

  let client: Client;
  test.beforeAll(async () => {
    client = new Client(dbConfig);
    await client.connect();
  });

  test.afterAll(async () => {
    await client.end();
  });

  test.beforeEach(async () => {
    // Clear users table and insert test user
    await client.query("TRUNCATE TABLE users RESTART IDENTITY CASCADE");
    const hashedPassword = await bcrypt.hash(
      testPassword,
      bcrypt.genSaltSync()
    );
    await client.query(
      "INSERT INTO users (username, password) VALUES ($1, $2)",
      [testUsername, hashedPassword]
    );
  });

  test.afterEach(async () => {
    // Clear users after test
    await client.query("TRUNCATE TABLE users RESTART IDENTITY CASCADE");
  });

  test("user can log in, change their password to a new password, log out, and log in with their new password", async ({
    page,
  }) => {
    await page.goto(FRONTEND_URL);

    await expect(page).toHaveURL(`${FRONTEND_URL}/login`);

    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    const map = page.locator("#map");
    await expect(map).toBeVisible({ timeout: 10000 });

    const mapBox = await map.boundingBox();
    if (!mapBox) throw new Error("Map element not found");

    // Navigate to settings via the username menu
    await page.hover("#usernameButton");
    await page.click("#settingsLink");

    // Ensure we're on the settings page
    await expect(page).toHaveURL(`${FRONTEND_URL}/settings`);

    // Fill out change password form and submit
    await page.fill('input[name="currentPassword"]', testPassword);
    await page.fill('input[name="newPassword"]', testNewPassword);
    await page.fill('input[name="confirmPassword"]', testNewPassword);
    await page.click("#submitButton");

    // Expect success message to appear
    await expect(
      page.locator("text=Password changed successfully")
    ).toBeVisible({ timeout: 5000 });

    // Return to the map page before logging out
    await page.click('a[href="/"]');
    await expect(page).toHaveURL(`${FRONTEND_URL}/`);

    const mapBeforeLogout = page.locator("#map");
    await expect(mapBeforeLogout).toBeVisible({ timeout: 10000 });

    // Log out and ensure we're back on the login page
    await page.hover("#usernameButton");
    await page.click("#logoutButton");
    await expect(page).toHaveURL(`${FRONTEND_URL}/login`);

    // Log in with the new password
    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testNewPassword);
    await page.locator("#submitButton").click();

    const mapAfter = page.locator("#map");
    await expect(mapAfter).toBeVisible({ timeout: 10000 });
  });

  test("changing password fails when new passwords do not match", async ({
    page,
  }) => {
    await page.goto(FRONTEND_URL);

    await expect(page).toHaveURL(`${FRONTEND_URL}/login`);

    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    const map = page.locator("#map");
    await expect(map).toBeVisible({ timeout: 10000 });

    // Navigate to settings
    await page.hover("#usernameButton");
    await page.click("#settingsLink");
    await expect(page).toHaveURL(`${FRONTEND_URL}/settings`);

    // Enter non-matching new passwords
    await page.fill('input[name="currentPassword"]', testPassword);
    await page.fill('input[name="newPassword"]', "first-new-password");
    await page.fill('input[name="confirmPassword"]', "different-new-password");
    await page.click("#submitButton");

    // Expect client-side validation error to appear
    await expect(page.locator("text=New passwords do not match")).toBeVisible({
      timeout: 5000,
    });

    // Stay on settings page
    await expect(page).toHaveURL(`${FRONTEND_URL}/settings`);
  });
});
