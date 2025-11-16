import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { vi, type Mock } from 'vitest';
import * as navigation from '$app/navigation';
import SettingsPage from './+page.svelte';

vi.mock('$app/navigation', () => ({
	goto: vi.fn()
}));

// Mock fetch globally
const mockFetch = vi.fn() as Mock;
global.fetch = mockFetch;

describe('Settings Page', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		localStorage.clear();
		mockFetch.mockReset();
	});

	it('redirects to login if user is not logged in', async () => {
		render(SettingsPage);

		await waitFor(() => {
			expect(navigation.goto).toHaveBeenCalledWith('/login');
		});
	});

	it('prefills cycling speed when logged in and shows controls', async () => {
		localStorage.setItem('token', 'fake-token');

		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({ username: 'alice', cyclingSpeed: 18 })
		});

		const { getByLabelText, getByRole, getByText } = render(SettingsPage);

		await waitFor(() => {
			expect(getByLabelText('Preferred Cycling Speed (km/h)')).toBeTruthy();
		});

		const input = getByLabelText('Preferred Cycling Speed (km/h)') as HTMLInputElement;
		expect(input.value).toBe('18');
		expect(getByRole('button', { name: 'Save Preferences' })).toBeTruthy();
		expect(getByText('Change Password')).toBeTruthy();
	});

	it('shows validation error for invalid speed', async () => {
		localStorage.setItem('token', 'fake-token');

		mockFetch.mockResolvedValueOnce({ ok: true, json: async () => ({ cyclingSpeed: 15 }) });

		const { getByLabelText, getByText, getByRole } = render(SettingsPage);
		const input = getByLabelText('Preferred Cycling Speed (km/h)') as HTMLInputElement;
		const saveButton = getByRole('button', { name: 'Save Preferences' }) as HTMLButtonElement;

		await waitFor(() => expect(input).toBeTruthy());

		await fireEvent.input(input, { target: { value: '-5' } });
		await fireEvent.click(saveButton);

		// The error message should appear when an invalid speed is entered.
		await waitFor(() => {
			expect(getByText(/Please enter a valid cycling speed/)).toBeTruthy();
		});

		// No POST should be attempted when validation fails (only the initial GET)
		expect(mockFetch).toHaveBeenCalledTimes(1);
	});

	it('successfully saves cycling speed', async () => {
		localStorage.setItem('token', 'fake-token');
		mockFetch
			.mockResolvedValueOnce({ ok: true, json: async () => ({ cyclingSpeed: 15 }) }) // initial GET
			.mockResolvedValueOnce({ ok: true }); // POST

		const { getByLabelText, getByRole, getByText } = render(SettingsPage);
		const input = await waitFor(
			() => getByLabelText('Preferred Cycling Speed (km/h)') as HTMLInputElement
		);
		const saveButton = getByRole('button', { name: 'Save Preferences' }) as HTMLButtonElement;

		await fireEvent.input(input, { target: { value: '22' } });
		await fireEvent.click(saveButton);

		await waitFor(() => {
			expect(getByText('Settings saved')).toBeTruthy();
		});

		// verify POST body
		expect(mockFetch).toHaveBeenCalledWith(
			expect.stringContaining('/api/user/settings'),
			expect.objectContaining({
				method: 'POST',
				headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
				body: JSON.stringify({ cyclingSpeed: 22 })
			})
		);
	});

	it('shows error message on failed save', async () => {
		localStorage.setItem('token', 'fake-token');
		mockFetch
			.mockResolvedValueOnce({ ok: true, json: async () => ({ cyclingSpeed: 15 }) })
			.mockResolvedValueOnce({ ok: false, text: async () => 'Save failed' });

		const { getByLabelText, getByRole, getByText } = render(SettingsPage);
		const input = await waitFor(
			() => getByLabelText('Preferred Cycling Speed (km/h)') as HTMLInputElement
		);
		const saveButton = getByRole('button', { name: 'Save Preferences' }) as HTMLButtonElement;

		await fireEvent.input(input, { target: { value: '30' } });
		await fireEvent.click(saveButton);

		await waitFor(() => {
			expect(getByText('Save failed')).toBeTruthy();
		});
	});
});
