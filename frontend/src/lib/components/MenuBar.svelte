<script lang="ts">
	import { getUsernameFromToken } from "$lib/utils/auth";
	import { onMount } from "svelte";

  let username: string | null = null;

  onMount(() => {
    username = getUsernameFromToken();
  })
</script>

<div class="flex items-center justify-between px-4 py-2 w-full">
  <div class="w-1/3"></div>
  <div class="w-1/3">
    <h1 class="text-3xl font-bold text-center">Edmonton Bike Map</h1>
  </div>
  <div class="w-1/3 text-right">
    {#if username}
      <h2>{username}</h2>
    {:else}
      <a href="/login" class="text-blue-500 underline hover:text-blue-700">
        Log In
      </a>
    {/if}
  </div>
</div>