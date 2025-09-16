import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach } from 'vitest';
import Sidebar from './Sidebar.svelte';
import { wayState } from '$lib/state.svelte';
import type { Review } from '$lib/types';
import * as reviewUtils from '$lib/utils/review';

vi.mock('./ReviewContainer.svelte', () => ({
	default: () => '<div data-testid="review-container"></div>'
}));

const mockReviews: Review[] = [
	{ wayID: 1, username: 'alice', rating: 3, comment: '', createdAt: '2025-09-15' },
	{ wayID: 1, username: 'bob', rating: 4, comment: '', createdAt: '2025-09-15' },
	{ wayID: 1, username: 'carol', rating: 5, comment: '', createdAt: '2025-09-15' }
];

// deterministic mock way
const mockWay = {
	id: 1,
	tags: {
		name: 'Test Route',
		surface: 'Asphalt',
		highway: 'cycleway',
		foot: 'yes',
		bicycle: 'designated',
		lcn: 'yes'
	}
};

describe('Sidebar.svelte', () => {
	beforeEach(() => {
		wayState.selectedWay = mockWay;
		vi.spyOn(reviewUtils, 'fetchReviews').mockResolvedValue(mockReviews);
	});

	it('renders route info correctly', async () => {
		const { getByText } = render(Sidebar);

		await waitFor(() => {
			expect(getByText('Shared Use Path')).toBeTruthy();
			expect(getByText('Yes')).toBeTruthy();
			expect(getByText('Asphalt')).toBeTruthy();
			expect(getByText('4 / 10')).toBeTruthy();
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
