import type { FeatureCollection } from '$lib/types';

/**
 * Retrieve adjacent ways for a specific way.
 * Adjacent ways are those that share at least one node with the specified way.
 *
 * @param wayId - ID of the way whose adjacent ways should be fetched
 * @returns A GeoJSON FeatureCollection containing adjacent way features
 * @throws Error - When the server responds with a non-OK status
 */
export async function getAdjacentWays(wayId: number): Promise<FeatureCollection> {
	const res = await fetch(`/api/adjacent-ways?id=${wayId}`);

	if (!res.ok) {
		throw new Error(`Failed to fetch adjacent ways: ${res.statusText}`);
	}

	return res.json();
}
