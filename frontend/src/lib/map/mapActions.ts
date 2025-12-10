/**
 * mapActions.ts
 *
 * Purpose:
 * Fetches a route between the start and end markers on a LeafletMap instance
 * and displays it on the map.
 *
 * Types:
 * - FindRouteOptions: object with `mapInstance` (LeafletMap | null) and optional `fetchFn`
 * - FetchFn: type alias for a fetch-like function
 *
 * Behavior:
 * - Validates that both start and end markers exist; alerts user if missing
 * - Builds URLSearchParams for backend route API
 * - Calls fetchRoute() to get GeoJSON from backend
 * - Calls applyRouteToMap() to remove any existing route and render new route layer
 * - Catches and alerts errors during fetch or display
 *
 * Notes:
 * - Default fetch function is browser global `fetch`, but can be overridden (useful for testing)
 * - Depends on LeafletMap class for marker retrieval and map updates
 * - Alerts are used for user feedback in case of missing points or errors
 */

import type { LeafletMap } from './LeafletMap';
import type { FeatureCollection } from 'geojson';

type FetchFn = typeof fetch;

interface FindRouteOptions {
	mapInstance: LeafletMap | null;
	fetchFn?: FetchFn; // fetch function to be used by findRoute
}

export async function findRoute({ mapInstance }: FindRouteOptions) {
	if (!mapInstance) return;

	const startLatLng = mapInstance.getStartLatLng();
	const endLatLng = mapInstance.getEndLatLng();

	if (!startLatLng || !endLatLng) {
		alert('Make sure you have selected both a start and an end point!');
		return;
	}

	const params = buildRouteParams(startLatLng, endLatLng);

	try {
		const geojson = await fetchRoute(params);
		applyRouteToMap(mapInstance, geojson);
	} catch (err: unknown) {
		if (err instanceof Error) {
			if (err.message === 'Unauthorized') throw err;

			alert('Error fetching or displaying route: ' + err.message);
			console.error(err);
		} else {
			alert('Error fetching or displaying route: ' + String(err));
			console.error(err);
		}
	}
}

function buildRouteParams(
	startLatLng: [number, number],
	endLatLng: [number, number]
): URLSearchParams {
	return new URLSearchParams({
		startLatitude: startLatLng[0].toString(),
		startLongitude: startLatLng[1].toString(),
		endLatitude: endLatLng[0].toString(),
		endLongitude: endLatLng[1].toString()
	});
}

async function fetchRoute(params: URLSearchParams) {
	const query = params.toString();
	const url = query ? `/api/route?${query}` : '/api/route';

	try {
		const res = await fetch(url);

		if (res.status === 401) throw new Error('Unauthorized');
		if (!res.ok) throw new Error('Failed to get route data');

		return await res.json();
	} catch (err) {
		console.error('Error fetching route:', err);
		throw err;
	}
}

function applyRouteToMap(mapInstance: LeafletMap, geojson: FeatureCollection) {
	mapInstance.removeRouteLayer();
	mapInstance.loadRouteLayer(geojson);
}
