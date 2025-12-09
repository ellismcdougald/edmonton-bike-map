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
