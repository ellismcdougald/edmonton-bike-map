import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { vi } from 'vitest';
import * as navigation from '$app/navigation';
import SettingsPage from './+page.svelte';

vi.mock('$app/navigation', () => ({
	goto: vi.fn()
}));

// Mock fetch globally
global.fetch = vi.fn();

describe('Settings Page', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		localStorage.clear();
		(global.fetch as any).mockReset();
	});

	it('redirects to login if user is not logged in', async () => {
		render(SettingsPage);

		await waitFor(() => {
			expect(navigation.goto).toHaveBeenCalledWith('/login');
		});
	});

	it('displays change password form when logged in', () => {
		localStorage.setItem('token', 'fake-token');

		const { getByPlaceholderText, getByRole } = render(SettingsPage);

		expect(getByPlaceholderText('Current Password')).toBeTruthy();
		expect(getByPlaceholderText('New Password')).toBeTruthy();
		expect(getByPlaceholderText('Confirm New Password')).toBeTruthy();
		expect(getByRole('button', { name: 'Change Password' })).toBeTruthy();
	});

	it('shows error when passwords do not match', async () => {
		localStorage.setItem('token', 'fake-token');

		const { getByPlaceholderText, getByText, getByRole } = render(SettingsPage);

		const currentPasswordInput = getByPlaceholderText('Current Password') as HTMLInputElement;
		const newPasswordInput = getByPlaceholderText('New Password') as HTMLInputElement;
		const confirmPasswordInput = getByPlaceholderText('Confirm New Password') as HTMLInputElement;
		const submitButton = getByRole('button', { name: 'Change Password' }) as HTMLButtonElement;

		await fireEvent.input(currentPasswordInput, { target: { value: 'oldpassword' } });
		await fireEvent.input(newPasswordInput, { target: { value: 'newpassword' } });
		await fireEvent.input(confirmPasswordInput, { target: { value: 'differentpassword' } });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(getByText('New passwords do not match')).toBeTruthy();
		});
	});

	it('successfully changes password', async () => {
		localStorage.setItem('token', 'fake-token');
		(global.fetch as any).mockResolvedValueOnce({
			ok: true
		});

		const { getByPlaceholderText, getByText, getByRole } = render(SettingsPage);

		const currentPasswordInput = getByPlaceholderText('Current Password') as HTMLInputElement;
		const newPasswordInput = getByPlaceholderText('New Password') as HTMLInputElement;
		const confirmPasswordInput = getByPlaceholderText('Confirm New Password') as HTMLInputElement;
		const submitButton = getByRole('button', { name: 'Change Password' }) as HTMLButtonElement;

		await fireEvent.input(currentPasswordInput, { target: { value: 'oldpassword' } });
		await fireEvent.input(newPasswordInput, { target: { value: 'newpassword' } });
		await fireEvent.input(confirmPasswordInput, { target: { value: 'newpassword' } });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(getByText('Password changed successfully')).toBeTruthy();
		});

		// Check that inputs are cleared
		expect(currentPasswordInput.value).toBe('');
		expect(newPasswordInput.value).toBe('');
		expect(confirmPasswordInput.value).toBe('');
	});

	it('shows error message on failed password change', async () => {
		localStorage.setItem('token', 'fake-token');
		(global.fetch as any).mockResolvedValueOnce({
			ok: false,
			text: async () => 'Current password is incorrect'
		});

		const { getByPlaceholderText, getByText, getByRole } = render(SettingsPage);

		const currentPasswordInput = getByPlaceholderText('Current Password') as HTMLInputElement;
		const newPasswordInput = getByPlaceholderText('New Password') as HTMLInputElement;
		const confirmPasswordInput = getByPlaceholderText('Confirm New Password') as HTMLInputElement;
		const submitButton = getByRole('button', { name: 'Change Password' }) as HTMLButtonElement;

		await fireEvent.input(currentPasswordInput, { target: { value: 'wrongpassword' } });
		await fireEvent.input(newPasswordInput, { target: { value: 'newpassword' } });
		await fireEvent.input(confirmPasswordInput, { target: { value: 'newpassword' } });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(getByText('Current password is incorrect')).toBeTruthy();
		});
	});
});
