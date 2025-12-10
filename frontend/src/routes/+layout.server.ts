import { redirect } from '@sveltejs/kit';

/**
 * Redirects requests for the root path to "/map".
 *
 * @param url - The request URL object whose `pathname` is inspected
 * @throws Throws a redirect with status `302` to `"/map"` when `url.pathname` is `"/"`
 */
export function load({ url }) {
	if (url.pathname === '/') {
		throw redirect(302, '/map');
	}
}