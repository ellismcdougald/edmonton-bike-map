import { render, fireEvent, waitFor } from '@testing-library/svelte';
import AddReviewPopup from '$lib/components/AddReviewPopup.svelte';
import { vi } from 'vitest';

describe('AddReviewPopup', () => {
	const wayId = 123;
	const closePopup = vi.fn();
	const onReviewAdded = vi.fn();
	const fetchMock = vi.fn();

	beforeEach(() => {
		fetchMock.mockReset();
		vi.stubGlobal('fetch', fetchMock);
	});

	afterEach(() => {
		vi.unstubAllGlobals();
	});

	it('calls submitReview and closes popup on successful submission', async () => {
		fetchMock.mockResolvedValue({ ok: true, json: async () => ({}) });
		const { container } = render(AddReviewPopup, { wayId, closePopup, onReviewAdded });

		const ratingInput = container.querySelector('#rating') as HTMLInputElement;
		const commentInput = container.querySelector('#comment') as HTMLTextAreaElement;
		const submitButton = container.querySelector('#submitButton') as HTMLButtonElement;

		await fireEvent.input(ratingInput, { target: { value: '8' } });
		await fireEvent.input(commentInput, { target: { value: 'Great ride!' } });

		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(fetchMock).toHaveBeenCalledWith('/api/reviews', expect.any(Object));
			expect(closePopup).toHaveBeenCalled();
		});
	});

	it('shows error message if submitReview throws', async () => {
		fetchMock.mockResolvedValue({ ok: false, text: async () => 'Network error' });

		const { container, getByText } = render(AddReviewPopup, { wayId, closePopup, onReviewAdded });

		const ratingInput = container.querySelector('#rating') as HTMLInputElement;
		const submitButton = container.querySelector('#submitButton') as HTMLButtonElement;

		await fireEvent.input(ratingInput, { target: { value: '7' } });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(getByText('Network error')).toBeTruthy();
		});
	});

	it('calls closePopup when Cancel button is clicked', async () => {
		const { container } = render(AddReviewPopup, { wayId, closePopup, onReviewAdded });
		const cancelButton = container.querySelector('#cancelButton') as HTMLButtonElement;

		await fireEvent.click(cancelButton);
		expect(closePopup).toHaveBeenCalled();
	});
});
