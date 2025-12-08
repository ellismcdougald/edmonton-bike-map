import type { RequestHandler } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

export const POST: RequestHandler = async ({ request, fetch }) => {
	try {
		const reviewData = await request.json();

		const res = await fetch(`${API_URL}/api/reviews`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(reviewData)
		});

		if (!res.ok) {
			const errorText = await res.text();
			return new Response(errorText || 'Failed to submit review', { status: res.status });
		}

		return new Response(null, { status: 201 });
	} catch (err) {
		return new Response(`Server error: ${err}`, { status: 500 });
	}
};
