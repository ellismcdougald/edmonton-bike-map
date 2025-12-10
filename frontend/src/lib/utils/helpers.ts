/**
 * Capitalizes the first character of the given string.
 *
 * If `str` is empty or otherwise falsy, it is returned unchanged.
 *
 * @param str - The string to transform; empty values are returned as-is
 * @returns The input string with its first character converted to upper case
 */
export function capitalizeFirstLetter(str: string): string {
	if (!str) return str;
	return str.charAt(0).toUpperCase() + str.slice(1);
}