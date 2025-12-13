import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ReviewContainer from './ReviewContainer.svelte';
import type { Review } from '$lib/types';
import * as reviewsApi from '$lib/api/client/reviews';
import * as waysApi from '$lib/api/client/ways';

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
		// Mock getAdjacentWays to prevent fetch errors in tests
		vi.spyOn(waysApi, 'getAdjacentWays').mockResolvedValue({
			type: 'FeatureCollection',
			features: []
		});
	});

	it('displays passed-in reviews', async () => {
		const { getByText } = render(ReviewContainer, { props: { wayId, reviews: mockReviews } });

		await waitFor(() => {
			expect(getByText(/Alice/)).toBeTruthy();
			expect(getByText('Nice ride')).toBeTruthy();
			expect(getByText('Rating: 8 / 10')).toBeTruthy();
		});
	});

	it('shows form when Add Review is clicked and submits', async () => {
		const createReviewSpy = vi.spyOn(reviewsApi, 'createReview').mockResolvedValue();
		const onReviewAdded = vi.fn();

		const { container } = render(ReviewContainer, {
			props: { wayId, reviews: mockReviews, onReviewAdded }
		});

		const addButton = container.querySelector<HTMLButtonElement>('#addReviewButton')!;
		expect(addButton).toBeInTheDocument();

		await fireEvent.click(addButton);

		const ratingInput = container.querySelector<HTMLInputElement>('#rating')!;
		const commentInput = container.querySelector<HTMLTextAreaElement>('#comment')!;
		const submitButton = container.querySelector<HTMLButtonElement>('#submitButton')!;

		ratingInput.value = '7';
		await fireEvent.input(ratingInput);
		commentInput.value = 'Great route!';
		await fireEvent.input(commentInput);

		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(createReviewSpy).toHaveBeenCalledWith([wayId], 7, 'Great route!');
			expect(onReviewAdded).toHaveBeenCalled();
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
