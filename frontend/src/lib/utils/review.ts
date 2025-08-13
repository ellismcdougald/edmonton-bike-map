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

export interface ReviewInput {
	wayId: number;
	userId: number;
	rating: number;
	comment?: string | null;
}

export async function submitReview({ wayId, userId, rating, comment }: ReviewInput) {
	const res = await fetch('http://localhost:8080/api/reviews', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ wayId, userId, rating, comment })
	});

	if (!res.ok) {
		throw new Error(`Failed to submit review: ${res.statusText}`);
	}

	return res.json();
}
