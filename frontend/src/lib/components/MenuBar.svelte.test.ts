import { render } from '@testing-library/svelte';
import MenuBar from '$lib/components/MenuBar.svelte';

describe('MenuBar', () => {
	it('renders settings link and logout form button with proper targets', () => {
		const { getByText } = render(MenuBar, { username: 'mock-username' });

		const settingsLink = getByText('Settings') as HTMLAnchorElement;
		expect(settingsLink.getAttribute('href')).toBe('/settings');
		expect(settingsLink.hasAttribute('data-sveltekit-reload')).toBe(true);

		const logoutButton = getByText('Log out') as HTMLButtonElement;
		expect(logoutButton.type).toBe('submit');
		const form = logoutButton.closest('form');
		expect(form).not.toBeNull();
		expect(form?.method.toUpperCase()).toBe('POST');
		expect(form?.action).toContain('/logout');
	});
});
