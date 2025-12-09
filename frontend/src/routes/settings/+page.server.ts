import { API_URL } from '$env/static/private';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const endpoint = API_URL ? `${API_URL}/api/settings` : '/api/settings';

	try {
		const res = await fetch(endpoint, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (res.ok) {
			const data = await res.json();
			return { cyclingSpeed: data.cyclingSpeed ?? 15 };
		}

		const loadError = res.status === 401 ? 'You must be logged in' : 'Could not load settings';
		return { cyclingSpeed: 15, loadError };
	} catch (err) {
		console.error('Failed to load settings', err);
		return { cyclingSpeed: 15, loadError: 'Could not load settings' };
	}
};
