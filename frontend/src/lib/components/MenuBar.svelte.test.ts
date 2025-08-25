import { render, fireEvent } from '@testing-library/svelte';
import { removeToken } from '$lib/utils/auth';
import { vi } from 'vitest';
import * as navigation from '$app/navigation';
import MenuBar from '$lib/components/MenuBar.svelte';

vi.mock('$lib/utils/auth', () => ({
	getUsernameFromToken: vi.fn(() => 'mock-username'),
	removeToken: vi.fn()
}));

vi.mock('$app/navigation', () => ({
	goto: vi.fn()
}));

describe('MenuBar', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('calls removeToken and redirects via goto on logout click', async () => {
		const { getByText } = render(MenuBar);

		// Since your dropdown appears on hover via Tailwind, we can query the logout button directly
		const logoutButton = getByText('Log out');

		await fireEvent.click(logoutButton);

		expect(removeToken).toHaveBeenCalled();
		expect(navigation.goto).toHaveBeenCalledWith('/login');
	});
});
