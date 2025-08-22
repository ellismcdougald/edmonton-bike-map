/**
 * utils/auth.ts
 *
 * Purpose:
 * Provides helper functions to extract user information from a JWT stored in localStorage.
 *
 * Functions:
 * - getUsernameFromToken(): returns the `username` from the JWT payload, or null if missing/invalid
 * - getUserIdFromToken(): returns the `userId` from the JWT payload, or null if missing/invalid
 *
 * Behavior:
 * - Reads token from localStorage key 'token'
 * - Decodes the JWT payload using atob and JSON.parse
 * - Safely handles missing tokens, malformed tokens, or missing fields by returning null
 *
 * Notes:
 * - Purely synchronous; does not communicate with backend
 * - Assumes token payload contains `username` and/or `userId`
 * - Useful for identifying the logged-in user on the frontend
 */

export function getUsernameFromToken(): string | null {
	const token = localStorage.getItem('token');
	if (!token) return null;

	try {
		const payload = JSON.parse(atob(token.split('.')[1])) as { username?: string };
		const username = payload.username;
		if (typeof username === 'string') {
			return username;
		} else {
			return null;
		}
	} catch {
		return null;
	}
}

export function getUserIdFromToken(): number | null {
	const token = localStorage.getItem('token');
	if (!token) return null;

	try {
		const payload = JSON.parse(atob(token.split('.')[1])) as { userId?: number };
		const userId = payload.userId;
		if (typeof userId === 'number') {
			return userId;
		} else {
			return null;
		}
	} catch {
		return null;
	}
}
