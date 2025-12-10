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

	// Return ways as a promise - this allows streaming
	return {
		username,
		ways: fetch(ALL_WAYS_ENDPOINT).then(async (res) => {
			if (res.status === 401) {
				throw redirect(302, '/login');
			}
			if (!res.ok) {
				throw new Error(`HTTP error! status: ${res.status}`);
			}
			return res.json();
		})
	};
};
