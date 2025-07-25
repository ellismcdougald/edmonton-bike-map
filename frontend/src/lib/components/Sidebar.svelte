<script lang="ts">
  import { selectedWay } from "$lib/stores/selectedWay";

  $: way = $selectedWay;

	let isVisible: boolean = true;

	function toggleSidebar() {
		isVisible = !isVisible;
	}
</script>

<div id="hide-sidebar-button-container" class="flex ml-auto">
	<button
		id="hide-sidebar-button"
    class="absolute bottom-2 right-2 bg-gray-300 hover:bg-gray-400 text-gray-800 px-3 py-1 rounded z-20"
		on:click={toggleSidebar}>{isVisible ? 'Hide' : 'Show'}</button
	>
</div>

{#if isVisible}
	<div id="sidebar-content" class="w-full p-2 h-full bg-white shadow-2xl pb-12">
		<h1 class="text-xl font-bold">{way?.tags?.name ? way?.tags?.name : "Unnamed Route"}</h1>

		<section class="mb-2">
			<h2 class="font-semibold">Type: <span class="font-normal">Bike Lane</span></h2>
			<h2 class="font-semibold">Surface: <span class="font-normal">Paved</span></h2>
			<h2 class="font-semibold">Rating: <span class="font-normal">8 / 10</span></h2>
		</section>

		<section id="sidebar-review-container" class="border-t border-dotted border-gray-600 pt-1">
			<h1 class="text-lg font-bold mb-2">Reviews:</h1>

			<div class="mb-2 border border-gray-600 p-1 rounded-xs">
				<div class="flex gap-2 text-sm text-gray-600 mb-1">
					<span>User: emcd84</span>
					<span>—</span>
					<span>Date: July 20, 2025</span>
					<span>—</span>
					<span>Rating: 10 / 10</span>
				</div>
				<p>Well-paved bike lane, separated from road by barrier. On a quiet road.</p>
			</div>

			<div class="mb-2 border border-gray-600 p-1 rounded-xs">
				<div class="flex gap-2 text-sm text-gray-600 mb-1">
					<span>User: js1234</span>
					<span>—</span>
					<span>Date: July 10, 2025</span>
					<span>—</span>
					<span>Rating: 6 / 10</span>
				</div>
				<p>Good bike lane but often blocked by construction.</p>
			</div>
		</section>
	</div>
{/if}

<style></style>
