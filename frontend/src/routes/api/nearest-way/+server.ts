import type { RequestHandler } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

export const GET: RequestHandler = async ({ url, fetch }) => {
	try {
		const lat = url.searchParams.get('lat');
		const lng = url.searchParams.get('lng');

		if (!lat || !lng) {
			return new Response('Missing lat or lng query parameters', { status: 400 });
		}

		const latNum = parseFloat(lat);
		const lngNum = parseFloat(lng);
		const isValidLat = Number.isFinite(latNum) && latNum >= -90 && latNum <= 90;
		const isValidLng = Number.isFinite(lngNum) && lngNum >= -180 && lngNum <= 180;

		if (!isValidLat || !isValidLng) {
			return new Response('Invalid numeric lat or lng', { status: 400 });
		}

		const params = new URLSearchParams({ lat: latNum.toString(), lng: lngNum.toString() });
		const res = await fetch(`${API_URL}/api/nearest-way?${params.toString()}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (res.status === 401) {
			return new Response('Unauthorized', { status: 401 });
		}

		if (res.status === 404) {
			return new Response('No ways found nearby', { status: 404 });
		}

		if (!res.ok) {
			return new Response('Upstream error', { status: 500 });
		}

		return new Response(await res.text(), { status: 200 });
	} catch (err) {
		console.error('Proxy nearest-way error:', err);
		return new Response('Network or server error', { status: 500 });
	}
};
