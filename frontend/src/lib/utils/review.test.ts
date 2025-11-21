import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { submitReview, computeAverageRating } from '$lib/utils/review';
import type { Review } from '$lib/types';

const apiUrl = import.meta.env.VITE_API_URL;

describe('submitReview', () => {
	let fetchMock: typeof globalThis.fetch;

	beforeEach(() => {
		fetchMock = vi.fn();
		globalThis.fetch = fetchMock;

		// Create a mock localStorage implementation
		globalThis.localStorage = {
			getItem: vi.fn(),
			setItem: vi.fn(),
			removeItem: vi.fn(),
			clear: vi.fn()
		} as unknown as Storage;
		vi.spyOn(globalThis.localStorage, 'getItem').mockReturnValue('test-token');
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
			headers: { 'Content-Type': 'application/json', Authorization: 'Bearer test-token' },
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

describe('computeAverageRating', () => {
	it('returns NaN for empty array', () => {
		const reviews: Review[] = [];
		expect(computeAverageRating(reviews)).toBeNaN();
	});

	it('returns the rating itself for single review', () => {
		const reviews: Review[] = [
			{ wayId: 1, username: 'alice', rating: 4, comment: 'Nice', createdAt: '2025-09-15' }
		];
		expect(computeAverageRating(reviews)).toBe(4);
	});

	it('calculates average rating rounded to 1 decimal', () => {
		const reviews: Review[] = [
			{ wayId: 1, username: 'alice', rating: 3, comment: '', createdAt: '2025-09-15' },
			{ wayId: 1, username: 'bob', rating: 4, comment: '', createdAt: '2025-09-15' },
			{ wayId: 1, username: 'carol', rating: 5, comment: '', createdAt: '2025-09-15' }
		];
		expect(computeAverageRating(reviews)).toBe(4); // (3+4+5)/3 = 4
	});

	it('rounds average to 1 decimal for non-integer result', () => {
		const reviews: Review[] = [
			{ wayId: 1, username: 'alice', rating: 3, comment: '', createdAt: '2025-09-15' },
			{ wayId: 1, username: 'bob', rating: 4, comment: '', createdAt: '2025-09-15' }
		];
		expect(computeAverageRating(reviews)).toBe(3.5); // (3+4)/2 = 3.5
	});
});
