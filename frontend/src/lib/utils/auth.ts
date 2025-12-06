export function getUsernameFromToken(token: string): string | null {
	try {
		const payload = JSON.parse(atob(token.split('.')[1])) as { username?: string };
		const username = payload.username;
		if (typeof username === 'string') {
			return username;
		} else {
			return null;
		}
	} catch {
		return null;
	}
}

export function getUserIdFromToken(): number | null {
	const token = localStorage.getItem('token');
	if (!token) return null;

	try {
		const payload = JSON.parse(atob(token.split('.')[1])) as { userId?: number };
		const userId = payload.userId;
		if (typeof userId === 'number') {
			return userId;
		} else {
			return null;
		}
	} catch {
		return null;
	}
}

export function removeToken(): void {
	localStorage.removeItem('token');
}
