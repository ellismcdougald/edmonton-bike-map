import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import Review from './Review.svelte';
import type { Review as ReviewType } from '$lib/types';
import * as reviewsApi from '$lib/api/client/reviews';

const review: ReviewType = {
	wayId: 10,
	username: 'me',
	rating: 6,
	comment: 'ok',
	createdAt: '2025-09-01'
};

describe('Review.svelte', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('calls delete API and onDeleted when current user deletes', async () => {
		const deleteSpy = vi.spyOn(reviewsApi, 'deleteReview').mockResolvedValue();
		const onDeleted = vi.fn();

		const { getByText } = render(Review, {
			props: { review, currentUser: 'me', wayId: review.wayId, onDeleted }
		});

		const deleteBtn = getByText('Delete');
		await fireEvent.click(deleteBtn);

		await waitFor(() => {
			expect(deleteSpy).toHaveBeenCalledWith(review.wayId);
			expect(onDeleted).toHaveBeenCalled();
		});
	});

	it('hides delete button for other users', async () => {
		const { queryByText } = render(Review, {
			props: { review, currentUser: 'someone-else', wayId: review.wayId }
		});

		await waitFor(() => {
			expect(queryByText('Delete')).toBeNull();
		});
	});
});
