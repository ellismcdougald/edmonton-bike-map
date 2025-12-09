import { goto } from '$app/navigation';

export async function createReview(wayId: number, rating: number, comment: string | null) {
	const res = await fetch('/api/reviews', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ wayId, rating, comment })
	});

	if (!res.ok) {
		if (res.status === 401) {
			await fetch('/logout');
			goto('/login');
		}
		const text = await res.text();
		throw new Error(text || 'Failed to submit review');
	}

	return;
}
