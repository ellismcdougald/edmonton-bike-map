import type { RequestHandler } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

export const GET: RequestHandler = async ({ fetch, url }) => {
	const qs = url.searchParams.toString();
	const query = qs ? `?${qs}` : '';

	try {
		const res = await fetch(`${API_URL}/api/route${query}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (res.status === 401) {
			return new Response('Unauthorized', { status: 401 });
		}

		if (!res.ok) {
			return new Response('Upstream error', { status: 500 });
		}

		return new Response(await res.text(), { status: 200 });
	} catch (err) {
		console.error('Proxy route error:', err);
		return new Response('Network or server error', { status: 500 });
	}
};
