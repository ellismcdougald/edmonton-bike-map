<script lang="ts">
	import '../../app.css';
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
