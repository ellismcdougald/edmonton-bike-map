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

async function getToken(page) {
  return page.evaluate(
    () =>
      new Promise<string>((resolve) => {
        const interval = setInterval(() => {
          const t = localStorage.getItem("token");
          if (t) {
            clearInterval(interval);
            resolve(t);
          }
        }, 50);
      })
  );
}

async function fetchWayId(page, token) {
  const apiRes = await page.request.get(`${API_URL}/api/all-ways`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!apiRes.ok())
    throw new Error(`Failed to fetch ways from API: ${apiRes.status()}`);
  const geojson = await apiRes.json();
  if (!geojson?.features?.length)
    throw new Error("No ways available from API to run reviews test against");

  const feature = geojson.features.find((f: any) => {
    const candidate = f.id ?? f.properties?.id;
    return candidate != null && !Number.isNaN(Number(candidate));
  });
  if (!feature)
    throw new Error(
      "No feature with numeric id found in /api/all-ways response"
    );
  return Number(feature.id ?? feature.properties.id);
}

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

  test("view a way, create new review, and it appears in the sidebar", async ({
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

    await page.goto(FRONTEND_URL);
    await expect(page).toHaveURL(`${FRONTEND_URL}/login`);
    await page.fill('input[name="username"]', testUsername);
    await page.fill('input[name="password"]', testPassword);
    await page.locator("#submitButton").click();

    const token = await getToken(page);
    const wayId = await fetchWayId(page, token);

    await page.goto(`${FRONTEND_URL}/?way=${wayId}`);
    await page.waitForLoadState("networkidle");
    await page.waitForSelector("#sidebar-content", {
      state: "attached",
      timeout: 20000,
    });

    const addButton = page.locator("#addReviewButton");
    await expect(addButton).toBeVisible({ timeout: 5000 });
    await addButton.click();
    await expect(page.locator("#addReviewPopup")).toBeVisible({
      timeout: 2000,
    });

    const uniqueComment = `e2e test comment ${Date.now()}`;
    await page.fill("#rating", "9");
    await page.fill("#comment", uniqueComment);
    await page.locator("#submitButton").click();

    await expect(page.locator("#addReviewPopup")).toHaveCount(0, {
      timeout: 5000,
    });
    await expect(page.locator(`text=${uniqueComment}`)).toBeVisible({
      timeout: 5000,
    });
    await expect(page.locator(`text=User: ${testUsername}`)).toBeVisible({
      timeout: 5000,
    });
  });

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

    await page.goto(FRONTEND_URL);
    await expect(page).toHaveURL(`${FRONTEND_URL}/login`);
    await page.fill('input[name="username"]', seededUsername);
    await page.fill('input[name="password"]', seededPassword);
    await page.locator("#submitButton").click();

    const token = await getToken(page);
    const wayId = await fetchWayId(page, token);

    const existingComment = `preexisting review ${Date.now()}`;
    await client.query(
      "INSERT INTO reviews (way_id, user_id, rating, comment, created_at) VALUES ($1, $2, $3, $4, now())",
      [wayId, seededUserId, 7, existingComment]
    );

    await page.goto(`${FRONTEND_URL}/?way=${wayId}`);
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
});
