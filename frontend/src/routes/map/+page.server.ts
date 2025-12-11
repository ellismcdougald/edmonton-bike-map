// +page.server.ts
import type { PageServerLoad } from './$types';
import { getUsernameFromToken } from '$lib/utils/auth';
import { fetchAllWays } from '$lib/api/server/ways';

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	const token: string | undefined = cookies.get('session');
	const username: string | null = token ? getUsernameFromToken(token) : null;

	// Return ways as a promise - this allows streaming
	return {
		username,
		ways: fetchAllWays(fetch)
	};
};
