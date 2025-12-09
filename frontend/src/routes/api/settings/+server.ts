import type { RequestHandler } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

export const GET: RequestHandler = async ({ fetch }) => {
	try {
		const res = await fetch(`${API_URL}/api/settings`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		return new Response(await res.text(), { status: res.status });
	} catch (err) {
		return new Response(`Server error: ${err}`, { status: 500 });
	}
};

export const POST: RequestHandler = async ({ request, fetch }) => {
	console.log('called');
	try {
		const body = await request.text();
		const res = await fetch(`${API_URL}/api/settings`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body
		});

		if (!res.ok) {
			const errorText = await res.text();
			return new Response(errorText || 'Failed to save settings', { status: res.status });
		}

		return new Response(null, { status: 200 });
	} catch (err) {
		return new Response(`Server error: ${err}`, { status: 500 });
	}
};
