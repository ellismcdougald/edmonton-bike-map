import { describe, it, expect } from 'vitest';
import { computeAverageRating } from '$lib/utils/review';
import type { Review } from '$lib/types';

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
