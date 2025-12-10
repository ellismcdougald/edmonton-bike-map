import { API_URL } from '$env/static/private';

/**
 * Fetches application settings from the API.
 *
 * @returns The parsed settings object returned by the API.
 * @throws Error with message 'unauthorized' when the response status is 401, 'failed' for other non-OK responses.
 */
export async function fetchSettings(fetch: typeof globalThis.fetch) {
	const res = await fetch(`${API_URL}/api/settings`);

	if (!res.ok) {
		throw new Error(res.status === 401 ? 'unauthorized' : 'failed');
	}

	return res.json();
}