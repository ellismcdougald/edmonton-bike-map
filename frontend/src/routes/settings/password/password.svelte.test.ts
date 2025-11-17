import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { vi, type Mock } from 'vitest';
import * as navigation from '$app/navigation';
import PasswordPage from './+page.svelte';

vi.mock('$app/navigation', () => ({
	goto: vi.fn()
}));

// Mock fetch globally
const mockFetch = vi.fn() as Mock;
global.fetch = mockFetch;

describe('Password Page', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		localStorage.clear();
		mockFetch.mockReset();
	});

	it('redirects to login if not authenticated', async () => {
		render(PasswordPage);

		await waitFor(() => {
			expect(navigation.goto).toHaveBeenCalledWith('/login');
		});
	});

	it('shows form when logged in and password change succeeds', async () => {
		localStorage.setItem('token', 'fake-token');
		mockFetch.mockResolvedValueOnce({ ok: true });

		const { getByPlaceholderText, getByRole, getByText } = render(PasswordPage);

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
	});

	it('shows error when passwords do not match', async () => {
		localStorage.setItem('token', 'fake-token');

		const { getByPlaceholderText, getByText, getByRole } = render(PasswordPage);

		const currentPasswordInput = getByPlaceholderText('Current Password') as HTMLInputElement;
		const newPasswordInput = getByPlaceholderText('New Password') as HTMLInputElement;
		const confirmPasswordInput = getByPlaceholderText('Confirm New Password') as HTMLInputElement;
		const submitButton = getByRole('button', { name: 'Change Password' }) as HTMLButtonElement;

		await fireEvent.input(currentPasswordInput, { target: { value: 'oldpassword' } });
		await fireEvent.input(newPasswordInput, { target: { value: 'newpassword' } });
		await fireEvent.input(confirmPasswordInput, { target: { value: 'different' } });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(getByText('New passwords do not match')).toBeTruthy();
		});
	});

	it('shows backend error when change fails', async () => {
		localStorage.setItem('token', 'fake-token');
		mockFetch.mockResolvedValueOnce({
			ok: false,
			text: async () => 'Current password is incorrect'
		});

		const { getByPlaceholderText, getByText, getByRole } = render(PasswordPage);

		const currentPasswordInput = getByPlaceholderText('Current Password') as HTMLInputElement;
		const newPasswordInput = getByPlaceholderText('New Password') as HTMLInputElement;
		const confirmPasswordInput = getByPlaceholderText('Confirm New Password') as HTMLInputElement;
		const submitButton = getByRole('button', { name: 'Change Password' }) as HTMLButtonElement;

		await fireEvent.input(currentPasswordInput, { target: { value: 'oldpassword' } });
		await fireEvent.input(newPasswordInput, { target: { value: 'newpassword' } });
		await fireEvent.input(confirmPasswordInput, { target: { value: 'newpassword' } });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(getByText('Current password is incorrect')).toBeTruthy();
		});
	});
});
