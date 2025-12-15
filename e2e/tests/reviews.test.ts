import { expect, test } from "@playwright/test";
import { Client } from "pg";
import * as dotenv from "dotenv";
import bcrypt from "bcrypt";

dotenv.config({ path: "../.env.test" });

function requireEnv(name: string): string {
  const v = process.env[name];
  if (!v)
    throw new Error(
      `Required environment variable ${name} is not set. Check .env.test`
    );
  return v;
}

const dbConfig = {
  user: requireEnv("POSTGRES_TEST_USER"),
  password: requireEnv("POSTGRES_TEST_PASSWORD"),
  database: requireEnv("POSTGRES_TEST_DB"),
  host: process.env.POSTGRES_TEST_HOST || "localhost",
  port: parseInt(requireEnv("POSTGRES_TEST_PORT") || "5434"),
};

const FRONTEND_URL = `http://localhost:${process.env.FRONTEND_PORT || 3001}`;
const API_URL = process.env.API_URL || "http://localhost:8080";

test.describe("reviews e2e", () => {
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
    await client.query("TRUNCATE TABLE users RESTART IDENTITY CASCADE");
    await client.query("TRUNCATE TABLE reviews RESTART IDENTITY CASCADE");
  });

  test.afterEach(async () => {
    await client.query("TRUNCATE TABLE users RESTART IDENTITY CASCADE");
    await client.query("TRUNCATE TABLE reviews RESTART IDENTITY CASCADE");
  });

  // View a way, create new review with just one way, and it appears in the sidebar
  test("view a way, create new review (single way), and it appears in the sidebar", async ({
    page,
  }) => {
    const hashedPassword = await bcrypt.hash(
      testPassword,
      bcrypt.genSaltSync()
    );
    await client.query(
      "INSERT INTO users (username, password) VALUES ($1, $2)",
      [testUsername, hashedPassword]
    );

    await page.goto(`${FRONTEND_URL}/login`);
    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    await page.goto(`${FRONTEND_URL}/map?way=20948648`);
    await page.waitForLoadState("networkidle");
    await page.waitForSelector("#sidebar-content", {
      state: "attached",
      timeout: 20000,
    });

    const addButton = page.locator("#addReviewButton");
    await expect(addButton).toBeVisible({ timeout: 5000 });
    await addButton.click();
    // Inline form replaces reviews; header changes to "Add Review"
    await expect(page.locator("text=Add Review")).toBeVisible({
      timeout: 2000,
    });

    const uniqueComment = `e2e test comment ${Date.now()}`;
    await page.fill("#rating", "9");
    await page.fill("#comment", uniqueComment);
    await page.locator("#submitButton").click();

    // After submit, inline form closes and header returns to "Reviews"
    await expect(page.locator("text=Reviews")).toBeVisible({
      timeout: 5000,
    });
    await expect(page.locator(`text=${uniqueComment}`)).toBeVisible({
      timeout: 5000,
    });
    await expect(page.locator(`text=User: ${testUsername}`)).toBeVisible({
      timeout: 5000,
    });
  });

  // View a way, create new review with multiple ways, and it appears in the sidebar for each of the ways
  test("create new review with multiple ways and it appears for each way", async ({
    page,
  }) => {
    // Seed user
    const hashedPassword = await bcrypt.hash(
      testPassword,
      bcrypt.genSaltSync()
    );
    await client.query(
      "INSERT INTO users (username, password) VALUES ($1, $2)",
      [testUsername, hashedPassword]
    );

    // Login
    await page.goto(`${FRONTEND_URL}/login`);
    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    const baseWayId = 20948648;
    await page.goto(`${FRONTEND_URL}/map?way=${baseWayId}`);
    await page.waitForLoadState("networkidle");
    await page.waitForSelector("#sidebar-content", {
      state: "attached",
      timeout: 20000,
    });

    // Enter Add Review mode
    const addButton = page.locator("#addReviewButton");
    await expect(addButton).toBeVisible({ timeout: 5000 });
    await addButton.click();
    await expect(page.locator("text=Add Review")).toBeVisible({
      timeout: 2000,
    });

    // Click an adjacent orange way (skip the original selected blue way)
    // Wait for at least one visible adjacent path (non-zero bounding box).
    await page.waitForFunction(
      () =>
        Array.from(
          document.querySelectorAll("svg path.adjacent-way.leaflet-interactive")
        ).some((p) => {
          const r = p.getBoundingClientRect();
          return r.width > 0 && r.height > 0;
        }),
      { timeout: 15000 }
    );
    await page.evaluate(() => {
      const paths = Array.from(
        document.querySelectorAll("svg path.adjacent-way.leaflet-interactive")
      );
      const target = paths.find((p) => {
        const r = p.getBoundingClientRect();
        return r.width > 0 && r.height > 0;
      });
      if (!target) throw new Error("No visible adjacent-way paths found");
      target.dispatchEvent(
        new MouseEvent("click", {
          bubbles: true,
          cancelable: true,
          view: window,
        })
      );
    });

    // Confirm that another way chip was added in the included ways list
    const removeButtons = page.locator(
      '#includedWays button[aria-label^="Remove way "]'
    );
    await expect(removeButtons).toHaveCount(1, { timeout: 5000 });
    const ariaLabel = await removeButtons.first().getAttribute("aria-label");
    const match = ariaLabel?.match(/Remove way (\d+)/);
    const additionalWayId = match ? parseInt(match[1], 10) : null;
    expect(additionalWayId).toBeTruthy();

    // Fill and submit review
    const uniqueComment = `e2e multi-way comment ${Date.now()}`;
    await page.fill("#rating", "8");
    await page.fill("#comment", uniqueComment);
    await page.locator("#submitButton").click();

    // Verify appears on base way page
    await expect(page.locator("text=Reviews")).toBeVisible({ timeout: 5000 });
    await expect(page.locator(`text=${uniqueComment}`)).toBeVisible({
      timeout: 5000,
    });
    await expect(page.locator(`text=User: ${testUsername}`)).toBeVisible({
      timeout: 5000,
    });

    // Navigate to the additional way and verify the same review appears
    if (additionalWayId) {
      await page.goto(`${FRONTEND_URL}/map?way=${additionalWayId}`);
      await page.waitForLoadState("networkidle");
      await page.waitForSelector("#sidebar-content", {
        state: "attached",
        timeout: 20000,
      });
      await expect(page.locator(`text=${uniqueComment}`)).toBeVisible({
        timeout: 5000,
      });
      await expect(page.locator(`text=User: ${testUsername}`)).toBeVisible({
        timeout: 5000,
      });
    }
  });

  // View the sidebar for the way and assert that the existing review is shown
  test("view a way with pre-existing reviews and assert they appear", async ({
    page,
  }) => {
    const seededUsername = `seed-user-${Date.now()}`;
    const seededPassword = "seed-password";
    const hashedPassword = await bcrypt.hash(
      seededPassword,
      bcrypt.genSaltSync()
    );
    const res = await client.query(
      "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
      [seededUsername, hashedPassword]
    );
    const seededUserId = res.rows[0].id;

    const wayId = 20948648;

    const existingComment = `preexisting review ${Date.now()}`;
    // Insert review (without way_id) and link it via review_ways
    const insertRes = await client.query(
      "INSERT INTO reviews (user_id, rating, comment, created_at) VALUES ($1, $2, $3, now()) RETURNING id",
      [seededUserId, 7, existingComment]
    );
    const reviewId = insertRes.rows[0].id as number;
    await client.query(
      "INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)",
      [reviewId, wayId]
    );

    await page.goto(`${FRONTEND_URL}/login`);
    await page.fill('input[name="username"]', seededUsername);
    await page.fill('input[name="password"]', seededPassword);
    await page.locator("#submitButton").click();

    await page.goto(`${FRONTEND_URL}/map?way=${wayId}`);
    await page.waitForLoadState("networkidle");
    await page.waitForSelector("#sidebar-content", {
      state: "attached",
      timeout: 20000,
    });

    await expect(page.locator(`text=${existingComment}`)).toBeVisible({
      timeout: 5000,
    });
    await expect(page.locator(`text=User: ${seededUsername}`)).toBeVisible({
      timeout: 5000,
    });
  });

  // User can delete their own review
  test("user can delete their own review and it disappears from sidebar", async ({
    page,
  }) => {
    const hashedPassword = await bcrypt.hash(
      testPassword,
      bcrypt.genSaltSync()
    );
    await client.query(
      "INSERT INTO users (username, password) VALUES ($1, $2)",
      [testUsername, hashedPassword]
    );

    await page.goto(`${FRONTEND_URL}/login`);
    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    await page.goto(`${FRONTEND_URL}/map?way=20948648`);
    await page.waitForLoadState("networkidle");
    await page.waitForSelector("#sidebar-content", {
      state: "attached",
      timeout: 20000,
    });

    // Add a review
    const addButton = page.locator("#addReviewButton");
    await expect(addButton).toBeVisible({ timeout: 5000 });
    await addButton.click();
    await expect(page.locator("text=Add Review")).toBeVisible({
      timeout: 2000,
    });

    const uniqueComment = `review to delete ${Date.now()}`;
    await page.fill("#rating", "6");
    await page.fill("#comment", uniqueComment);
    await page.locator("#submitButton").click();

    // Verify review appears
    await expect(page.locator("text=Reviews")).toBeVisible({
      timeout: 5000,
    });
    await expect(page.locator(`text=${uniqueComment}`)).toBeVisible({
      timeout: 5000,
    });

    // Click delete button
    const deleteButton = page.locator('button:has-text("Delete")');
    await expect(deleteButton).toBeVisible({ timeout: 5000 });
    await deleteButton.click();

    // Verify review is removed from sidebar
    await expect(page.locator(`text=${uniqueComment}`)).not.toBeVisible({
      timeout: 5000,
    });

    // Verify add review button appears again (since user no longer has review)
    await expect(page.locator("#addReviewButton")).toBeVisible({
      timeout: 5000,
    });
  });

  // User cannot delete other users' reviews (button not shown)
  test("user cannot delete other users' reviews", async ({ page }) => {
    const otherUsername = `other-user-${Date.now()}`;
    const otherPassword = "other-password";
    const otherHashed = await bcrypt.hash(otherPassword, bcrypt.genSaltSync());
    const otherRes = await client.query(
      "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
      [otherUsername, otherHashed]
    );
    const otherUserId = otherRes.rows[0].id;

    const wayId = 20948648;
    const otherComment = `other user review ${Date.now()}`;
    const insertRes = await client.query(
      "INSERT INTO reviews (user_id, rating, comment, created_at) VALUES ($1, $2, $3, now()) RETURNING id",
      [otherUserId, 8, otherComment]
    );
    const reviewId = insertRes.rows[0].id as number;
    await client.query(
      "INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)",
      [reviewId, wayId]
    );

    // Login as testUsername (different user)
    const hashedPassword = await bcrypt.hash(
      testPassword,
      bcrypt.genSaltSync()
    );
    await client.query(
      "INSERT INTO users (username, password) VALUES ($1, $2)",
      [testUsername, hashedPassword]
    );

    await page.goto(`${FRONTEND_URL}/login`);
    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    await page.goto(`${FRONTEND_URL}/map?way=${wayId}`);
    await page.waitForLoadState("networkidle");
    await page.waitForSelector("#sidebar-content", {
      state: "attached",
      timeout: 20000,
    });

    // Verify other user's review appears
    await expect(page.locator(`text=${otherComment}`)).toBeVisible({
      timeout: 5000,
    });

    // Verify delete button is NOT visible for other user's review
    const deleteButtons = page.locator('button:has-text("Delete")');
    await expect(deleteButtons).toHaveCount(0);
  });
});
