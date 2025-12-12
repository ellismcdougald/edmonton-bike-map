<script lang="ts">
	import { createReview } from '$lib/api/client/reviews';
	import { getAdjacentWays } from '$lib/api/client/ways';
	import { wayState } from '$lib/state.svelte';
	import type { WayFeature, WayFeatureGeoJSON } from '$lib/types';
	import { onDestroy } from 'svelte';

	let { wayId, onSubmitted }: { wayId: number; onSubmitted?: () => void } = $props();

	let rating: number | null = $state(null);
	let comment: string | null = $state(null);
	let errorMsg: string = $state('');
	let isSubmitting: boolean = $state(false);
	let adjacentWays: WayFeatureGeoJSON[] = $state([]);

	let selectedWay: WayFeature | null = $derived(wayState.selectedWay);

	$effect(() => {
		// Re-load adjacent ways whenever wayId changes
		loadAdjacentWays();
	});

	onDestroy(() => {
		wayState.adjacentWays = [];
	});

	async function loadAdjacentWays() {
		try {
			const featureCollection = await getAdjacentWays(wayId);
			adjacentWays = featureCollection.features;
			wayState.adjacentWays = adjacentWays;
		} catch (err) {
			console.error('Failed to load adjacent ways:', err);
			adjacentWays = [];
			wayState.adjacentWays = [];
		}
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

		isSubmitting = true;
		errorMsg = '';
		try {
			await createReview(wayId, ratingNum, comment);
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
		<div id="includedWays" class="border border-gray-300 rounded px-3 py-2">
			{#if selectedWay}
				<ul class="flex flex-wrap gap-2">
					<li class="px-2 py-1 rounded-full bg-gray-200 text-sm">
						{selectedWay.tags?.name ? selectedWay.tags.name : `Way #${selectedWay.id}`}
					</li>
				</ul>
			{:else}
				<p class="text-sm text-gray-600">No way selected yet.</p>
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
