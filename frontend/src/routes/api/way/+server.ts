import type { RequestHandler } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

export const GET: RequestHandler = async ({ url, fetch }) => {
	try {
		const id = url.searchParams.get('id');

		if (!id) {
			return new Response('Missing id query parameter', { status: 400 });
		}

		const res = await fetch(`${API_URL}/api/way?id=${id}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (res.status === 401) {
			return new Response('Unauthorized', { status: 401 });
		}

		if (res.status === 404) {
			return new Response('Way not found', { status: 404 });
		}

		if (!res.ok) {
			return new Response('Upstream error', { status: 500 });
		}

		return new Response(await res.text(), { status: 200 });
	} catch (err) {
		console.error('Proxy way error:', err);
		return new Response('Network or server error', { status: 500 });
	}
};
