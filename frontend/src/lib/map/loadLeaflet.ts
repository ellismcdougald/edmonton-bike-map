export async function loadLeaflet() {
	const L = await import('leaflet'); // browser-only import
	const { LeafletMap } = await import('./LeafletMap'); // your class

	// Return a new class that injects L in the constructor
	return {
		LeafletMap: class extends LeafletMap {
			constructor() {
				super(L); // inject Leaflet
			}
		}
	};
}
