import { redirect } from '@sveltejs/kit';

const API_URL = process.env.API_URL;

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
