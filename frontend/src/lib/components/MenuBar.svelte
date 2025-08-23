<!--
  Header.svelte

  Purpose:
  Displays the app header with title and user login state.

  Props:
  - none

  State:
  - username (string | null): derived from JWT; shows if user is logged in

  Behavior:
  - On mount, extracts username from JWT via getUsernameFromToken()
  - If username exists, displays it on the right side
  - If no username, shows "Log In" link

  Notes:
  - Depends on $lib/utils/auth
-->

<script lang="ts">
	import { getUsernameFromToken } from '$lib/utils/auth';
	import { onMount } from 'svelte';

	let username: string | null = null;

	onMount(() => {
		username = getUsernameFromToken();
	});
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
			<a href="/login" class="text-blue-500 underline hover:text-blue-700"> Log In </a>
		{/if}
	</div>
</div>
