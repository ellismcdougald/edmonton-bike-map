import { API_URL } from '$env/static/private';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

const endpoint = API_URL ? `${API_URL}/api/change-password` : '/api/change-password';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.token) {
		throw redirect(303, '/login');
	}
	return {};
};

export const actions: Actions = {
	default: async ({ request, fetch, locals }) => {
		if (!locals.token) {
			throw redirect(303, '/login');
		}

		const formData = await request.formData();
		const currentPassword = String(formData.get('currentPassword') ?? '');
		const newPassword = String(formData.get('newPassword') ?? '');
		const confirmPassword = String(formData.get('confirmPassword') ?? '');

		if (!currentPassword || !newPassword) {
			return fail(400, { error: 'Current password and new password are required' });
		}

		if (newPassword !== confirmPassword) {
			return fail(400, { error: 'New passwords do not match' });
		}

		const res = await fetch(endpoint, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ currentPassword, newPassword })
		});

		if (res.ok) {
			return { success: true };
		}

		const text = await res.text();
		return fail(res.status, { error: text || 'Failed to change password' });
	}
};
