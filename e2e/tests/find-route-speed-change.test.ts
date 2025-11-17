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

test.describe("route time updates with cycling speed", () => {
  const testUsername = "test-user";
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

  test("time decreases after increasing cycling speed", async ({ page }) => {
    await page.goto(FRONTEND_URL);

    await expect(page).toHaveURL(`${FRONTEND_URL}/login`);

    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    const map = page.locator("#map");
    await expect(map).toBeVisible({ timeout: 10000 });
    const mapBox = await map.boundingBox();
    if (!mapBox) throw new Error("Map element not found");

    // Choose two points and remember their coordinates so we can reuse them
    const startX = mapBox.x + mapBox.width / 2;
    const startY = mapBox.y + mapBox.height / 2;
    const endX = mapBox.x + mapBox.width / 3;
    const endY = mapBox.y + mapBox.height / 3;

    await page.locator("#selectStartButton").click();
    await page.mouse.click(startX, startY);

    await page.locator("#selectEndButton").click();
    await page.mouse.click(endX, endY);

    await page.locator("#findRouteButton").click();

    const timeControl = page.locator(".time-control");
    await expect(timeControl).toBeVisible({ timeout: 10000 });
    const timeText1 = await timeControl.textContent();
    expect(timeText1).toBeTruthy();
    const minutes1Match = timeText1!.match(/(\d+)\s*min/);
    expect(minutes1Match).not.toBeNull();
    const minutes1 = parseInt(minutes1Match![1], 10);

    // Go to settings and increase cycling speed
    await page.hover("#usernameButton");
    await page.click("#settingsLink");
    await expect(page).toHaveURL(`${FRONTEND_URL}/settings`);

    // Change preferred cycling speed to a faster value
    await page.fill("#cyclingSpeed", "40");
    await page.click('button:has-text("Save Preferences")');
    await expect(page.locator("text=Settings saved")).toBeVisible({
      timeout: 5000,
    });

    // Return to map and re-run route between the same two points
    await page.click('a[href="/"]');
    await expect(page).toHaveURL(`${FRONTEND_URL}/`);

    // re-select same positions and find route again
    await page.locator("#selectStartButton").click();
    await page.mouse.click(startX, startY);
    await page.locator("#selectEndButton").click();
    await page.mouse.click(endX, endY);
    await page.locator("#findRouteButton").click();

    await expect(timeControl).toBeVisible({ timeout: 10000 });
    const timeText2 = await timeControl.textContent();
    expect(timeText2).toBeTruthy();
    const minutes2Match = timeText2!.match(/(\d+)\s*min/);
    expect(minutes2Match).not.toBeNull();
    const minutes2 = parseInt(minutes2Match![1], 10);

    // The new time should be less than the previous time after increasing speed
    expect(minutes2).toBeLessThan(minutes1);
  });
});
