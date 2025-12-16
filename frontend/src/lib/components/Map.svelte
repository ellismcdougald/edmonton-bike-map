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
	import { findRoutes } from '$lib/map/mapActions';
	import type { WayFeatureGeoJSON } from '$lib/types';

	let mapInstance: InstanceType<typeof import('$lib/map/LeafletMap').LeafletMap> | null = null;
	let mode: MapModeState = $state({ selectStartActive: false, selectEndActive: false });
	let loadError: string | null = $state(null);
	let isLoadingWay: boolean = $state(false);
	let nearestWayAbortController: AbortController | null = null;
	let selectedWayAbortController: AbortController | null = null;

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

			const wayFeature = (await response.json()) as WayFeatureGeoJSON;
			if (typeof wayFeature?.properties?.id !== 'undefined') {
				const idNum = Number(wayFeature.properties.id);
				wayState.selectedWay = {
					id: idNum,
					tags: wayFeature.properties as Record<string, string>
				};
				await loadSelectedWayHighlight(idNum);
			}
		} catch (err) {
			console.error('Error restoring selection from URL:', err);
		}
	}

	async function fetchWayGeojson(
		id: number,
		signal?: AbortSignal
	): Promise<WayFeatureGeoJSON | null> {
		try {
			const res = await fetch(`/api/way?id=${id}`, { signal });
			if (!res.ok) {
				return null;
			}
			const data = (await res.json()) as WayFeatureGeoJSON;
			if (!data || data.type !== 'Feature' || !data.geometry) return null;
			return data;
		} catch (err) {
			if (err instanceof DOMException && err.name === 'AbortError') return null;
			console.error('Error fetching way GeoJSON:', err);
			return null;
		}
	}

	async function loadSelectedWayHighlight(id: number) {
		if (!mapInstance) return;

		if (selectedWayAbortController) {
			selectedWayAbortController.abort();
		}

		const controller = new AbortController();
		selectedWayAbortController = controller;
		const feature = await fetchWayGeojson(id, controller.signal);
		if (controller !== selectedWayAbortController) return;
		if (feature) {
			mapInstance.loadSelectedWayLayer(feature);
		}
	}

	// API call to get nearest way
	async function fetchNearestWay(lat: number, lng: number, signal?: AbortSignal) {
		isLoadingWay = true;
		try {
			const response = await fetch(`/api/nearest-way?lat=${lat}&lng=${lng}`.toString(), {
				signal
			});
			if (!response.ok) {
				if (response.status === 404) {
					loadError = 'No ways found at this location';
				} else {
					loadError = 'Failed to find nearby way';
				}
				return null;
			}

			const data = await response.json();
			if (
				!data ||
				typeof data !== 'object' ||
				typeof (data as { id?: unknown }).id !== 'number' ||
				typeof (data as { tags?: unknown }).tags !== 'object' ||
				(data as { tags?: unknown }).tags === null ||
				Array.isArray((data as { tags?: unknown }).tags)
			) {
				loadError = 'Unexpected response from server';
				return null;
			}

			loadError = null;
			return {
				id: (data as { id: number }).id,
				tags: (data as { tags: Record<string, string> }).tags
			};
		} catch (err) {
			if (err instanceof DOMException && err.name === 'AbortError') {
				return null;
			}

			const message =
				err instanceof Error && err.message ? err.message : 'Error fetching nearest way';
			loadError = message;
			console.error('Error fetching nearest way:', err);
			return null;
		} finally {
			isLoadingWay = false;
		}
	}

	// Map interactions:
	async function onMapClick(latlng: [number, number]): Promise<void> {
		const [lat, lng] = latlng;
		if (!mapInstance) return;

		// Disable normal map click behavior while Add Review is active
		if (wayState.isAddReviewActive) {
			return;
		}

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
			if (nearestWayAbortController) {
				nearestWayAbortController.abort();
			}

			const controller = new AbortController();
			nearestWayAbortController = controller;
			const wayData = await fetchNearestWay(lat, lng, controller.signal);
			if (controller !== nearestWayAbortController) {
				return;
			}
			if (wayData) {
				wayState.selectedWay = {
					id: wayData.id,
					tags: wayData.tags
				};
				await loadSelectedWayHighlight(wayData.id);
				updateUrlWithWay(wayData.id);
			}
		}
	}

	async function handleFindRouteClick() {
		if (!mapInstance) return;
		try {
			await findRoutes({ mapInstance, k: 3 });
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
			mapInstance.removeSelectedWayLayer();
		}
	}

	$effect(() => {
		// Explicitly track mapInstance so effect re-runs when it's initialized
		const map = mapInstance;
		const adjacentWays = wayState.adjacentWays;
		const additionalSelectedWayIds = wayState.additionalSelectedWayIds;
		const selectedWay = wayState.selectedWay;
		const clickHandler = wayState.onAdjacentWayClick;

		if (!map) return;

		const selectedWayId = selectedWay?.id;

		// Filter out ways that have been selected (convert IDs to numbers for comparison)
		// Also exclude the original selected way (shown in blue)
		const unselectedAdjacentWays = adjacentWays.filter((way) => {
			const wayId = way.properties?.id;
			if (wayId === undefined) return true;
			const numWayId = Number(wayId);
			// Exclude if it's in additional selected or is the original selected way
			return !additionalSelectedWayIds.includes(numWayId) && numWayId !== selectedWayId;
		});

		if (unselectedAdjacentWays.length > 0) {
			const featureCollection: GeoJSON.FeatureCollection = {
				type: 'FeatureCollection',
				features: unselectedAdjacentWays
			};
			map.loadAdjacentWaysLayer(featureCollection, clickHandler ?? undefined);
		} else {
			map.removeAdjacentWaysLayer();
		}

		// Show selected adjacent ways in green (but NOT the original selected way)
		if (additionalSelectedWayIds.length > 0) {
			const selectedAdjacentWays = adjacentWays.filter((way) => {
				const wayId = way.properties?.id;
				if (wayId === undefined) return false;
				const numWayId = Number(wayId);
				// Include only if in additional selected AND not the original selected way
				return additionalSelectedWayIds.includes(numWayId) && numWayId !== selectedWayId;
			});
			if (selectedAdjacentWays.length > 0) {
				const selectedFeatureCollection: GeoJSON.FeatureCollection = {
					type: 'FeatureCollection',
					features: selectedAdjacentWays
				};
				map.loadAdditionalSelectedWaysLayer(selectedFeatureCollection);
			} else {
				map.removeAdditionalSelectedWaysLayer();
			}
		} else {
			map.removeAdditionalSelectedWaysLayer();
		}
	});

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

		{#if isLoadingWay}
			<div
				role="status"
				class="absolute top-2 right-2 z-20 bg-blue-600 text-white px-3 py-2 rounded shadow"
			>
				Finding nearest way...
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
