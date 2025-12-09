import { API_URL } from '$env/static/private';

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
