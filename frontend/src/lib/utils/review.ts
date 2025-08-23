/**
 * utils/review.ts
 *
 * Purpose:
 * Provides helper functions to fetch and submit user reviews for ways/routes.
 *
 * Types:
 * - ReviewData: object with `wayId`, `userId`, `rating`, and optional `comment`
 * - FetchFn: type alias for a fetch-like function (used for dependency injection/testing)
 *
 * Functions:
 * - fetchReviews(wayId): retrieves all reviews for a given way from the backend
 *     - Returns an array of Review objects
 *     - Returns an empty array on fetch failure
 * - submitReview(reviewData, fetchFn?): sends a new review to the backend
 *     - Throws an error if submission fails
 *     - Returns the created review as JSON
 *
 * Behavior:
 * - fetchReviews uses the browser's fetch API to GET /api/reviews
 * - submitReview uses fetch to POST /api/reviews with JSON payload
 * - Errors are handled gracefully and logged in fetchReviews; submitReview throws errors for the caller to handle
 *
 * Notes:
 * - Default fetch function is browser fetch, but can be overridden (useful for testing)
 * - Depends on backend endpoints: /api/reviews
 */

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
