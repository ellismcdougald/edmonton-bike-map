/**
 * loadLeaflet.ts
 *
 * Purpose:
 * Dynamically imports Leaflet (browser-only) and the LeafletMap wrapper class,
 * then returns a version of LeafletMap with Leaflet injected into the constructor.
 *
 * Behavior:
 * - Uses dynamic import() to load 'leaflet' only in the browser
 * - Imports LeafletMap class
 * - Returns an object with LeafletMap class extended to automatically receive Leaflet
 *
 * Usage:
 * const { LeafletMap } = await loadLeaflet();
 * const mapInstance = new LeafletMap();
 *
 * Notes:
 * - This avoids server-side rendering errors in SvelteKit
 * - LeafletMap instances returned are ready to use with Leaflet already injected
 */

export async function loadLeaflet() {
	const L = await import('leaflet');
	const { LeafletMap } = await import('./LeafletMap');

	return {
		LeafletMap: class extends LeafletMap {
			constructor() {
				super(L); // inject Leaflet
			}
		}
	};
}
