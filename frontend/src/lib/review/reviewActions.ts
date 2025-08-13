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
