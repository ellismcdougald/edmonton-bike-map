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
