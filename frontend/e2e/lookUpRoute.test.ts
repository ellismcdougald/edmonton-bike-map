import { expect, test } from '@playwright/test';

test('user can log in, then select a start and end point and find a route between those two points', async ({
	page
}) => {
	const testUsername = 'test-user';
	const testPassword = 'test-password';

	await page.goto('./');

	await expect(page).toHaveURL('/login');

	await page.fill('input[name="username"]', testUsername);
	await page.fill('input[name="password"]', testPassword);
	await page.locator('#submitButton').click();

	await expect(page).toHaveURL('/');

	const map = page.locator('#map');
	await expect(map).toBeVisible();
	const mapBox = await map.boundingBox();
	if (!mapBox) throw new Error('Map element not found');

	const startButton = page.locator('#selectStartButton');
	await startButton.click();

	await page.mouse.click(mapBox.x + mapBox.width / 2, mapBox.y + mapBox.height / 2);

	await page.locator('#selectEndButton').click();
	await page.mouse.click(mapBox.x + mapBox.width / 3, mapBox.y + mapBox.height / 3);

	await page.locator('#findRouteButton').click();

	const routePaths = page.locator('path.leaflet-interactive');
	await expect(routePaths.first()).toBeVisible({ timeout: 5000 });
	const count = await routePaths.count();
	expect(count).toBeGreaterThan(0);
});
