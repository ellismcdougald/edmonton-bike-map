<!--
  ReviewContainer.svelte
  ----------------------
  Displays a list of reviews for the given wayId.
  Fetches reviews from the backend whenever wayId changes.
  Allows users to open a popup to add a new review.

  Props:
    - wayId (number): The ID of the selected way to fetch reviews for.

  State:
    - addReviewActive (boolean): Controls visibility of the AddReviewPopup.
    - reviews (ReviewObj[]): List of reviews fetched from the backend.

  Behavior:
    - Automatically fetches reviews when wayId changes.
    - Opens AddReviewPopup on button click.
    - Passes wayId and closePopup callback to AddReviewPopup.
    - Renders each review using the Review component.

  Improvements:
    - Could implement a loading screen while reviews are being fetched.$$render
    - Could add a message when no reviews are present
-->

<script lang="ts">
	import Review from './Review.svelte';
	import AddReviewPopup from './AddReviewPopup.svelte';
  import type { Review as ReviewObj } from '$lib/types';
  
  let { wayId }: { wayId: number } = $props();

	let addReviewActive: boolean = $state(false);
  let reviews: ReviewObj[] = $state([])

  async function loadReviews(wayId: number) {
    try {
      const res = await fetch(`http://localhost:8080/api/reviews?wayID=${wayId}`)
      if (!res.ok) throw new Error(`HTTP error! Status: ${res.status}`)
      reviews = await res.json();
    } catch (e) {
      console.error("Failed to fetch reviews: ", e)
      return []
    }
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
			}}
      wayId={wayId}
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

<style></style>
