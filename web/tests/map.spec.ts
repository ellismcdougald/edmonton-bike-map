import { test, expect } from "@playwright/test";

test("map loads", async ({ page }) => {
  await page.goto("http://localhost:8080");

  const mapContainer = page.locator("#map");
  await expect(mapContainer).toBeVisible();
});

test("displays route", async ({ page }) => {
  await page.goto("http://localhost:8080");

  const mapContainer = page.locator("#map");
  await expect(mapContainer).toBeVisible();

  // Select start location
  await page.click("#selectStartButton");
  await page.click("#map", { position: { x: 250, y: 250 } });

  // Select end location
  await page.click("#selectEndButton");
  await page.click("#map", { position: { x: 300, y: 300 } });

  // Display route
  await page.click("#findRouteButton");

  // Route is displayed
  const routeLine = page.locator(
    ".leaflet-pane.leaflet-route-pane path.leaflet-interactive"
  );
  await expect(routeLine).toHaveCount(1);
});
