import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import SettingsPage from './+page.svelte';

const fetchMock = vi.fn();

describe('Settings Page', () => {
	beforeEach(() => {
		fetchMock.mockReset();
		vi.stubGlobal('fetch', fetchMock);
	});

	afterEach(() => {
		vi.unstubAllGlobals();
	});

	it('prefills cycling speed and shows controls', async () => {
		const { getByLabelText, getByRole, getByText } = render(SettingsPage, {
			data: { cyclingSpeed: 18 }
		});

		const input = getByLabelText('Preferred Cycling Speed (km/h)') as HTMLInputElement;
		expect(input.value).toBe('18');
		expect(getByRole('button', { name: 'Save Preferences' })).toBeTruthy();
		expect(getByText('Change Password')).toBeTruthy();
	});

	it('shows validation error for invalid speed', async () => {
		const { getByLabelText, getByText, getByRole } = render(SettingsPage, {
			data: { cyclingSpeed: 15 }
		});
		const input = getByLabelText('Preferred Cycling Speed (km/h)') as HTMLInputElement;
		const saveButton = getByRole('button', { name: 'Save Preferences' }) as HTMLButtonElement;

		await fireEvent.input(input, { target: { value: '-5' } });
		await fireEvent.click(saveButton);

		await waitFor(() => {
			expect(getByText(/Please enter a valid cycling speed/)).toBeTruthy();
		});

		expect(fetchMock).not.toHaveBeenCalled();
	});

	it('successfully saves cycling speed', async () => {
		fetchMock.mockResolvedValueOnce({ ok: true });

		const { getByLabelText, getByRole, getByText } = render(SettingsPage, {
			data: { cyclingSpeed: 15 }
		});
		const input = getByLabelText('Preferred Cycling Speed (km/h)') as HTMLInputElement;
		const saveButton = getByRole('button', { name: 'Save Preferences' }) as HTMLButtonElement;

		await fireEvent.input(input, { target: { value: '22' } });
		await fireEvent.click(saveButton);

		await waitFor(() => {
			expect(getByText('Settings saved')).toBeTruthy();
		});

		expect(fetchMock).toHaveBeenCalledWith(
			'/api/settings',
			expect.objectContaining({
				method: 'POST',
				headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
				body: JSON.stringify({ cyclingSpeed: 22 })
			})
		);
	});

	it('shows error message on failed save', async () => {
		fetchMock.mockResolvedValueOnce({ ok: false, text: async () => 'Save failed' });

		const { getByLabelText, getByRole, getByText } = render(SettingsPage, {
			data: { cyclingSpeed: 15 }
		});
		const input = getByLabelText('Preferred Cycling Speed (km/h)') as HTMLInputElement;
		const saveButton = getByRole('button', { name: 'Save Preferences' }) as HTMLButtonElement;

		await fireEvent.input(input, { target: { value: '30' } });
		await fireEvent.click(saveButton);

		await waitFor(() => {
			expect(getByText('Save failed')).toBeTruthy();
		});
	});
});
