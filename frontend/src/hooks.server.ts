import type { Handle, HandleFetch } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('session');
	event.locals.token = token;
	return await resolve(event);
};

export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
	const token = event.locals.token;
	if (token) {
		request.headers.set('Authorization', `Bearer ${token}`);
	}
	return fetch(request);
};
