import { API_URL } from '$env/static/private';

/**
 * Change the user's password by calling the backend change-password endpoint.
 *
 * @param fetch - A fetch-like function used to perform the HTTP request (e.g., `window.fetch`).
 * @param currentPassword - The user's current password.
 * @param newPassword - The new password to set for the user.
 * @returns `true` when the password was changed successfully.
 * @throws Error with the server response text, or `"Failed to change password"` if the response body is empty, when the request fails.
 */
export async function changePassword(
	fetch: typeof globalThis.fetch,
	currentPassword: string,
	newPassword: string
) {
	const endpoint = `${API_URL}/api/change-password`;

	const res = await fetch(endpoint, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ currentPassword, newPassword })
	});

	if (!res.ok) {
		const errorText = await res.text();
		throw new Error(errorText || 'Failed to change password');
	}

	return true;
}