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

test.describe("route finding with a test user", () => {
  const testUsername = "test-user";
  const testPassword = "test-password";

  let client: Client;
  test.beforeAll(async () => {
    client = new Client(dbConfig);
    await client.connect();

    const result = await client.query(
      "SELECT current_database(), current_user"
    );
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

  test("user can log in, then select a start and end point and find a route between those two points", async ({
    page,
  }) => {
    await page.goto(`${FRONTEND_URL}/login`);

    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    // Wait for map to load and ways data to be streamed
    const map = page.locator("#map");
    await expect(map).toBeVisible({ timeout: 10000 });
    await page.waitForLoadState("networkidle");

    const mapBox = await map.boundingBox();
    if (!mapBox) throw new Error("Map element not found");

    // Click start button and ensure it's active before clicking map
    const startButton = page.locator("#selectStartButton");
    await startButton.click();
    await page.waitForTimeout(300); // Allow button state to update

    await page.mouse.click(
      mapBox.x + mapBox.width / 2,
      mapBox.y + mapBox.height / 2
    );

    // Wait for marker to appear
    const markers = page.locator(".leaflet-marker-icon");
    await expect(markers).toHaveCount(1, { timeout: 5000 });

    // Click end button and wait for it to be active
    const endButton = page.locator("#selectEndButton");
    await endButton.click();
    await page.waitForTimeout(300);

    await page.mouse.click(
      mapBox.x + mapBox.width / 3,
      mapBox.y + mapBox.height / 3
    );

    // Wait for second marker
    await expect(markers).toHaveCount(2, { timeout: 5000 });

    await page.locator("#findRouteButton").click();

    // Should display 3 alternative routes
    const routePaths = page.locator("path.leaflet-interactive");
    await expect(routePaths.first()).toBeVisible({ timeout: 10000 });
    const count = await routePaths.count();
    expect(count).toBeGreaterThanOrEqual(3);

    // Verify colors: one blue route (fastest) and at least two green routes
    let blueRoutes = 0;
    let greenRoutes = 0;
    let bluePathIndex = -1;
    for (let i = 0; i < count; i++) {
      const path = routePaths.nth(i);
      const stroke = await path.getAttribute("stroke");
      if (stroke === "blue") {
        blueRoutes++;
        bluePathIndex = i;
      }
      if (stroke === "green") greenRoutes++;
    }
    expect(blueRoutes).toBe(1);
    expect(greenRoutes).toBeGreaterThanOrEqual(2);

    const distanceControl = page.locator(".distance-control");
    await expect(distanceControl).toBeVisible({ timeout: 5000 });
    const distanceText = await distanceControl.textContent();
    expect(distanceText).toMatch(/km$/); // ends with "km"
    const timeControl = page.locator(".time-control");
    await expect(timeControl).toBeVisible({ timeout: 5000 });
    const timeText = await timeControl.textContent();
    expect(timeText).toMatch(/min$/); // ends with "min"
  });

  test("user can log in, then select a start and end point, find a route between those two points, then reset the map", async ({
    page,
  }) => {
    await page.goto(`${FRONTEND_URL}/login`);

    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    // Wait for map to load and ways data to be streamed
    const map = page.locator("#map");
    await expect(map).toBeVisible({ timeout: 10000 });
    await page.waitForLoadState("networkidle");

    const mapBox = await map.boundingBox();
    if (!mapBox) throw new Error("Map element not found");

    const markers = page.locator(".leaflet-marker-icon");

    // Click start button and wait for active state
    const startButton = page.locator("#selectStartButton");
    await startButton.click();
    await page.waitForTimeout(300);

    await page.mouse.click(
      mapBox.x + mapBox.width / 2,
      mapBox.y + mapBox.height / 2
    );

    await expect(markers).toHaveCount(1, { timeout: 5000 });

    // Click end button and wait for active state
    const endButton = page.locator("#selectEndButton");
    await endButton.click();
    await page.waitForTimeout(300);

    await page.mouse.click(
      mapBox.x + mapBox.width / 3,
      mapBox.y + mapBox.height / 3
    );

    await expect(markers).toHaveCount(2, { timeout: 5000 });

    await page.locator("#findRouteButton").click();

    // Should display 3 alternative routes
    const routePaths = page.locator("path.leaflet-interactive");
    await expect(routePaths.first()).toBeVisible({ timeout: 10000 });
    const count = await routePaths.count();
    expect(count).toBeGreaterThanOrEqual(3);

    // Verify colors: one blue route (fastest) and at least two green routes
    let blueRoutes = 0;
    let greenRoutes = 0;
    let greenPathIndex = -1;
    for (let i = 0; i < count; i++) {
      const path = routePaths.nth(i);
      const stroke = await path.getAttribute("stroke");
      if (stroke === "blue") blueRoutes++;
      if (stroke === "green") {
        greenRoutes++;
        if (greenPathIndex === -1) greenPathIndex = i;
      }
    }
    expect(blueRoutes).toBe(1);
    expect(greenRoutes).toBeGreaterThanOrEqual(2);

    await page.locator("#resetButton").click();

    // Wait for route controls to be removed (distance and time controls disappear when route is removed)
    const distanceControl = page.locator(".distance-control");
    const timeControl = page.locator(".time-control");
    await expect(distanceControl).toHaveCount(0, { timeout: 5000 });
    await expect(timeControl).toHaveCount(0, { timeout: 5000 });

    // Wait for markers to be removed
    await expect(markers).toHaveCount(0, { timeout: 5000 });
  });
});
