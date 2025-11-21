/**
 * utils/review.ts
 *
 * Purpose:
 * Provides helper functions to fetch, submit, and process user reviews for ways/routes.
 *
 * Types:
 * - ReviewData: object with `wayId`, `userId`, `rating`, and optional `comment`
 * - ReviewObj: type representing a review retrieved from the backend
 * - FetchFn: type alias for a fetch-like function (used for dependency injection/testing)
 *
 * Functions:
 * - fetchReviews(wayId): retrieves all reviews for a given way from the backend
 *     - Returns an array of ReviewObj
 *     - Returns an empty array on fetch failure
 * - submitReview(reviewData, fetchFn?): sends a new review to the backend
 *     - Throws an error if submission fails
 *     - Returns nothing (caller handles result)
 * - computeAverageRating(reviews): computes the average rating of an array of ReviewObj
 *     - Returns a number rounded to 1 decimal place
 *     - Returns NaN if the array is empty
 *
 * Behavior:
 * - fetchReviews uses the browser's fetch API to GET /api/reviews
 * - submitReview uses fetch to POST /api/reviews with JSON payload
 * - computeAverageRating sums ratings and rounds the average to one decimal
 * - Errors are logged in fetchReviews; submitReview throws errors for the caller
 *
 * Notes:
 * - Default fetch function is browser fetch, but can be overridden (useful for testing)
 * - Depends on backend endpoint: /api/reviews
 */

import type { Review as ReviewObj } from '$lib/types';

const apiUrl = import.meta.env.VITE_API_URL;

export interface ReviewData {
	wayId: number;
	userId: number;
	rating: number;
	comment?: string | null;
}

interface ReviewResponse {
	WayID: string;
	UserID: string;
	Rating: number;
	Comment?: string;
	CreatedAt: string;
	Username: string;
}

export async function fetchReviews(wayId: number): Promise<ReviewObj[]> {
	try {
		const token = localStorage.getItem('token');
		const res = await fetch(`${apiUrl}/api/reviews?wayID=${wayId}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			}
		});
		if (res.status === 401) throw new Error('Unauthorized');
		if (!res.ok) throw new Error(`HTTP error! Status: ${res.status}`);
		const data: ReviewResponse[] = await res.json();

		return data.map((r: ReviewResponse) => ({
			wayId: Number(r.WayID),
			rating: r.Rating,
			comment: r.Comment ?? '',
			createdAt: r.CreatedAt,
			username: r.Username
		}));
	} catch (e) {
		console.error('Failed to fetch reviews: ', e);
		return [];
	}
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
