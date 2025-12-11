import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ReviewContainer from './ReviewContainer.svelte';
import type { Review } from '$lib/types';

describe('ReviewContainer.svelte', () => {
	const wayId = 1;

	const mockReviews: Review[] = [
		{
			username: 'Alice',
			rating: 8,
			comment: 'Nice ride',
			createdAt: '2025-08-18T00:00:00Z',
			wayId: 1
		}
	];

	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('displays passed-in reviews', async () => {
		const { getByText } = render(ReviewContainer, { props: { wayId, reviews: mockReviews } });

		await waitFor(() => {
			expect(getByText(/Alice/)).toBeTruthy();
			expect(getByText('Nice ride')).toBeTruthy();
			expect(getByText('Rating: 8 / 10')).toBeTruthy();
		});
	});

	it('opens AddReviewPopup when Add Review button is clicked', async () => {
		const { container } = render(ReviewContainer, { props: { wayId, reviews: mockReviews } });

		const addButton = container.querySelector<HTMLButtonElement>('#addReviewButton')!;
		expect(addButton).toBeInTheDocument();

		await fireEvent.click(addButton);

		await waitFor(() => {
			const popup = container.querySelector('#addReviewPopup');
			expect(popup).toBeInTheDocument();
		});
	});

	it('hides Add Review for guests and shows prompt', async () => {
		const { container, getByText } = render(ReviewContainer, {
			props: { wayId, reviews: mockReviews, canReview: false }
		});

		await waitFor(() => {
			expect(getByText('Log in to add a review.')).toBeTruthy();
		});

		const addButton = container.querySelector('#addReviewButton');
		expect(addButton).toBeNull();
	});
});
