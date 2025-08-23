import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach } from 'vitest';
import Sidebar from './Sidebar.svelte';
import { wayState } from '$lib/state.svelte';

vi.mock('./ReviewContainer.svelte', () => ({
	default: () => '<div data-testid="review-container"></div>'
}));

// deterministic mock way
const mockWay = {
	id: 1,
	tags: {
		name: 'Test Route',
		surface: 'asphalt',
		highway: 'cycleway',
		foot: 'yes',
		bicycle: 'designated',
		lcn: 'yes'
	}
};

describe('Sidebar.svelte', () => {
	beforeEach(() => {
		wayState.selectedWay = mockWay;
	});

	it('renders route info correctly', async () => {
		const { getByText } = render(Sidebar);

		await waitFor(() => {
			expect(getByText('Shared Use Path')).toBeTruthy();
			expect(getByText('Yes')).toBeTruthy();
			expect(getByText('asphalt')).toBeTruthy();
			expect(getByText('8 / 10')).toBeTruthy();
		});
	});

	it('toggles sidebar visibility when hide button is clicked', async () => {
		const { container } = render(Sidebar);

		const hideButton = container.querySelector('#hide-sidebar-button');
		expect(hideButton).toBeTruthy();

		await fireEvent.click(hideButton!); // click to hide
		await waitFor(() => {
			const sidebarContent = container.querySelector('#sidebar-content');
			expect(sidebarContent).toBeNull();
		});

		await fireEvent.click(hideButton!); // click to show
		await waitFor(() => {
			const sidebarContent = container.querySelector('#sidebar-content');
			expect(sidebarContent).toBeTruthy();
		});
	});
});
