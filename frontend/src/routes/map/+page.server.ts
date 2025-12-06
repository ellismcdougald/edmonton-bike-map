// +page.server.ts
import type { PageServerLoad } from './$types';
import { getUsernameFromToken } from '$lib/utils/auth';

const API_URL = process.env.API_URL;
const ALL_WAYS_ENDPOINT = `${API_URL}/api/all-ways`;

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	// Decode username
	const token: string | undefined = cookies.get('session');
	const username: string | null = token ? getUsernameFromToken(token) : null;

	// Get ways
	const waysResponse = await fetch(ALL_WAYS_ENDPOINT);
	if (!waysResponse.ok) throw new Error(`HTTP error! status: ${waysResponse.status}`);
	const ways = await waysResponse.json();
	return { username, ways };
};
