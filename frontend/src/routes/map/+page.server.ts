// +page.server.ts
import type { PageServerLoad } from './$types';
import { getUsernameFromToken } from '$lib/utils/auth';

export const load: PageServerLoad = async ({ cookies }) => {
	const token: string | undefined = cookies.get('session');
	const username: string | null = token ? getUsernameFromToken(token) : null;

	return {
		username
	};
};
