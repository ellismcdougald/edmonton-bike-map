/**
 * Extracts the `username` field from the payload of a JWT-formatted token.
 *
 * @param token - JWT-like token whose payload is the middle (base64-encoded) segment
 * @returns The `username` string if present and a string, `null` otherwise.
 */
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
