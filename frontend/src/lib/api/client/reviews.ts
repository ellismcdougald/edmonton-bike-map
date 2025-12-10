import { goto } from '$app/navigation';
import type { Review as ReviewObj } from '$lib/types';

/**
 * Submit a new review for a way.
 *
 * Sends the review data to the backend; on authorization failure it logs out and navigates to the login page.
 *
 * @param wayId - Identifier of the way being reviewed
 * @param rating - Numerical rating value for the way
 * @param comment - Optional textual comment for the review, or `null` if none
 * @throws Error - When the server responds with a non-OK, non-401 status; message contains server-provided text or 'Failed to submit review'
 */
export async function createReview(wayId: number, rating: number, comment: string | null) {
	const res = await fetch('/api/reviews', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ wayId, rating, comment })
	});

	if (!res.ok) {
		if (res.status === 401) {
			await fetch('/logout', { method: 'POST' });
			goto('/login');
			return;
		}
		const text = await res.text();
		throw new Error(text || 'Failed to submit review');
	}
}

/**
 * Retrieve reviews for a specific way.
 *
 * If the server responds with 401, the client will post to /logout, navigate to /login, and return an empty array. For other non-OK responses an Error is thrown with the response status text.
 *
 * @param wayId - ID of the way whose reviews should be fetched
 * @returns An array of `ReviewObj` for the specified way (empty array after redirect on 401)
 */
export async function fetchReviews(wayId: number): Promise<ReviewObj[]> {
	const res = await fetch(`/api/reviews/${wayId}`);

	if (!res.ok) {
		if (res.status === 401) {
			await fetch('/logout', { method: 'POST' });
			goto('/login');
			return [] as ReviewObj[];
		}
		throw new Error(`Failed to fetch reviews: ${res.statusText}`);
	}

	return res.json();
}