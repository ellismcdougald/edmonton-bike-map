<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type { WayFeature } from '$lib/types';
	import type { MapModeState } from '$lib/map/mapModes';
	import { toggleSelectStart, toggleSelectEnd } from '$lib/map/mapModes';
	import { wayState } from '$lib/state.svelte';
	import { loadLeaflet } from '$lib/map/LeafletMap';
	import { findRoute } from '$lib/map/mapActions';

	let mapInstance: InstanceType<typeof import('$lib/map/LeafletMap').LeafletMap> | null = null;
	const allWaysEndpoint = 'http://localhost:8080/api/all-ways';
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
			onclick={handleSelectStartClick}>
			Select Start Location
		</button>

		<button
			type="button"
			id="selectEndButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow"
			class:active={mode.selectEndActive}
			onclick={handleSelectEndClick}>
			Select End Location
		</button>

		<button
			type="button"
			id="findRouteButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow"
			onclick={() => findRoute({ mapInstance })}>
			Find Route
		</button>

		<button
			type="button"
			id="resetButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow">
			Reset
		</button>
	</div>
</div>
