import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { fetchSettings } from '$lib/api/server/settings';

export const load: PageServerLoad = async ({ fetch, locals }) => {
	if (!locals.token) {
		throw redirect(303, '/login');
	}

	try {
		const { cyclingSpeed } = await fetchSettings(fetch);
		return { cyclingSpeed: cyclingSpeed ?? 15 };
	} catch (err) {
		if (err instanceof Error && err.message === 'unauthorized') {
			throw redirect(303, '/login');
		}

		return {
			cyclingSpeed: 15,
			loadError:
				err instanceof Error && err.message === 'unauthorized'
					? 'You must be logged in'
					: 'Could not load settings'
		};
	}
};
