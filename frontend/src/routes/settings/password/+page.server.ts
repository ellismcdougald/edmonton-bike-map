import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { changePassword } from '$lib/api/user';

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

		// Basic validation stays here
		if (!currentPassword || !newPassword) {
			return fail(400, { error: 'Current password and new password are required' });
		}

		if (newPassword !== confirmPassword) {
			return fail(400, { error: 'New passwords do not match' });
		}

		try {
			await changePassword(fetch, currentPassword, newPassword);
			return { success: true };
		} catch (err) {
			return fail(400, {
				error: err instanceof Error ? err.message : 'Failed to change password'
			});
		}
	}
};
