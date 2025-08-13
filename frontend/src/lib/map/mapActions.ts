import type { LeafletMap } from './LeafletMap';
import type { FeatureCollection } from 'geojson';

type FetchFn = typeof fetch;

interface FindRouteOptions {
	mapInstance: LeafletMap | null;
	fetchFn?: FetchFn; // fetch function to be used by findRoute
}

export async function findRoute({ mapInstance, fetchFn = fetch }: FindRouteOptions) {
	if (!mapInstance) return;

	const startLatLng = mapInstance.getStartLatLng();
	const endLatLng = mapInstance.getEndLatLng();

	if (!startLatLng || !endLatLng) {
		alert('Make sure you have selected both a start and an end point!');
		return;
	}

	const params = buildRouteParams(startLatLng, endLatLng);

	try {
		const geojson = await fetchRoute(params, fetchFn);
		applyRouteToMap(mapInstance, geojson);
	} catch (err: unknown) {
		if (err instanceof Error) {
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

async function fetchRoute(params: URLSearchParams, fetchFn: FetchFn) {
	const res = await fetchFn(`http://localhost:8080/api/route?${params.toString()}`);
	if (!res.ok) throw new Error('Failed to get route data');
	return res.json();
}

function applyRouteToMap(mapInstance: LeafletMap, geojson: FeatureCollection) {
	mapInstance.removeRouteLayer();
	mapInstance.loadRouteLayer(geojson);
}
