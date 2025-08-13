import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { submitReview } from './review';

describe('submitReview', () => {
	let fetchMock: ReturnType<typeof vi.fn>;

	beforeEach(() => {
		fetchMock = vi.fn();
		globalThis.fetch = fetchMock;
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('calls fetch with correct URL and payload', async () => {
		const mockResponse = {
			ok: true,
			json: async () => ({ success: true })
		} as Response;

		fetchMock.mockResolvedValueOnce(mockResponse);

		const reviewData = { wayId: 1, userId: 42, rating: 8, comment: 'Great route' };

		await submitReview(reviewData, fetchMock as typeof fetch);

		expect(fetchMock).toHaveBeenCalledTimes(1);
		expect(fetchMock).toHaveBeenCalledWith('http://localhost:8080/api/reviews', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(reviewData)
		});
	});

	it('throws an error if fetch response is not ok', async () => {
		const mockResponse = { ok: false, statusText: 'Bad Request' } as Response;

		fetchMock.mockResolvedValueOnce(mockResponse);

		const reviewData = { wayId: 1, userId: 42, rating: 8, comment: 'Great route' };

		await expect(submitReview(reviewData, fetchMock as typeof fetch)).rejects.toThrow(
			'Failed to submit review: Bad Request'
		);
	});
});
