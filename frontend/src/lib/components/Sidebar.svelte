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

	let way: WayFeature | null = $derived(wayState.selectedWay);

	let isVisible: boolean = $state(true);

	function toggleSidebar(): void {
		isVisible = !isVisible;
	}

	// determineRouteType uses the tags for a route to determine its road type (arterial, collector, local, etc) and bike infastructure (if possible)
	// TODO: simplify and improve logic, continue to ensure that all common types are classified
	function determineRouteType(tags: Record<string, string>): string {
		if (tags.highway == 'residential') {
			return 'Residential Street';
		} else if (tags.highway == 'tertiary') {
			if (tags.cycleway == 'lane') {
				return 'Collector Road with Bike Lane';
			} else {
				return 'Collector Road';
			}
		} else if (tags.highway == 'path') {
			if (tags.bicycle == 'designated' && tags.foot == 'designated') {
				return 'Shared Use Path';
			} else if (tags.bicycle == 'designated') {
				return 'Bike Path';
			} else if (tags.foot == 'designated') {
				return 'Footpath';
			}
		} else if (tags.highway == 'footway') {
			if (tags.bicycle == 'yes' || tags.bicycle == 'designated') {
				return 'Shared Use Path';
			} else {
				return 'Footpath';
			}
		} else if (tags.highway == 'cycleway') {
			if (tags.foot == 'designated' || tags.foot == 'yes') {
				return 'Shared Use Path';
			} else {
				return 'Cycleway';
			}
		} else if (tags.highway == 'secondary') {
			if (
				tags.cycleway == 'share_busway' ||
				tags['cycleway:left'] == 'share_busway' ||
				tags['cycleway:right'] == 'share_busway'
			) {
				return 'Arterial Road With Bus-Bike Lane';
			}
			return 'Arterial Road';
		} else if (tags.highway == 'primary') {
			return 'Major Arterial Road';
		} else if (tags.highway == 'unclassified') {
			if (
				tags.cycleway == 'separate' ||
				tags['cycleway:left'] == 'separate' ||
				tags['cycleway:right'] == 'separate'
			) {
				return 'Local Road with Separated Bike Lane';
			}
			return 'Local Road';
		}
		return 'Unknown';
	}

	// determineBicycleRoute uses the tags to determine if a route is part of the local bicycle network
	function determineBicycleRoute(tags: Record<string, string>): string {
		if (tags.bicycle == 'designated' || tags.lcn == 'yes') {
			return 'Yes';
		} else {
			return 'No';
		}
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
