import { goto } from '$app/navigation';

/**
 * Updates the stored cycling speed setting on the server.
 *
 * Sends the new `cyclingSpeed` to the backend; if the server responds with 401 it posts to `/logout` and navigates to the login page.
 *
 * @param cyclingSpeed - New cycling speed value to save in user settings.
 * @throws Error - If the server responds with a non-OK status; message is the response body or 'Failed to save settings'.
 */
export async function updateCyclingSpeed(cyclingSpeed: number) {
	const res = await fetch('/api/settings', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ cyclingSpeed })
	});

	if (!res.ok) {
		if (res.status === 401) {
			await fetch('/logout', { method: 'POST' });
			goto('/login');
		}
		const text = await res.text();
		throw new Error(text || 'Failed to save settings');
	}

	return;
}
