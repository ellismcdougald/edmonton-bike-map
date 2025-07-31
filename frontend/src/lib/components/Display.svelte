<!--
  Layout component that displays the map and sidebar.

  Structure:
  - The map fills the main container area.
  - The sidebar is positioned absolutely on the top right, overlaying the map.

  Behavior:
  - Sidebar has fixed width (25%) and full height.
  - Map takes full height and width of remaining area behind sidebar.

  Usage:
  - Import this component to render the main map view with sidebar controls.
-->

<script lang="ts">
  import Map from './Map.svelte';
  import Sidebar from './Sidebar.svelte';
</script>

<style></style>

<div id="display-container" class="relative h-93/100">
  <div id="map-container" class="h-full w-full relative z-0">
    <Map />
  </div>
  <div id="sidebar-container" class="absolute top-0 right-0 w-1/4 h-full z-10">
    <Sidebar />
  </div>
</div>