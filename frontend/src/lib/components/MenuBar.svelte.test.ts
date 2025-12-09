import { render } from '@testing-library/svelte';
import MenuBar from '$lib/components/MenuBar.svelte';

describe('MenuBar', () => {
	it('renders settings and logout links with proper targets', () => {
		const { getByText } = render(MenuBar, { username: 'mock-username' });

		const settingsLink = getByText('Settings') as HTMLAnchorElement;
		expect(settingsLink.getAttribute('href')).toBe('/settings');
		expect(settingsLink.hasAttribute('data-sveltekit-reload')).toBe(true);

		const logoutLink = getByText('Log out') as HTMLAnchorElement;
		expect(logoutLink.getAttribute('href')).toBe('/logout');
	});
});
