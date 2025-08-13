import type { Review as ReviewObj } from '$lib/types';

export async function fetchReviews(wayId: number): Promise<ReviewObj[]> {
	try {
		const res = await fetch(`http://localhost:8080/api/reviews?wayID=${wayId}`);
		if (!res.ok) throw new Error(`HTTP error! Status: ${res.status}`);
		return await res.json();
	} catch (e) {
		console.error('Failed to fetch reviews: ', e);
		return [];
	}
}

export interface ReviewData {
	wayId: number;
	userId: number;
	rating: number;
	comment?: string | null;
}

type FetchFn = typeof fetch;

export async function submitReview(reviewData: ReviewData, fetchFn: FetchFn = fetch) {
	const res = await fetchFn('http://localhost:8080/api/reviews', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(reviewData)
	});

	if (!res.ok) {
		throw new Error(`Failed to submit review: ${res.statusText}`);
	}

	return res.json();
}
