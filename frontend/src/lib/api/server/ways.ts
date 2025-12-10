import { redirect } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

/**
 * Retrieve all ways from the server and return the parsed JSON response.
 *
 * @param fetch - The fetch function used to perform the HTTP request.
 * @returns The parsed JSON body returned by the `/api/all-ways` endpoint.
 * @throws Redirect to `/login` with status 302 when the server responds with 401.
 * @throws Error when the response has a non-OK HTTP status.
 */
export async function fetchAllWays(fetch: typeof global.fetch) {
	const response = await fetch(`${API_URL}/api/all-ways`);

	if (response.status === 401) {
		throw redirect(302, '/login');
	}

	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}

	return response.json();
}