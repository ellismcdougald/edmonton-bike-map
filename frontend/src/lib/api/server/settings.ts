import { API_URL } from '$env/static/private';

export async function fetchSettings(fetch: typeof globalThis.fetch) {
	const res = await fetch(`${API_URL}/api/settings`);

	if (!res.ok) {
		throw new Error(res.status === 401 ? 'unauthorized' : 'failed');
	}

	return res.json();
}
