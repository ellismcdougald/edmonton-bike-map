import { goto } from '$app/navigation';
import type { Review as ReviewObj } from '$lib/types';

export async function fetchReviews(wayId: number): Promise<ReviewObj[]> {
	const res = await fetch(`/api/reviews/${wayId}`);

	if (!res.ok) {
		if (res.status === 401) {
			await fetch('/logout');
			goto('/login');
		}
		throw new Error(`Failed to fetch reviews: ${res.statusText}`);
	}

	return res.json();
}

export async function createReview(wayId: number, rating: number, comment: string | null) {
	const res = await fetch('/api/reviews', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ wayId, rating, comment })
	});

	if (!res.ok) {
		if (res.status === 401) {
			await fetch('/logout');
			goto('/login');
		}
		const text = await res.text();
		throw new Error(text || 'Failed to submit review');
	}

	return;
}
