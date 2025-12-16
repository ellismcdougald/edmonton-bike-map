/**
 * mapActions.ts
 *
 * Purpose:
 * Fetches a route between the start and end markers on a LeafletMap instance
 * and displays it on the map.
 *
 * Types:
 * - FindRouteOptions: object with `mapInstance` (LeafletMap | null)
 *
 * Behavior:
 * - Validates that both start and end markers exist; alerts user if missing
 * - Builds URLSearchParams for backend route API
 * - Calls fetchRoute() to get GeoJSON from backend
 * - Calls applyRouteToMap() to remove any existing route and render new route layer
 * - Catches and alerts errors during fetch or display
 *
 * Notes:
 * - Uses browser global `fetch` to communicate with the backend
 * - Depends on LeafletMap class for marker retrieval and map updates
 * - Alerts are used for user feedback in case of missing points or errors
 */

import type { LeafletMap } from './LeafletMap';
import type { FeatureCollection } from 'geojson';

interface FindRouteOptions {
	mapInstance: LeafletMap | null;
}

/**
 * Fetches a route between the map's selected start and end points and renders it on the provided LeafletMap.
 *
 * If `mapInstance` is null the function does nothing. If either start or end point is missing an alert is shown and no request is made. A thrown `Error` with message `"Unauthorized"` is rethrown; other errors cause an alert to the user and are logged to the console.
 *
 * @param mapInstance - The LeafletMap to read start/end points from and to render the route; may be null
 */
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

/**
 * Create URLSearchParams containing start and end coordinates for a route query.
 *
 * @param startLatLng - Tuple [latitude, longitude] of the start point
 * @param endLatLng - Tuple [latitude, longitude] of the end point
 * @returns URLSearchParams with keys `startLatitude`, `startLongitude`, `endLatitude`, and `endLongitude` whose values are the corresponding coordinates as strings
 */
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

/**
 * Fetches a route GeoJSON from the backend using the provided query parameters.
 *
 * @param params - URLSearchParams containing route query keys (e.g., `startLatitude`, `startLongitude`, `endLatitude`, `endLongitude`)
 * @returns The GeoJSON FeatureCollection describing the route
 * @throws `Error('Unauthorized')` if the server responds with HTTP 401
 * @throws `Error('Failed to get route data')` if the response is not OK (non-2xx and not 401)
 * @throws Any underlying network or parsing error encountered during fetch
 */
async function fetchRoute(params: URLSearchParams): Promise<FeatureCollection> {
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

/**
 * Replace the current route layer on the given map with the provided GeoJSON route.
 *
 * @param mapInstance - The LeafletMap instance to update
 * @param geojson - A GeoJSON FeatureCollection representing the route to render
 */
function applyRouteToMap(mapInstance: LeafletMap, geojson: FeatureCollection) {
	mapInstance.removeRouteLayer();
	mapInstance.loadRouteLayer(geojson);
}

interface FindRoutesOptions {
	mapInstance: LeafletMap | null;
	k?: number;
}

/**
 * Fetches multiple routes between the map's selected start and end points and renders them on the provided LeafletMap.
 *
 * @param mapInstance - The LeafletMap to read start/end points from and to render the routes
 * @param k - Number of alternative routes to find (default: 3)
 */
export async function findRoutes({ mapInstance, k = 3 }: FindRoutesOptions) {
	if (!mapInstance) return;

	const startLatLng = mapInstance.getStartLatLng();
	const endLatLng = mapInstance.getEndLatLng();

	if (!startLatLng || !endLatLng) {
		alert('Make sure you have selected both a start and an end point!');
		return;
	}

	const params = buildRouteParams(startLatLng, endLatLng);
	params.set('k', k.toString());

	try {
		const geojson = await fetchRoutes(params);
		applyRouteToMap(mapInstance, geojson);
	} catch (err: unknown) {
		if (err instanceof Error) {
			if (err.message === 'Unauthorized') throw err;

			alert('Error fetching or displaying routes: ' + err.message);
			console.error(err);
		} else {
			alert('Error fetching or displaying routes: ' + String(err));
			console.error(err);
		}
	}
}

/**
 * Fetches multiple routes GeoJSON from the backend using the provided query parameters.
 *
 * @param params - URLSearchParams containing route query keys including k parameter
 * @returns The GeoJSON FeatureCollection describing multiple routes
 */
async function fetchRoutes(params: URLSearchParams): Promise<FeatureCollection> {
	const query = params.toString();
	const url = query ? `/api/routes?${query}` : '/api/routes';

	try {
		const res = await fetch(url);

		if (res.status === 401) throw new Error('Unauthorized');
		if (!res.ok) throw new Error('Failed to get routes data');

		return await res.json();
	} catch (err) {
		console.error('Error fetching routes:', err);
		throw err;
	}
}
