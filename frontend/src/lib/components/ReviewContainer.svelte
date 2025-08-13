<script lang="ts">
	import Review from './Review.svelte';
	import AddReviewPopup from './AddReviewPopup.svelte';
	import type { Review as ReviewObj } from '$lib/types';
	import { fetchReviews } from '$lib/review/reviewActions';

	let { wayId }: { wayId: number } = $props();
	let addReviewActive: boolean = $state(false);
	let reviews: ReviewObj[] = $state([]);

	async function loadReviews(wayId: number) {
		reviews = await fetchReviews(wayId);
	}

	$effect(() => {
		loadReviews(wayId);
	});
</script>

<div id="review-container" class="border-t">
	<div class="flex justify-between items-center mt-1 mb-2">
		<h1 class="text-2xl font-bold mb-2 mt-1">Reviews:</h1>
		<button
			class="bg-blue-600 text-white px-2 py-1 rounded hover:bg-blue-700"
			onclick={() => {
				addReviewActive = true;
			}}>Add Review</button
		>
	</div>

	{#if addReviewActive}
		<AddReviewPopup
			closePopup={() => {
				addReviewActive = false;
				loadReviews(wayId);
			}}
			{wayId}
		/>
	{/if}

	{#if reviews}
		<div>
			{#each reviews as review (review.createdAt + review.rating)}
				<Review {review} />
			{/each}
		</div>
	{/if}
</div>
