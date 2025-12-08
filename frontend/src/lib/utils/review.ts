import type { Review as ReviewObj } from '$lib/types';

export function computeAverageRating(reviews: ReviewObj[]): number {
	if (reviews.length == 0) return NaN;

	let total = 0;
	for (const review of reviews) {
		total += review.rating;
	}

	const avg = total / reviews.length;
	return Math.round(avg * 10) / 10;
}
