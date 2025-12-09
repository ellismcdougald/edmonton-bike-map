<!--
  Sidebar.svelte
  Displays metadata and reviews for the currently selected way in a toggleable sidebar.
-->

<script lang="ts">
	import ReviewContainer from './ReviewContainer.svelte';
	import { wayState } from '$lib/state.svelte';
	import type { WayFeature } from '$lib/types';
	import { determineRouteType, determineBicycleRoute } from '$lib/utils/route';
	import type { Review as ReviewObj } from '$lib/types';
	import { computeAverageRating } from '$lib/utils/review';
	import { capitalizeFirstLetter } from '$lib/utils/helpers';
	import { fetchReviews } from '$lib/api/reviews';

	let way: WayFeature | null = $derived(wayState.selectedWay);
	let reviews: ReviewObj[] = $state([]);

	async function loadReviews(wayId: number) {
		try {
			reviews = await fetchReviews(wayId);
		} catch (error) {
			console.error('Error loading reviews:', error);
		}
	}

	let sidebarLoaded: boolean = $state(false);
	let isVisible: boolean = $state(true);
	$effect(() => {
		if (way) {
			isVisible = true;
			loadReviews(way.id);
			sidebarLoaded = true;
		}
	});

	function toggleSidebar(): void {
		isVisible = !isVisible;
	}
</script>

{#if sidebarLoaded}
	<div id="hide-sidebar-button-container" class="flex ml-auto">
		<button
			id="hide-sidebar-button"
			class="absolute bottom-2 right-2 bg-gray-300 hover:bg-gray-400 text-gray-800 px-3 py-1 rounded z-20"
			onclick={toggleSidebar}>{isVisible ? 'Hide' : 'Show'}</button
		>
	</div>
{/if}

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
				Surface: <span class="font-normal"
					>{way.tags.surface ? capitalizeFirstLetter(way.tags.surface) : 'Unknown'}</span
				>
			</h2>
			<h2 class="font-semibold">
				Average Rating: <span class="font-normal"
					>{reviews.length > 0 ? `${computeAverageRating(reviews)} / 10` : 'TBD'}</span
				>
			</h2>
		</section>

		<ReviewContainer wayId={way.id} {reviews} onReviewAdded={() => loadReviews(way.id)} />
	</div>
{/if}

<style></style>
