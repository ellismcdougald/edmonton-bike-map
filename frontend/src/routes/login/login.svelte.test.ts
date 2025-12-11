import { render } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import LoginPage from './+page.svelte';

describe('Login page', () => {
	it('links to map for guest viewing', () => {
		const { getByText } = render(LoginPage, { form: { error: '' } });

		const guestLink = getByText('View map as guest') as HTMLAnchorElement;
		expect(guestLink).toBeInTheDocument();
		expect(guestLink.getAttribute('href')).toBe('/map');
	});
});
