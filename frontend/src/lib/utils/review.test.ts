import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { submitReview } from '$lib/utils/review';

const apiUrl = import.meta.env.VITE_API_URL;

describe('submitReview', () => {
	let fetchMock: typeof globalThis.fetch;

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

		fetchMock = vi.fn().mockResolvedValueOnce(mockResponse);
		globalThis.fetch = fetchMock;

		const reviewData = { wayId: 1, userId: 42, rating: 8, comment: 'Great route' };

		await submitReview(reviewData); // just call normally

		expect(fetchMock).toHaveBeenCalledTimes(1);
		expect(fetchMock).toHaveBeenCalledWith(`${apiUrl}/api/reviews`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(reviewData)
		});
	});

	it('throws an error if fetch response is not ok', async () => {
		const mockResponse = { ok: false, statusText: 'Bad Request' } as Response;

		fetchMock = vi.fn().mockResolvedValueOnce(mockResponse);
		globalThis.fetch = fetchMock;

		const reviewData = { wayId: 1, userId: 42, rating: 8, comment: 'Great route' };

		await expect(submitReview(reviewData)).rejects.toThrow('Failed to submit review: Bad Request');
	});
});
