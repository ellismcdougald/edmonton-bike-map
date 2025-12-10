import { goto } from '$app/navigation';

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
