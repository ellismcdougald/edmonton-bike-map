// +page.server.ts
import type { PageServerLoad } from './$types';
import { getUsernameFromToken } from '$lib/utils/auth';
import { redirect } from '@sveltejs/kit';

const API_URL = process.env.API_URL;
const ALL_WAYS_ENDPOINT = `${API_URL}/api/all-ways`;

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	// Decode username
	const token: string | undefined = cookies.get('session');
	if (!token) {
		throw redirect(302, '/login');
	}
	const username: string | null = token ? getUsernameFromToken(token) : null;

	// Get ways
	const waysResponse = await fetch(ALL_WAYS_ENDPOINT);
	if (!waysResponse.ok) {
		if (waysResponse.status === 401) {
			throw redirect(302, '/login');
		}
		throw new Error(`HTTP error! status: ${waysResponse.status}`);
	}

	const ways = await waysResponse.json();
	return { username, ways };
};
