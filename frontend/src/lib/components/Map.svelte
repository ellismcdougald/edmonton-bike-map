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
	import { goto } from '$app/navigation';
	import type { WayFeature, WayFeatureCollection } from '$lib/types';
	import type { MapModeState } from '$lib/map/mapModes';
	import { toggleSelectStart, toggleSelectEnd } from '$lib/map/mapModes';
	import { wayState } from '$lib/state.svelte';
	import { loadLeaflet } from '$lib/map/loadLeaflet';
	import { findRoute } from '$lib/map/mapActions';

	let mapInstance: InstanceType<typeof import('$lib/map/LeafletMap').LeafletMap> | null = null;
	let mode: MapModeState = $state({ selectStartActive: false, selectEndActive: false });
	let loadError: string | null = $state(null);

	let { ways }: { ways: Promise<GeoJSON.GeoJsonObject> } = $props();

	// URL Helpers:
	function updateUrlWithWay(id: number) {
		try {
			const path = `${window.location.pathname}?way=${encodeURIComponent(String(id))}`;
			goto(path, { replaceState: true, noScroll: true });
		} catch {
			/* ignore */
		}
	}

	// Map setup:
	async function initializeMap() {
		await createMap();
		try {
			const waysData = await ways;
			mapInstance?.loadInfoLayer(waysData, { color: 'red', weight: 1, opacity: 0 }, handleWayClick);
			restoreSelectionFromUrl();
		} catch (err) {
			const message =
				err instanceof Error && err.message ? err.message : 'Unable to load ways data';
			loadError = message;
			console.error('Error loading ways:', err);
		}
	}

	async function createMap() {
		const { LeafletMap } = await loadLeaflet();
		mapInstance = new LeafletMap();
		mapInstance.addTileLayer(
			'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
			19,
			10,
			'© OpenStreetMap contributors'
		);
		mapInstance.onMapClick(onMapClick);
	}

	async function restoreSelectionFromUrl() {
		if (typeof window === 'undefined' || !window.location?.search) return;

		const match = window.location.search.match(/[?&]way=([^&]+)/);
		const targetId = match?.[1] ? Number(decodeURIComponent(match[1])) : null;
		if (!targetId) return;

		try {
			const waysData = await ways;
			if (waysData?.type === 'FeatureCollection') {
				const featureCollection = waysData as WayFeatureCollection;
				const feature = featureCollection.features.find((f) => {
					const id = f.id ?? f.properties?.id;
					return Number(id) === targetId;
				});

				if (feature) {
					wayState.selectedWay = {
						id: Number(feature.id ?? feature.properties.id),
						tags: Object.fromEntries(
							Object.entries(feature.properties || {}).map(([k, v]) => [k, String(v)])
						)
					};
				}
			}
		} catch (err) {
			console.error('Error restoring selection:', err);
		}
	}

	// Map interactions:
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

	async function handleFindRouteClick() {
		if (!mapInstance) return;
		try {
			await findRoute({ mapInstance });
		} catch (err) {
			if (err instanceof Error && err.message === 'Unauthorized') {
				await fetch('/logout', { method: 'POST' });
				goto('/login');
			}
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

	function handleWayClick(way: WayFeature) {
		way.id = Number(way.id);
		wayState.selectedWay = way;
		updateUrlWithWay(way.id);
	}

	onMount(async () => {
		try {
			await initializeMap();
		} catch (err) {
			console.error('Error initializing map:', err);
		}
	});

	onDestroy(() => {
		mapInstance?.map?.remove();
	});
</script>

<div id="map-content" class="w-full h-full">
	<div id="map" class="w-full h-9/10 relative">
		{#if loadError}
			<div
				role="alert"
				class="absolute top-2 left-2 z-20 bg-red-600 text-white px-3 py-2 rounded shadow"
			>
				{loadError}
			</div>
		{/if}
	</div>
	<div id="controls" class="h-1/10 flex gap-4 p-2 w-6/10">
		<button
			type="button"
			id="selectStartButton"
			class="text-white py-2 px-4 rounded transition flex-grow bg-blue-600 hover:bg-blue-700"
			class:bg-blue-800={mode.selectStartActive}
			onclick={handleSelectStartClick}
			class:active={mode.selectStartActive}
		>
			Select Start Location
		</button>

		<button
			type="button"
			id="selectEndButton"
			class="text-white py-2 px-4 rounded transition flex-grow bg-blue-600 hover:bg-blue-700"
			class:bg-blue-800={mode.selectEndActive}
			onclick={handleSelectEndClick}
			class:active={mode.selectEndActive}
		>
			Select End Location
		</button>

		<button
			type="button"
			id="findRouteButton"
			class="text-white py-2 px-4 rounded transition flex-grow bg-blue-600 hover:bg-blue-700"
			onclick={handleFindRouteClick}
		>
			Find Route
		</button>

		<button
			type="button"
			id="resetButton"
			class="text-white py-2 px-4 rounded transition flex-grow bg-blue-600 hover:bg-blue-700"
			onclick={resetMap}
		>
			Reset
		</button>
	</div>

	<style>
		.distance-control,
		.time-control {
			background-color: rgba(255, 255, 255, 0.9);
			padding: 0.5rem 1rem;
			border-radius: 0.5rem;
			box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
			font-size: 0.875rem;
			font-weight: 500;
			text-align: center;
		}

		button.active {
			box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.5);
			transform: scale(1.05);
		}
	</style>
</div>
