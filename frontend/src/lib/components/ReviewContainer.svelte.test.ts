import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ReviewContainer from './ReviewContainer.svelte';
import { submitReview } from '$lib/utils/review';

vi.mock('$lib/utils/review', () => ({
	fetchReviews: vi.fn().mockResolvedValue([
		{
			username: 'Alice',
			rating: 8,
			comment: 'Nice ride',
			createdAt: '2025-08-18T00:00:00Z',
			wayID: 1
		}
	]),
	submitReview: vi.fn()
}));

describe('ReviewContainer.svelte', () => {
	let wayId = 1;

	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('loads and displays reviews', async () => {
		const { getByText } = render(ReviewContainer, { wayId });

		await waitFor(() => {
			expect(getByText(/Alice/)).toBeTruthy();
			expect(getByText('Nice ride')).toBeTruthy();
			expect(getByText('Rating: 8 / 10')).toBeTruthy();
		});
	});

	it('opens AddReviewPopup when Add Review button is clicked', async () => {
		const { container, queryByText } = render(ReviewContainer, { wayId });

		const addButton = container.querySelector<HTMLButtonElement>('#addReviewButton')!;
		expect(addButton).toBeInTheDocument();

		await fireEvent.click(addButton);

		await waitFor(() => {
			// The button still exists
			expect(addButton).toBeInTheDocument();

			const popup = container.querySelector('#addReviewPopup');
			expect(popup).toBeInTheDocument();
		});
	});
});
