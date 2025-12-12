import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AddReview from './AddReview.svelte';
import * as reviewsApi from '$lib/api/client/reviews';
import type { WayFeature } from '$lib/types';

// Provide a hoisted mock for wayState without referencing test-scoped vars
vi.mock('$lib/state.svelte', () => {
	const wayState: { selectedWay: WayFeature | null } = { selectedWay: null };
	return { wayState };
});
import { wayState } from '$lib/state.svelte';

describe('AddReview.svelte', () => {
	const wayId = 42;

	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('renders Included Ways pill using selectedWay name when available', async () => {
		// set selectedWay via the mocked state
		wayState.selectedWay = { id: wayId, tags: { name: 'Sample Way' } } as WayFeature;

		const { container, getByText } = render(AddReview, { props: { wayId } });
		await waitFor(() => {
			expect(getByText('Included Ways')).toBeTruthy();
			const pill = container.querySelector('#includedWays');
			expect(pill?.textContent).toContain('Sample Way');
		});
	});

	it('submits rating/comment, calls API and onSubmitted', async () => {
		const createReviewSpy = vi.spyOn(reviewsApi, 'createReview').mockResolvedValue();
		const onSubmitted = vi.fn();

		// provide a basic selectedWay
		wayState.selectedWay = { id: wayId, tags: {} } as WayFeature;

		const { container } = render(AddReview, { props: { wayId, onSubmitted } });

		const ratingInput = container.querySelector<HTMLInputElement>('#rating')!;
		const commentInput = container.querySelector<HTMLTextAreaElement>('#comment')!;
		const submitButton = container.querySelector<HTMLButtonElement>('#submitButton')!;

		ratingInput.value = '9';
		await fireEvent.input(ratingInput);
		commentInput.value = 'Love this way!';
		await fireEvent.input(commentInput);

		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(createReviewSpy).toHaveBeenCalledWith(wayId, 9, 'Love this way!');
			expect(onSubmitted).toHaveBeenCalled();
		});
	});
});
