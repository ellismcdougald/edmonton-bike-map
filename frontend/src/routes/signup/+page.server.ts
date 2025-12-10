import { fail, redirect } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

export const actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const username = data.get('username');
		const password = data.get('password');

		if (!username || !password) {
			return fail(400, { error: 'Username and password are required' });
		}

		const res = await fetch(`${API_URL}/api/signup`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ username, password })
		});

		if (!res.ok) {
			const errorText = await res.text();
			return fail(400, { error: errorText || 'Signup failed' });
		}

		throw redirect(303, '/login');
	}
};
