// +page.server.ts
import type { PageServerLoad } from './$types';

const API_URL = process.env.API_URL;
const ALL_WAYS_ENDPOINT = `${API_URL}/api/all-ways`;

export const load: PageServerLoad = async ({ fetch }) => {
	const waysResponse = await fetch(ALL_WAYS_ENDPOINT);
	if (!waysResponse.ok) throw new Error(`HTTP error! status: ${waysResponse.status}`);
	const ways = await waysResponse.json();
	return { ways };
};
