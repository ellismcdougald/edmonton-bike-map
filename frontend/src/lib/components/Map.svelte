<script lang="ts">
  import {onMount, onDestroy} from 'svelte';
  import type { LeafletMap } from '$lib/Map';
  import type { WayFeature } from '$lib/types';
  import { selectedWay } from '$lib/stores/selectedWay';

  let mapInstance: LeafletMap | null

  const allWaysEndpoint = "http://localhost:8080/api/all-ways"

  onMount(async () => {
    const {LeafletMap } = await import('$lib/Map');

    mapInstance = new LeafletMap();
    mapInstance.addTileLayer(
      "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
      19,
      10,
      "© OpenStreetMap contributors"
    );

    try {
      const res = await fetch(allWaysEndpoint);
      if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
      const geojson = await res.json()
      mapInstance.loadInfoLayer(geojson, 
        { color: 'red', weight: 1, opacity: 0 }, 
        (way: WayFeature) => {
          selectedWay.set(way);
      });
    } catch (err) {
      console.error("Error loading info layer.");
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
    <button type="button" id="selectStartButton" class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow">Select Start Location</button>
    <button type="button" id="selectEndButton" class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow">Select End Location</button>
    <button type="button" id="findRouteButton" class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow">Find Route</button>
    <button type="button" id="resetButton" class="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition flex-grow">Reset</button>
  </div>
</div>