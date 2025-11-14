import { test, expect } from '@playwright/test';
import { Client } from 'pg';
import * as dotenv from 'dotenv';
import bcrypt from 'bcrypt';

const dotenvResult = dotenv.config({ path: '../../.env' });

if (dotenvResult.error) {
	console.error('Error loading .env file:', dotenvResult.error);
} else {
	console.log('dotenv.config() result (parsed variables):', dotenvResult.parsed);
}

console.log('After dotenv.config, checking env vars:');
console.log(`POSTGRES_TEST_USER: "${process.env.POSTGRES_TEST_USER}"`);
console.log(`POSTGRES_TEST_PASSWORD: "${process.env.POSTGRES_TEST_PASSWORD}"`);
console.log(`POSTGRES_TEST_DB: "${process.env.POSTGRES_TEST_DB}"`);
console.log(`POSTGRES_TEST_PORT: "${process.env.POSTGRES_TEST_PORT}"`);

const parsedPort = parseInt(process.env.POSTGRES_TEST_PORT || '5434');
const dbConfig = {
	user: process.env.POSTGRES_TEST_USER,
	password: process.env.POSTGRES_TEST_PASSWORD,
	database: process.env.POSTGRES_TEST_DB,
	host: 'localhost',
	port: isNaN(parsedPort) ? 5434 : parsedPort
};

console.log('Final DB Config:', dbConfig);

test.describe('Change Password Feature', () => {
	const testUsername = 'testuser';
	const testPassword = 'password123';
	const NEW_PASSWORD = 'newpassword123';

	let client: Client;
	test.beforeAll(async () => {
		try {
			client = new Client(dbConfig);
			await client.connect();
			console.log('Database connected successfully in changePassword.test.ts');
		} catch (error) {
			console.error('Failed to connect to database in changePassword.test.ts:', error);
			throw error; // Re-throw to fail the test early
		}
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

	// Helper function for login
	async function login(page, username, password) {
		await page.goto('/login');
		await page.fill('input[name="username"]', username);
		await page.fill('input[name="password"]', password);
		await page.click('button[type="submit"]');
		// Explicitly wait for the URL to change to '/' after login.
		// This can provide a clearer timeout error if the navigation doesn't occur.
		console.log(`Current URL before waiting for navigation: ${page.url()}`);
		await page.waitForURL('/', { timeout: 60000 }); // Increased timeout for debugging
		console.log('Navigation to "/" detected.');
		// Increase timeout for URL assertion as navigation might take longer.
		await expect(page).toHaveURL('/', { timeout: 10000 }); // Assert the URL after waiting
		await expect(page.locator(`text=${username}`)).toBeVisible(); // Verify user is logged in
	}

	// Helper function for logout
	async function logout(page) {
		await page.click('[data-testid="nav-logout-button"]');
		await expect(page).toHaveURL('/login'); // Assuming logout redirects to login
	}

	test('user can log in, change password, and log in with new password', async ({ page }) => {
		// 1. Log in with initial password
		await test.step('Log in with initial password', async () => {
			await login(page, testUsername, testPassword);
		});

		// 2. Navigate to settings page
		await test.step('Navigate to settings page', async () => {
			await page.click('[data-testid="nav-settings-link"]');
			await expect(page).toHaveURL('/settings');
		});

		// 3. Change password to new password
		await test.step('Change password to new password', async () => {
			await page.fill('input[name="currentPassword"]', testPassword);
			await page.fill('input[name="newPassword"]', NEW_PASSWORD);
			await page.fill('input[name="confirmNewPassword"]', NEW_PASSWORD);
			await page.click('button[data-testid="change-password-submit-button"]');
			await expect(page.locator('[data-testid="password-change-success-message"]')).toBeVisible();
			await expect(page.locator('[data-testid="password-change-success-message"]')).toHaveText('Password changed successfully');
		});

		// 4. Log out
		await test.step('Log out after password change', async () => {
			await logout(page);
		});

		// 5. Log in with the new password
		await test.step('Log in with new password', async () => {
			await login(page, testUsername, NEW_PASSWORD);
		});

		// 6. Navigate back to settings and change password to original for test cleanup
		await test.step('Change password back to original for cleanup', async () => {
			await page.click('[data-testid="nav-settings-link"]');
			await expect(page).toHaveURL('/settings');

			await page.fill('input[name="currentPassword"]', NEW_PASSWORD);
			await page.fill('input[name="newPassword"]', testPassword);
			await page.fill('input[name="confirmNewPassword"]', testPassword);
			await page.click('button[data-testid="change-password-submit-button"]');
			await expect(page.locator('[data-testid="password-change-success-message"]')).toBeVisible();
			await expect(page.locator('[data-testid="password-change-success-message"]')).toHaveText('Password changed successfully');
		});

		// 7. Log out again
		await test.step('Log out after cleanup', async () => {
			await logout(page);
		});

		// Optional: Verify login with initial password still works after cleanup
		await test.step('Verify login with initial password after cleanup', async () => {
			await login(page, testUsername, testPassword);
			await logout(page);
		});
	});
});
