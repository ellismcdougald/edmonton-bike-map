import { expect, test } from '@playwright/test';
import { Client } from 'pg';
import * as dotenv from 'dotenv';
import bcrypt from 'bcrypt';

dotenv.config({ path: '../.env' });

const dbConfig = {
	user: process.env.POSTGRES_TEST_USER,
	password: process.env.POSTGRES_TEST_PASSWORD,
	database: process.env.POSTGRES_TEST_DB,
	host: 'localhost',
	port: parseInt(process.env.POSTGRES_TEST_PORT || '5434')
};

test.describe('route finding with a test user', () => {
	const testUsername = 'test-user';
	const testPassword = 'test-password';

	let client: Client;
	test.beforeAll(async () => {
		const dbUrl = `postgres://${dbConfig.user}:${dbConfig.password}@${dbConfig.host}:${dbConfig.port}/${dbConfig.database}?sslmode=disable`;

		client = new Client(dbConfig);
		await client.connect();

		const result = await client.query('SELECT current_database(), current_user');
		console.log('Frontend DB:', result.rows[0]);
	});

	test.afterAll(async () => {
		await client.end();
	});

	test.beforeEach(async () => {
		// Clear users table and insert test user
		await client.query('TRUNCATE TABLE users RESTART IDENTITY CASCADE');
		const hashedPassword = await bcrypt.hash(testPassword, bcrypt.genSaltSync());
		await client.query('INSERT INTO users (username, password) VALUES ($1, $2)', [
			testUsername,
			hashedPassword
		]);
	});

	test.afterEach(async () => {
		// Clear users after test
		await client.query('TRUNCATE TABLE users RESTART IDENTITY CASCADE');
	});

	test('user can log in, then select a start and end point and find a route between those two points', async ({
		page
	}) => {
		await page.goto('http://localhost:4173/');

		await expect(page).toHaveURL('http://localhost:4173/login');

		await page.fill('input[name="username"]', testUsername);
		await page.fill('input[name="password"]', testPassword);
		await page.locator('#submitButton').click();

		const map = page.locator('#map');
		await expect(map).toBeVisible({ timeout: 10000 });
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

		const distanceControl = page.locator('.distance-control');
		await expect(distanceControl).toBeVisible({ timeout: 5000 });
		const distanceText = await distanceControl.textContent();
		expect(distanceText).toMatch(/km$/); // ends with "km"
		const timeControl = page.locator('.time-control');
		await expect(timeControl).toBeVisible({ timeout: 5000 });
		const timeText = await timeControl.textContent();
		expect(timeText).toMatch(/min$/); // ends with "min"
	});

	test('user can log in, then select a start and end point, find a route between those two points, then reset the map', async ({
		page
	}) => {
		await page.goto('http://localhost:4173/');

		await expect(page).toHaveURL('http://localhost:4173/login');

		await page.fill('input[name="username"]', testUsername);
		await page.fill('input[name="password"]', testPassword);
		await page.locator('#submitButton').click();

		const map = page.locator('#map');
		await expect(map).toBeVisible({ timeout: 10000 });
		const mapBox = await map.boundingBox();
		if (!mapBox) throw new Error('Map element not found');

		const startButton = page.locator('#selectStartButton');
		await startButton.click();

		await page.mouse.click(mapBox.x + mapBox.width / 2, mapBox.y + mapBox.height / 2);

		const markers = page.locator('.leaflet-marker-icon');
		expect(await markers.count()).toBe(1);

		const endButton = page.locator('#selectEndButton');
		await endButton.click();
		await page.mouse.click(mapBox.x + mapBox.width / 3, mapBox.y + mapBox.height / 3);

		expect(await markers.count()).toBe(2);

		await page.locator('#findRouteButton').click();

		const routePaths = page.locator('path.leaflet-interactive');
		await expect(routePaths.first()).toBeVisible({ timeout: 5000 });
		const count = await routePaths.count();
		expect(count).toBeGreaterThan(0);

		await page.locator('#resetButton').click();

		const countAfterReset = await routePaths.count();
		expect(countAfterReset).toBe(0);

		expect(await markers.count()).toBe(0);
	});
});
