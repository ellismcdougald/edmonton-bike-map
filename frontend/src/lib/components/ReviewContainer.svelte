<!--
  Reviews.svelte

  Purpose:
  Container for displaying and managing reviews for a specific wayId.

  Props:
  - wayId (number): The ID of the way/route to fetch and display reviews for

  State:
  - addReviewActive (boolean): Tracks whether the AddReviewPopup is open
  - reviews (Review[]): Array of reviews for the given wayId

  Behavior:
  - On mount and whenever wayId changes, fetches reviews from backend
  - Displays "Add Review" button that toggles AddReviewPopup
  - After a new review is added, reloads reviews list
  - Renders each review using Review.svelte

  Notes:
  - Relies on fetchReviews from $lib/utils/review
  - Uses review.createdAt + review.rating as key in {#each} loop
-->

<script lang="ts">
	import Review from './Review.svelte';
	import AddReviewPopup from './AddReviewPopup.svelte';
	import type { Review as ReviewObj } from '$lib/types';
	import { fetchReviews } from '$lib/utils/review';

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
			id="addReviewButton"
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
