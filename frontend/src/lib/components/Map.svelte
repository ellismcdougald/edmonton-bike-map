<script lang="ts">
  import {onMount, onDestroy} from 'svelte';
  import type { LeafletMap } from '$lib/Map';

  let mapInstance: LeafletMap | null

  onMount(async () => {
    const {LeafletMap } = await import('$lib/Map');

    mapInstance = new LeafletMap();
    mapInstance.addTileLayer(
      "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
      19,
      10,
      "© OpenStreetMap contributors"
    );

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