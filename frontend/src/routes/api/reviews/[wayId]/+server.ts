import type { RequestHandler } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

interface ReviewResponse {
	WayID: string;
	UserID: string;
	Rating: number;
	Comment?: string;
	CreatedAt: string;
	Username: string;
}

export const GET: RequestHandler = async ({ params, fetch }) => {
	const wayId = Number(params.wayId);
	if (isNaN(wayId)) {
		return new Response('Invalid Way ID', { status: 400 });
	}

	try {
		const res = await fetch(`${API_URL}/api/reviews?wayID=${wayId}`, {
			method: 'GET',
			headers: { 'Content-Type': 'application/json' }
		});

		if (!res.ok) {
			const errorText = await res.text();
			return new Response(errorText || 'Failed to fetch reviews', { status: res.status });
		}

		const data: ReviewResponse[] = await res.json();
		const mappedData = data.map((r: ReviewResponse) => ({
			wayId: Number(r.WayID),
			rating: r.Rating,
			comment: r.Comment ?? '',
			createdAt: r.CreatedAt,
			username: r.Username
		}));

		return new Response(JSON.stringify(mappedData), {
			headers: { 'Content-Type': 'application/json' }
		});
	} catch (err) {
		return new Response(`Server error: ${err}`, { status: 500 });
	}
};

export const DELETE: RequestHandler = async ({ params, locals, fetch }) => {
	const wayId = Number(params.wayId);
	if (isNaN(wayId)) {
		return new Response('Invalid Way ID', { status: 400 });
	}

	const token = locals.token;
	const headers: Record<string, string> = {};
	if (token) {
		headers.Authorization = `Bearer ${token}`;
	}

	const res = await fetch(`${API_URL}/api/reviews/${wayId}`, {
		method: 'DELETE',
		headers
	});

	if (!res.ok) {
		const text = await res.text();
		return new Response(text || 'Failed to delete review', { status: res.status });
	}

	return new Response(null, { status: 204 });
};
