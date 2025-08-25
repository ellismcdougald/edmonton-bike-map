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
	import { getUsernameFromToken, removeToken } from '$lib/utils/auth';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let username: string | null = null;

	onMount(() => {
		username = getUsernameFromToken();
	});

	function logout() {
		removeToken();
		goto('/login');
	}
</script>

<div class="flex items-center justify-between px-4 py-2 w-full">
	<div class="w-1/3"></div>
	<div class="w-1/3">
		<h1 class="text-3xl font-bold text-center">Edmonton Bike Map</h1>
	</div>
	<div class="w-1/3 flex justify-end">
		<div class="relative inline-block group">
			<button
				id="usernameButton"
				class="text-black text-center px-4 py-1 hover:underline focus:outline-none"
			>
				{username}
			</button>
			<div
				class="absolute left-0 mt-2 w-full bg-white text-gray-800 rounded shadow-lg z-50
             opacity-0 invisible translate-y-1 group-hover:visible group-hover:opacity-100 group-hover:translate-y-0
             transition-all duration-200 ease-out"
			>
				<button
					id="logoutButton"
					class="block w-full px-4 py-2 text-center hover:bg-gray-100"
					on:click={logout}
				>
					Log out
				</button>
			</div>
		</div>
	</div>
</div>
