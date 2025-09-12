<!--
  Map.svelte

  Purpose:
  Main interactive map component. Provides UI for selecting start/end points
  and triggering route-finding via the backend.

  Props:
  - none (map is self-contained, but interacts with global state)

  State:
  - mode (MapModeState): tracks whether user is selecting start or end point
  - mapInstance (LeafletMap | null): holds Leaflet map object once initialized

  Behavior:
  - Initializes Leaflet map on mount with OpenStreetMap tiles
  - Loads all ways from backend (`/api/all-ways`) into an info layer
  - On map click:
    - If selectStartActive → places start marker, toggles mode, shows info layer
    - If selectEndActive → places end marker, toggles mode, shows info layer
  - Control buttons:
    - "Select Start Location": toggles start selection mode
    - "Select End Location": toggles end selection mode
    - "Find Route": calls findRoute() with current mapInstance
    - "Reset": (placeholder — reset behavior not yet implemented)

  Notes:
  - Depends on $lib/map/loadLeaflet (LeafletMap wrapper)
  - Depends on $lib/map/mapModes for mode toggling
  - Depends on $lib/map/mapActions for route finding
  - Updates global wayState.selectedWay when a way is clicked
-->

<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type { WayFeature } from '$lib/types';
	import type { MapModeState } from '$lib/map/mapModes';
	import { toggleSelectStart, toggleSelectEnd } from '$lib/map/mapModes';
	import { wayState } from '$lib/state.svelte';
	import { loadLeaflet } from '$lib/map/loadLeaflet';
	import { findRoute } from '$lib/map/mapActions';

  const apiUrl = import.meta.env.VITE_API_URL;

	let mapInstance: InstanceType<typeof import('$lib/map/LeafletMap').LeafletMap> | null = null;
	const allWaysEndpoint = `${apiUrl}/api/all-ways`;
	let mode: MapModeState = { selectStartActive: false, selectEndActive: false };

	function onMapClick(latlng: [number, number]): void {
		const [lat, lng] = latlng;
		if (!mapInstance) return;

		if (mode.selectStartActive) {
			mapInstance.removeStartMarker();
			mapInstance.addStartMarker([lat, lng]);
			mode = toggleSelectStart(mode);
			mapInstance.showInfoLayer();
		} else if (mode.selectEndActive) {
			mapInstance.removeEndMarker();
			mapInstance.addEndMarker([lat, lng]);
			mode = toggleSelectEnd(mode);
			mapInstance.showInfoLayer();
		}
	}

	function handleSelectStartClick() {
		mode = toggleSelectStart(mode);
		mapInstance?.[mode.selectStartActive ? 'hideInfoLayer' : 'showInfoLayer']();
	}

	function handleSelectEndClick() {
		mode = toggleSelectEnd(mode);
		mapInstance?.[mode.selectEndActive ? 'hideInfoLayer' : 'showInfoLayer']();
	}

	function resetMap() {
		if (mapInstance) {
			mapInstance.reset();
			mode = { selectStartActive: false, selectEndActive: false };
		}
	}

	onMount(async () => {
		const { LeafletMap } = await loadLeaflet();
		mapInstance = new LeafletMap();

		mapInstance.addTileLayer(
			'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
			19,
			10,
			'© OpenStreetMap contributors'
		);

		mapInstance.onMapClick(onMapClick);

		try {
			const res = await fetch(allWaysEndpoint);
			if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
			const geojson = await res.json();
			mapInstance.loadInfoLayer(
				geojson,
				{ color: 'red', weight: 1, opacity: 0 },
				(way: WayFeature) => {
					way.id = Number(way.id);
					wayState.selectedWay = way;
				}
			);
		} catch (err) {
			console.error('Error loading info layer:', err);
		}
	});

	onDestroy(() => {
		if (mapInstance?.map) {
			mapInstance.map.remove();
		}
	});
</script>

<div id="map-content" class="w-full h-full">
	<div id="map" class="w-full h-9/10"></div>
	<div id="controls" class="h-1/10 flex gap-4 p-2 w-6/10">
		<button
			type="button"
			id="selectStartButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow"
			class:active={mode.selectStartActive}
			onclick={handleSelectStartClick}
		>
			Select Start Location
		</button>

		<button
			type="button"
			id="selectEndButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow"
			class:active={mode.selectEndActive}
			onclick={handleSelectEndClick}
		>
			Select End Location
		</button>

		<button
			type="button"
			id="findRouteButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow"
			onclick={() => findRoute({ mapInstance })}
		>
			Find Route
		</button>

		<button
			type="button"
			id="resetButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow"
			onclick={() => resetMap()}
		>
			Reset
		</button>
	</div>
</div>
