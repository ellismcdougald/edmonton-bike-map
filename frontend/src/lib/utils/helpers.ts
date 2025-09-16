/**
 * utils/helpers.ts
 *
 * Provides general helper functions.
 *
 * Functions:
 * - capitalizeFirstLetter(str): capitalizes the first letter of the given string
 *   - Returns a new string
 */

export function capitalizeFirstLetter(str: string): string {
	if (!str) return str;
	return str.charAt(0).toUpperCase() + str.slice(1);
}
