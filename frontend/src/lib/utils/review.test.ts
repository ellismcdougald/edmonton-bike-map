import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { submitReview } from './review';
import type { Mock } from 'vitest';

describe('submitReview', () => {
	let fetchMock: Mock;

	beforeEach(() => {
		fetchMock = vi.fn() as Mock;
		// @ts-ignore
		globalThis.fetch = fetchMock;
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('calls fetch with correct URL and payload', async () => {
		fetchMock.mockResolvedValueOnce({
			ok: true,
			json: async () => ({ success: true })
		});

		const reviewData = { wayId: 1, userId: 42, rating: 8, comment: 'Great route' };

		await submitReview(reviewData, fetchMock);

		expect(fetchMock).toHaveBeenCalledTimes(1);
		expect(fetchMock).toHaveBeenCalledWith('http://localhost:8080/api/reviews', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(reviewData)
		});
	});

	it('throws an error if fetch response is not ok', async () => {
		fetchMock.mockResolvedValueOnce({ ok: false, statusText: 'Bad Request' });

		const reviewData = { wayId: 1, userId: 42, rating: 8, comment: 'Great route' };

		await expect(submitReview(reviewData, fetchMock)).rejects.toThrow(
			'Failed to submit review: Bad Request'
		);
	});
});
