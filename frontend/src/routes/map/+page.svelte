<script lang="ts">
	import '../../app.css';
	import Display from '$lib/components/Display.svelte';
	import MenuBar from '$lib/components/MenuBar.svelte';

	export let data;
</script>

<MenuBar />
<Display ways={data.ways} />
