import type { Review as ReviewObj } from '$lib/types';

const apiUrl = import.meta.env.VITE_API_URL;

export interface ReviewData {
	wayId: number;
	userId: number;
	rating: number;
	comment?: string | null;
}

type FetchFn = typeof fetch;

export async function submitReview(reviewData: ReviewData, fetchFn: FetchFn = fetch) {
	const token = localStorage.getItem('token');
	const res = await fetchFn(`${apiUrl}/api/reviews`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
		body: JSON.stringify(reviewData)
	});

	if (!res.ok) {
		throw new Error(`Failed to submit review: ${res.statusText}`);
	}
}

export function computeAverageRating(reviews: ReviewObj[]): number {
	if (reviews.length == 0) return NaN;

	let total = 0;
	for (const review of reviews) {
		total += review.rating;
	}

	const avg = total / reviews.length;
	return Math.round(avg * 10) / 10;
}
