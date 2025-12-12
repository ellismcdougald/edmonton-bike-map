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
  - On map click (when not selecting start/end markers):
    - Queries the backend for the nearest way to the clicked coordinates
    - Updates the sidebar with the way information
  - On map click (when selecting markers):
    - If selectStartActive → places start marker, toggles mode
    - If selectEndActive → places end marker, toggles mode
  - Control buttons:
    - "Select Start Location": toggles start selection mode
    - "Select End Location": toggles end selection mode
    - "Find Route": calls findRoute() with current mapInstance
    - "Reset": clears markers and selections

  Notes:
  - Depends on $lib/map/loadLeaflet (LeafletMap wrapper)
  - Depends on $lib/map/mapModes for mode toggling
  - Depends on $lib/map/mapActions for route finding
  - Uses /api/nearest-way to find ways by coordinates
  - Updates global wayState.selectedWay when a way is clicked
-->

<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import type { MapModeState } from '$lib/map/mapModes';
	import { toggleSelectStart, toggleSelectEnd } from '$lib/map/mapModes';
	import { wayState } from '$lib/state.svelte';
	import { loadLeaflet } from '$lib/map/loadLeaflet';
	import { findRoute } from '$lib/map/mapActions';

	let mapInstance: InstanceType<typeof import('$lib/map/LeafletMap').LeafletMap> | null = null;
	let mode: MapModeState = $state({ selectStartActive: false, selectEndActive: false });
	let loadError: string | null = $state(null);

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
		await restoreSelectionFromUrl();
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

	// Restore way selection from URL parameter
	async function restoreSelectionFromUrl() {
		if (typeof window === 'undefined' || !window.location?.search) return;

		const match = window.location.search.match(/[?&]way=([^&]+)/);
		const targetId = match?.[1] ? Number(decodeURIComponent(match[1])) : null;
		if (!targetId) return;

		try {
			const response = await fetch(`/api/way?id=${targetId}`);
			if (!response.ok) {
				console.error('Failed to load way from URL');
				return;
			}

			const wayData = await response.json();
			wayState.selectedWay = {
				id: wayData.id,
				tags: wayData.tags
			};
		} catch (err) {
			console.error('Error restoring selection from URL:', err);
		}
	}

	// API call to get nearest way
	async function fetchNearestWay(lat: number, lng: number) {
		try {
			const response = await fetch(`/api/nearest-way?lat=${lat}&lng=${lng}`);
			if (!response.ok) {
				if (response.status === 404) {
					loadError = 'No ways found at this location';
				} else {
					loadError = 'Failed to find nearby way';
				}
				return null;
			}

			const data = await response.json();
			loadError = null;
			return data as { id: number; tags: Record<string, string> };
		} catch (err) {
			const message =
				err instanceof Error && err.message ? err.message : 'Error fetching nearest way';
			loadError = message;
			console.error('Error fetching nearest way:', err);
			return null;
		}
	}

	// Map interactions:
	async function onMapClick(latlng: [number, number]): Promise<void> {
		const [lat, lng] = latlng;
		if (!mapInstance) return;

		if (mode.selectStartActive) {
			mapInstance.removeStartMarker();
			mapInstance.addStartMarker([lat, lng]);
			mode = toggleSelectStart(mode);
		} else if (mode.selectEndActive) {
			mapInstance.removeEndMarker();
			mapInstance.addEndMarker([lat, lng]);
			mode = toggleSelectEnd(mode);
		} else {
			// Normal map click - find nearest way
			const wayData = await fetchNearestWay(lat, lng);
			if (wayData) {
				wayState.selectedWay = {
					id: wayData.id,
					tags: wayData.tags
				};
				updateUrlWithWay(wayData.id);
			}
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
	}

	function handleSelectEndClick() {
		mode = toggleSelectEnd(mode);
	}

	function resetMap() {
		if (mapInstance) {
			mapInstance.reset();
			mode = { selectStartActive: false, selectEndActive: false };
		}
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
