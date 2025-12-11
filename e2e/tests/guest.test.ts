import { expect, test } from "@playwright/test";
import { Client } from "pg";
import * as dotenv from "dotenv";
import bcrypt from "bcrypt";

dotenv.config({ path: "../.env.test" });

function requireEnv(name: string): string {
  const value = process.env[name];
  if (!value)
    throw new Error(
      `Required environment variable ${name} is not set. Check .env.test`
    );
  return value;
}

const dbConfig = {
  user: requireEnv("POSTGRES_TEST_USER"),
  password: requireEnv("POSTGRES_TEST_PASSWORD"),
  database: requireEnv("POSTGRES_TEST_DB"),
  host: process.env.POSTGRES_TEST_HOST || "localhost",
  port: parseInt(requireEnv("POSTGRES_TEST_PORT") || "5434"),
};

const FRONTEND_URL = `http://localhost:${process.env.FRONTEND_PORT || 3001}`;

test.describe("guest flows", () => {
  const testUsername = "guest-flow-user";
  const testPassword = "guest-flow-password";
  const wayIdWithReviews = 20948648;

  let client: Client;

  test.beforeAll(async () => {
    client = new Client(dbConfig);
    await client.connect();
  });

  test.afterAll(async () => {
    await client.end();
  });

  test.beforeEach(async () => {
    await client.query("TRUNCATE TABLE reviews RESTART IDENTITY CASCADE");
    await client.query("TRUNCATE TABLE users RESTART IDENTITY CASCADE");
  });

  test.afterEach(async () => {
    await client.query("TRUNCATE TABLE reviews RESTART IDENTITY CASCADE");
    await client.query("TRUNCATE TABLE users RESTART IDENTITY CASCADE");
  });

  test("guest can access map and find a route", async ({ page }) => {
    await page.goto(`${FRONTEND_URL}/map`);

    const map = page.locator("#map");
    await expect(map).toBeVisible({ timeout: 10000 });
    await page.waitForLoadState("networkidle");

    const mapBox = await map.boundingBox();
    if (!mapBox) throw new Error("Map element not found");

    const markers = page.locator(".leaflet-marker-icon");

    const startButton = page.locator("#selectStartButton");
    await startButton.click();
    await page.waitForTimeout(300);
    await page.mouse.click(
      mapBox.x + mapBox.width / 2,
      mapBox.y + mapBox.height / 2
    );
    await expect(markers).toHaveCount(1, { timeout: 5000 });

    const endButton = page.locator("#selectEndButton");
    await endButton.click();
    await page.waitForTimeout(300);
    await page.mouse.click(
      mapBox.x + mapBox.width / 3,
      mapBox.y + mapBox.height / 3
    );
    await expect(markers).toHaveCount(2, { timeout: 5000 });

    await page.locator("#findRouteButton").click();

    const distanceControl = page.locator(".distance-control");
    await expect(distanceControl).toBeVisible({ timeout: 10000 });
  });

  test("guest can navigate to login and authenticate", async ({ page }) => {
    const hashedPassword = await bcrypt.hash(
      testPassword,
      bcrypt.genSaltSync()
    );
    await client.query(
      "INSERT INTO users (username, password) VALUES ($1, $2)",
      [testUsername, hashedPassword]
    );

    await page.goto(`${FRONTEND_URL}/map`);

    const loginLink = page.locator("#loginLink");
    await expect(loginLink).toBeVisible({ timeout: 5000 });
    await loginLink.click();

    await expect(page).toHaveURL(`${FRONTEND_URL}/login`);

    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    await expect(page).toHaveURL(`${FRONTEND_URL}/map`);

    // Wait for map and menu to finish hydrating before asserting auth UI
    await expect(page.locator("#map")).toBeVisible({ timeout: 15000 });
    await page.waitForSelector("#usernameButton", {
      state: "visible",
      timeout: 15000,
    });
  });

  test("guest can view reviews but cannot add one", async ({ page }) => {
    const seededPassword = "seed-password";
    const seededUsername = "seed-user";
    const hashedPassword = await bcrypt.hash(
      seededPassword,
      bcrypt.genSaltSync()
    );
    const res = await client.query(
      "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
      [seededUsername, hashedPassword]
    );
    const seededUserId = res.rows[0].id;

    const existingComment = `guest-visible review ${Date.now()}`;
    await client.query(
      "INSERT INTO reviews (way_id, user_id, rating, comment, created_at) VALUES ($1, $2, $3, $4, now())",
      [wayIdWithReviews, seededUserId, 7, existingComment]
    );

    await page.goto(`${FRONTEND_URL}/map?way=${wayIdWithReviews}`);
    await page.waitForLoadState("networkidle");
    await page.waitForSelector("#sidebar-content", {
      state: "attached",
      timeout: 20000,
    });

    await expect(page.locator(`text=${existingComment}`)).toBeVisible({
      timeout: 5000,
    });
    await expect(page.locator("#addReviewButton")).toHaveCount(0);
    await expect(page.locator("text=Log in to add a review.")).toBeVisible({
      timeout: 5000,
    });
  });
});
