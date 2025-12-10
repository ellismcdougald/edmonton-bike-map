// +page.server.ts
import type { PageServerLoad } from './$types';
import { getUsernameFromToken } from '$lib/utils/auth';
import { fetchAllWays } from '$lib/api/server/ways';
import { redirect } from '@sveltejs/kit';

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
		ways: fetchAllWays(fetch)
	};
};
