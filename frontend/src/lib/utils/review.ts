import type { Review as ReviewObj } from '$lib/types';

/**
 * Calculate the average rating from a list of reviews, rounded to one decimal place.
 *
 * @param reviews - Array of review objects whose `rating` values will be averaged
 * @returns The average of all `rating` values rounded to one decimal place, or `NaN` if `reviews` is empty
 */
export function computeAverageRating(reviews: ReviewObj[]): number {
	if (reviews.length == 0) return NaN;

	let total = 0;
	for (const review of reviews) {
		total += review.rating;
	}

	const avg = total / reviews.length;
	return Math.round(avg * 10) / 10;
}
