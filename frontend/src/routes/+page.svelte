<!--
  src/routes/+page.svelte

  Purpose:
  Main app route for authenticated users. Displays the app layout with MenuBar and Display components.

  State:
  - none (stateless route; purely layout and authentication check)

  Behavior:
  - On mount, checks for JWT token in localStorage
      - If token is missing, redirects user to '/login'
  - Renders MenuBar and Display components

  Notes:
  - Depends on $lib/components/MenuBar.svelte and $lib/components/Display.svelte
  - Uses SvelteKit's `goto` for client-side navigation
  - Enforces simple authentication guard via localStorage token
-->

<script lang="ts">
	import '../app.css';
	import Display from '$lib/components/Display.svelte';
	import MenuBar from '$lib/components/MenuBar.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(() => {
		const token = localStorage.getItem('token');
		if (!token) {
			goto('/login');
		}
	});
</script>

<MenuBar />
<Display />
