import { goto } from '$app/navigation';
import { API_URL } from '$env/static/private';

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
			await fetch('/logout');
			goto('/login');
		}
		const text = await res.text();
		throw new Error(text || 'Failed to save settings');
	}

	return;
}

export async function fetchSettings(fetch: typeof globalThis.fetch) {
	const res = await fetch(`${API_URL}/api/settings`);

	if (!res.ok) {
		throw new Error(res.status === 401 ? 'unauthorized' : 'failed');
	}

	return res.json();
}
