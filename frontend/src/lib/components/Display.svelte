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