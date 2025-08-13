<!--
  Sidebar.svelte

  State:
  - way: WayFeature | null (derived from wayState.selectedWay)
  - isVisible: boolean (sidebar visibility)

  Behaviour:
  - Toggles sidebar visibility on button click
  - Displays metadata about the selected way
  - Passes way.id to ReviewContainer
  - Determines route type and bicycle route status using helper functions

  Issues:
  - Clean up and improve the logic for determineRouteType and determineBicycleRoute
-->

<script lang="ts">
	import ReviewContainer from './ReviewContainer.svelte';
	import { wayState } from '$lib/state.svelte';
	import type { WayFeature } from '$lib/types';
	import { determineRouteType, determineBicycleRoute } from '$lib/utils/route';

	let way: WayFeature | null = $derived(wayState.selectedWay);

	let isVisible: boolean = $state(true);

	function toggleSidebar(): void {
		isVisible = !isVisible;
	}
</script>

<div id="hide-sidebar-button-container" class="flex ml-auto">
	<button
		id="hide-sidebar-button"
		class="absolute bottom-2 right-2 bg-gray-300 hover:bg-gray-400 text-gray-800 px-3 py-1 rounded z-20"
		onclick={toggleSidebar}>{isVisible ? 'Hide' : 'Show'}</button
	>
</div>

{#if isVisible && way}
	<div id="sidebar-content" class="w-full p-2 h-full bg-white shadow-2xl pb-12">
		<h1 class="text-xl font-bold">{way.tags?.name ? way?.tags?.name : 'Unnamed Route'}</h1>

		<section class="mb-2">
			<h2 class="font-semibold">
				Type: <span class="font-normal">{determineRouteType(way.tags)}</span>
			</h2>
			<h2 class="font-semibold">
				Bicycle Route: <span class="font-normal">{determineBicycleRoute(way.tags)}</span>
			</h2>
			<h2 class="font-semibold">
				Surface: <span class="font-normal">{way.tags.surface ? way.tags.surface : 'Unknown'}</span>
			</h2>
			<h2 class="font-semibold">Rating: <span class="font-normal">8 / 10</span></h2>
		</section>

		<ReviewContainer wayId={way.id} />
	</div>
{/if}

<style></style>
