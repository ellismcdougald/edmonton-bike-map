import { test, expect } from '@playwright/test';

test.describe('Change Password Feature', () => {
	const USERNAME = 'testuser'; // Ensure this user exists in your test environment
	const INITIAL_PASSWORD = 'password123';
	const NEW_PASSWORD = 'newpassword123';

	// Helper function for login
	async function login(page, username, password) {
		await page.goto('/login');
		await page.fill('input[name="username"]', username);
		await page.fill('input[name="password"]', password);
		await page.click('button[type="submit"]');
		await expect(page).toHaveURL('/'); // Assuming successful login redirects to home
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
			await login(page, USERNAME, INITIAL_PASSWORD);
		});

		// 2. Navigate to settings page
		await test.step('Navigate to settings page', async () => {
			await page.click('[data-testid="nav-settings-link"]');
			await expect(page).toHaveURL('/settings');
		});

		// 3. Change password to new password
		await test.step('Change password to new password', async () => {
			await page.fill('input[name="currentPassword"]', INITIAL_PASSWORD);
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
			await login(page, USERNAME, NEW_PASSWORD);
		});

		// 6. Navigate back to settings and change password to original for test cleanup
		await test.step('Change password back to original for cleanup', async () => {
			await page.click('[data-testid="nav-settings-link"]');
			await expect(page).toHaveURL('/settings');

			await page.fill('input[name="currentPassword"]', NEW_PASSWORD);
			await page.fill('input[name="newPassword"]', INITIAL_PASSWORD);
			await page.fill('input[name="confirmNewPassword"]', INITIAL_PASSWORD);
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
			await login(page, USERNAME, INITIAL_PASSWORD);
			await logout(page);
		});
	});
});
