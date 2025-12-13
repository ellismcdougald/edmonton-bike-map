<script lang="ts">
	import { createReview } from '$lib/api/client/reviews';
	import { getAdjacentWays } from '$lib/api/client/ways';
	import { wayState } from '$lib/state.svelte';
	import type { WayFeature, WayFeatureGeoJSON } from '$lib/types';
	import { onDestroy } from 'svelte';
	import { SvelteSet } from 'svelte/reactivity';

	let { wayId, onSubmitted }: { wayId: number; onSubmitted?: () => void } = $props();

	let rating: number | null = $state(null);
	let comment: string | null = $state(null);
	let errorMsg: string = $state('');
	let isSubmitting: boolean = $state(false);
	let adjacentWays: WayFeatureGeoJSON[] = $state([]);
	let selectedWayIds: number[] = $state([wayId]);
	let adjacentLoadSeq = 0;

	let selectedWay: WayFeature | null = $derived(wayState.selectedWay);

	$effect(() => {
		// Reset selection to just the initial wayId when it changes
		selectedWayIds = [wayId];
		wayState.additionalSelectedWayIds = [];
		wayState.isAddReviewActive = true;

		// Set up click handler for adjacent ways
		wayState.onAdjacentWayClick = toggleWaySelection;
	});

	$effect(() => {
		// Sync additional selected way IDs to state (excluding the original wayId)
		wayState.additionalSelectedWayIds = selectedWayIds.filter((id) => id !== wayId);

		// Load adjacent ways for all selected ways
		loadAdjacentWays();
	});

	onDestroy(() => {
		wayState.adjacentWays = [];
		wayState.additionalSelectedWayIds = [];
		wayState.onAdjacentWayClick = null;
		wayState.isAddReviewActive = false;
	});

	async function loadAdjacentWays() {
		const seq = ++adjacentLoadSeq;
		try {
			// Fetch adjacent ways for all selected ways
			const allAdjacentWaysPromises = selectedWayIds.map((id) => getAdjacentWays(id));
			const allFeatureCollections = await Promise.all(allAdjacentWaysPromises);
			if (seq !== adjacentLoadSeq) return;

			// Combine all features and deduplicate by way ID
			const seenIds = new SvelteSet<number>();
			const uniqueFeatures: WayFeatureGeoJSON[] = [];

			for (const featureCollection of allFeatureCollections) {
				for (const feature of featureCollection.features) {
					const featureId = feature.properties?.id;
					if (featureId !== undefined && !seenIds.has(Number(featureId))) {
						seenIds.add(Number(featureId));
						uniqueFeatures.push(feature);
					}
				}
			}

			adjacentWays = uniqueFeatures;
			wayState.adjacentWays = adjacentWays;
		} catch (err) {
			if (seq !== adjacentLoadSeq) return;
			console.error('Failed to load adjacent ways:', err);
			adjacentWays = [];
			wayState.adjacentWays = [];
		}
	}

	function toggleWaySelection(id: number) {
		if (selectedWayIds.includes(id)) {
			selectedWayIds = selectedWayIds.filter((wid) => wid !== id);
		} else {
			selectedWayIds = [...selectedWayIds, id];
		}
	}

	function getWayName(id: number): string {
		if (id === wayId && selectedWay?.tags?.name) {
			return String(selectedWay.tags.name);
		}
		const adjacent = adjacentWays.find((w) => w.properties?.id === id);
		if (adjacent?.properties?.tags?.name) {
			return String(adjacent.properties.tags.name);
		}
		return `Way #${id}`;
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();

		if (rating == null) {
			errorMsg = 'Please provide a rating!';
			return;
		}
		const ratingNum = Number(rating);
		if (!Number.isFinite(ratingNum)) {
			errorMsg = 'Please provide a valid numeric rating!';
			return;
		}

		if (selectedWayIds.length === 0) {
			errorMsg = 'Please select at least one way to review!';
			return;
		}

		isSubmitting = true;
		errorMsg = '';
		try {
			await createReview(selectedWayIds, ratingNum, comment);
			rating = null;
			comment = null;
			if (onSubmitted) onSubmitted();
		} catch (err: unknown) {
			errorMsg = err instanceof Error ? err.message : String(err);
		} finally {
			isSubmitting = false;
		}
	}
</script>

<form class="space-y-3 mb-4" onsubmit={handleSubmit}>
	<div>
		<label class="block mb-1 font-semibold" for="includedWays">Included Ways</label>
		<p class="text-xs text-gray-600 mb-2">
			Click orange adjacent ways on the map to include them in your review.
		</p>
		<div id="includedWays" class="border border-gray-300 rounded px-3 py-2">
			{#if selectedWayIds.length > 0}
				<ul class="flex flex-wrap gap-2">
					{#each selectedWayIds as wid (wid)}
						<li class="px-2 py-1 rounded-full bg-blue-200 text-sm flex items-center gap-1">
							<span>{getWayName(wid)}</span>
							{#if wid !== wayId}
								<button
									type="button"
									class="text-red-600 hover:text-red-800 ml-1"
									onclick={() => toggleWaySelection(wid)}
									aria-label="Remove way {wid}"
								>
									×
								</button>
							{/if}
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-sm text-gray-600">No ways selected.</p>
			{/if}
		</div>
	</div>

	<div>
		<label class="block mb-1 font-semibold" for="rating">Rating</label>
		<input
			id="rating"
			type="number"
			min="1"
			max="10"
			placeholder="5"
			class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
			required
			bind:value={rating}
		/>
	</div>

	<div>
		<label class="block mb-1 font-semibold" for="comment">Review</label>
		<textarea
			id="comment"
			rows="3"
			placeholder="Write your review here..."
			class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
			bind:value={comment}
		></textarea>
	</div>

	{#if errorMsg}
		<p class="text-red-600 text-sm">{errorMsg}</p>
	{/if}

	<div class="flex justify-end">
		<button
			type="submit"
			id="submitButton"
			class="px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-60"
			disabled={isSubmitting}
		>
			Submit
		</button>
	</div>
</form>
