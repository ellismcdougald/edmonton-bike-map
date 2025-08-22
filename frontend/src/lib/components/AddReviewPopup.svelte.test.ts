import { render, fireEvent, waitFor } from '@testing-library/svelte';
import AddReviewPopup from '$lib/components/AddReviewPopup.svelte';
import { submitReview } from '$lib/utils/review';
import { getUserIdFromToken } from '$lib/utils/auth';
import { vi } from 'vitest';

vi.mock('$lib/utils/auth', () => ({
	getUserIdFromToken: vi.fn(() => 'mock-user-id')
}));

vi.mock('$lib/utils/review', () => ({
	submitReview: vi.fn()
}));

describe('AddReviewPopup', () => {
	const wayId = 123;
	const closePopup = vi.fn();

	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('calls submitReview and closes popup on successful submission', async () => {
		const { container } = render(AddReviewPopup, { wayId, closePopup });

		const ratingInput = container.querySelector('#rating') as HTMLInputElement;
		const commentInput = container.querySelector('#comment') as HTMLTextAreaElement;
		const submitButton = container.querySelector('#submitButton') as HTMLButtonElement;

		await fireEvent.input(ratingInput, { target: { value: '8' } });
		await fireEvent.input(commentInput, { target: { value: 'Great ride!' } });

		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(submitReview).toHaveBeenCalledWith({
				wayId,
				userId: 'mock-user-id',
				rating: 8,
				comment: 'Great ride!'
			});
			expect(closePopup).toHaveBeenCalled();
		});
	});

	it('shows error message if submitReview throws', async () => {
		(submitReview as any).mockRejectedValueOnce(new Error('Network error'));

		const { container, getByText } = render(AddReviewPopup, { wayId, closePopup });

		const ratingInput = container.querySelector('#rating') as HTMLInputElement;
		const submitButton = container.querySelector('#submitButton') as HTMLButtonElement;

		await fireEvent.input(ratingInput, { target: { value: '7' } });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(getByText('Network error')).toBeTruthy();
		});
	});

	it('calls closePopup when Cancel button is clicked', async () => {
		const { container } = render(AddReviewPopup, { wayId, closePopup });
		const cancelButton = container.querySelector('#cancelButton') as HTMLButtonElement;

		await fireEvent.click(cancelButton);
		expect(closePopup).toHaveBeenCalled();
	});
});
