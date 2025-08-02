<!--
  MapComponent.svelte

  State:
  - mapInstance: LeafletMap | null — reference to the map instance
  - selectStartActive: boolean — whether "select start" mode is active
  - selectEndActive: boolean — whether "select end" mode is active

  Behaviour:
  - Initializes LeafletMap on mount, loads tile and info layers
  - Listens for map clicks to set start/end markers depending on mode
  - Toggles select start/end modes with buttons, ensuring only one active at a time
  - Fetches route from backend API with start/end coordinates on "Find Route"
  - Cleans up map instance on destroy
-->

<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type { LeafletMap } from '$lib/Map';
	import type { WayFeature } from '$lib/types';
	import { wayState } from '$lib/state.svelte';

	let mapInstance: LeafletMap | null;

	const allWaysEndpoint = 'http://localhost:8080/api/all-ways';

	let selectStartActive: boolean = false;
	let selectEndActive: boolean = false;

	function onMapClick(latlng: [number, number]): void {
		const [lat, lng] = latlng;

		if (!mapInstance) return;

		if (selectStartActive) {
			mapInstance.removeStartMarker();
			mapInstance.addStartMarker([lat, lng]);
			toggleSelectStartActive();
			mapInstance.showInfoLayer();
		} else if (selectEndActive) {
			mapInstance.removeEndMarker();
			mapInstance.addEndMarker([lat, lng]);
			toggleSelectEndActive();
			mapInstance.showInfoLayer();
		}
	}

	function toggleSelectStartActive() {
		if (selectStartActive) {
			selectStartActive = false;
			mapInstance?.showInfoLayer();
		} else {
			selectStartActive = true;
			selectEndActive = false;
			mapInstance?.hideInfoLayer();
		}
	}

	function toggleSelectEndActive() {
		if (selectEndActive) {
			selectEndActive = false;
			mapInstance?.showInfoLayer();
		} else {
			selectEndActive = true;
			selectStartActive = false;
			mapInstance?.hideInfoLayer();
		}
	}

	function findRoute(): void {
		if (!mapInstance) return;

		const startLatLng: [number, number] | null = mapInstance.getStartLatLng();
		const endLatLng: [number, number] | null = mapInstance.getEndLatLng();

		if (!startLatLng || !endLatLng) {
			alert('Make sure you have selected both a start and an end point!');
			return;
		}

		const params = new URLSearchParams({
			startLatitude: startLatLng[0].toString(),
			startLongitude: startLatLng[1].toString(),
			endLatitude: endLatLng[0].toString(),
			endLongitude: endLatLng[1].toString()
		});

		console.log(params.toString());

		fetch(`http://localhost:8080/api/route?${params.toString()}`)
			.then((res) => {
				if (!res.ok) throw new Error('Failed to get route data');
				return res.json();
			})
			.then((geojson) => {
				if (!mapInstance) return;

				mapInstance.removeRouteLayer();
				mapInstance.loadRouteLayer(geojson);
			})
			.catch((err) => {
				alert('Error fetching or displaying route: ' + err.message);
				console.error(err);
			});
	}

	onMount(async () => {
		const { LeafletMap } = await import('$lib/Map');

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
		} catch {
			console.error('Error loading info layer.');
		}

		let selectStartButton = document.getElementById('selectStartButton');
		selectStartButton?.addEventListener('click', () => {
			toggleSelectStartActive();
			selectStartButton.classList.toggle('active', selectStartActive);
		});

		let selectEndButton = document.getElementById('selectEndButton');
		selectEndButton?.addEventListener('click', () => {
			toggleSelectEndActive();
			selectEndButton.classList.toggle('active', selectEndActive);
		});

		let findRouteButton = document.getElementById('findRouteButton');
		findRouteButton?.addEventListener('click', findRoute);
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
			>Select Start Location</button
		>
		<button
			type="button"
			id="selectEndButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow"
			>Select End Location</button
		>
		<button
			type="button"
			id="findRouteButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow"
			>Find Route</button
		>
		<button
			type="button"
			id="resetButton"
			class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow"
			>Reset</button
		>
	</div>
</div>
