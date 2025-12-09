import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import PasswordPage from './+page.svelte';

describe('Password Page', () => {
	it('renders form fields', () => {
		const { getByPlaceholderText, getByRole } = render(PasswordPage, { form: null });
		expect(getByPlaceholderText('Current Password')).toBeTruthy();
		expect(getByPlaceholderText('New Password')).toBeTruthy();
		expect(getByPlaceholderText('Confirm New Password')).toBeTruthy();
		expect(getByRole('button', { name: 'Change Password' })).toBeTruthy();
	});

	it('shows success message and clears fields when form success is returned', async () => {
		const { getByPlaceholderText, getByText, rerender } = render(PasswordPage, { form: null });

		const currentPasswordInput = getByPlaceholderText('Current Password') as HTMLInputElement;
		const newPasswordInput = getByPlaceholderText('New Password') as HTMLInputElement;
		const confirmPasswordInput = getByPlaceholderText('Confirm New Password') as HTMLInputElement;

		await fireEvent.input(currentPasswordInput, { target: { value: 'oldpassword' } });
		await fireEvent.input(newPasswordInput, { target: { value: 'newpassword' } });
		await fireEvent.input(confirmPasswordInput, { target: { value: 'newpassword' } });

		await rerender({ form: { success: true } });

		await waitFor(() => {
			expect(getByText('Password changed successfully')).toBeTruthy();
			expect(currentPasswordInput.value).toBe('');
			expect(newPasswordInput.value).toBe('');
			expect(confirmPasswordInput.value).toBe('');
		});
	});

	it('shows error message from form data', async () => {
		const { getByText, rerender } = render(PasswordPage, { form: null });

		await rerender({ form: { error: 'Current password is incorrect' } });

		await waitFor(() => {
			expect(getByText('Current password is incorrect')).toBeTruthy();
		});
	});
});
