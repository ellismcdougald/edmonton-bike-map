import type { PageServerLoad } from './$types';
import { fetchSettings } from '$lib/api/server/settings';

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		const { cyclingSpeed } = await fetchSettings(fetch);
		return { cyclingSpeed: cyclingSpeed ?? 15 };
	} catch (err) {
		return {
			cyclingSpeed: 15,
			loadError:
				err instanceof Error && err.message === 'unauthorized'
					? 'You must be logged in'
					: 'Could not load settings'
		};
	}
};
