import { expect, test } from "@playwright/test";
import { Client } from "pg";
import * as dotenv from "dotenv";

dotenv.config({ path: "../.env.test" });

const dbConfig = {
  user: process.env.POSTGRES_TEST_USER,
  password: process.env.POSTGRES_TEST_PASSWORD,
  database: process.env.POSTGRES_TEST_DB,
  host: "localhost",
  port: parseInt(process.env.POSTGRES_TEST_PORT || "5434"),
};

const FRONTEND_URL = `http://localhost:${process.env.FRONTEND_PORT || 3001}`;

test.describe("sign up", () => {
  const testUsername = `test-user-${Date.now()}`;
  const testPassword = "test-password";

  let client: Client;
  test.beforeAll(async () => {
    client = new Client(dbConfig);
    await client.connect();
  });

  test.afterAll(async () => {
    await client.end();
  });

  test.beforeEach(async () => {
    await client.query("TRUNCATE TABLE users RESTART IDENTITY CASCADE");
  });

  test.afterEach(async () => {
    await client.query("TRUNCATE TABLE users RESTART IDENTITY CASCADE");
  });

  test("user can sign up, log in, look up a route, then log out", async ({
    page,
  }) => {
    await page.goto(`${FRONTEND_URL}/login`);

    const signupLink = page.locator("#signupLink");
    await signupLink.click();

    await expect(page).toHaveURL(`${FRONTEND_URL}/signup`);

    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    await expect(page).toHaveURL(`${FRONTEND_URL}/login`);

    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    await expect(page).toHaveURL(`${FRONTEND_URL}/map`);

    // Wait for map to load and ways data to be streamed
    const map = page.locator("#map");
    await expect(map).toBeVisible({ timeout: 10000 });
    await page.waitForLoadState("networkidle");

    const mapBox = await map.boundingBox();
    if (!mapBox) throw new Error("Map element not found");

    // Click start button and wait for active state
    await page.locator("#selectStartButton").click();
    await page.waitForTimeout(300);

    await page.mouse.click(
      mapBox.x + mapBox.width / 2,
      mapBox.y + mapBox.height / 2
    );

    // Click end button and wait for active state
    await page.locator("#selectEndButton").click();
    await page.waitForTimeout(300);

    await page.mouse.click(
      mapBox.x + mapBox.width / 3,
      mapBox.y + mapBox.height / 3
    );

    await page.locator("#findRouteButton").click();

    // Check for route controls instead of paths (more reliable)
    const distanceControl = page.locator(".distance-control");
    await expect(distanceControl).toBeVisible({ timeout: 10000 });

    // Verify logout
    await page.locator("#usernameButton").hover();
    await page.locator("#logoutButton").click();

    await expect(page).toHaveURL(`${FRONTEND_URL}/login`);
  });
});
