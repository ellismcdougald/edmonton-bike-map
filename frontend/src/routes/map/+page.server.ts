// +page.server.ts
import type { PageServerLoad } from './$types';

const API_URL = process.env.VITE_API_URL;
const ALL_WAYS_ENDPOINT = `${API_URL}/api/all-ways`;

export const load: PageServerLoad = async ({ fetch }) => {
	const res = await fetch(ALL_WAYS_ENDPOINT);
	if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
	const ways = await res.json();
	console.log('Fetched ways:', ways);
	return { ways };
};
