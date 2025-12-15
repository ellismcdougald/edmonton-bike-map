<!--
	Reviews.svelte

	Purpose:
	Container for displaying and managing reviews for a specific wayId.

	Props:
	- wayId (number): The ID of the way/route to fetch and display reviews for

	State:
	- reviews (Review[]): Array of reviews for the given wayId
	- Inline add-review form state (rating, comment, error, isSubmitting)

	Behavior:
	- Displays reviews list
	- Shows inline form (instead of popup) for adding a review when allowed
	- Calls onReviewAdded after successful submit so parent can refresh reviews

	Notes:
	- Relies on createReview from $lib/api/client/reviews
	- Uses review.createdAt + review.rating as key in {#each} loop
-->

<script lang="ts">
	import type { Review as ReviewObj } from '$lib/types';
	import AddReview from './AddReview.svelte';
	import ViewReviews from './ViewReviews.svelte';

	let {
		wayId,
		reviews,
		onReviewAdded,
		canReview = true,
		username
	}: {
		wayId: number;
		reviews: ReviewObj[];
		onReviewAdded?: () => void;
		canReview?: boolean;
		username?: string | null;
	} = $props();

	let showForm: boolean = $state(false);

	// selectedWay is consumed within AddReview; no local binding needed here

	function resetForm() {
		showForm = false;
	}

	function handleSubmitted() {
		if (onReviewAdded && typeof onReviewAdded === 'function') onReviewAdded();
		showForm = false;
	}

	// Determine if the current user has already reviewed this way
	const hasReviewed: boolean = $derived(
		Boolean(username) && Array.isArray(reviews) && reviews.some((r) => r.username === username)
	);

	const canAddReview: boolean = $derived(Boolean(canReview) && !hasReviewed);
</script>

<div id="review-container" class="border-t">
	<!-- Header at the top: title + action button -->
	<div class="flex justify-between items-center mt-1 mb-2">
		<h1 class="text-2xl font-bold mb-2 mt-1">{showForm ? 'Add Review' : 'Reviews'}</h1>
		{#if canAddReview}
			{#if showForm}
				<button
					type="button"
					id="cancelAddReviewButton"
					class="px-3 py-1 rounded border"
					onclick={resetForm}
				>
					Back to reviews
				</button>
			{:else}
				<button
					type="button"
					id="addReviewButton"
					class="bg-blue-600 text-white px-2 py-1 rounded hover:bg-blue-700"
					onclick={() => {
						showForm = true;
					}}
				>
					Add Review
				</button>
			{/if}
		{:else}
			<span class="text-sm text-gray-600">
				{#if !username}
					Log in to add a review.
				{:else}
					You’ve already reviewed this route.
				{/if}
			</span>
		{/if}
	</div>

	{#if showForm}
		<AddReview {wayId} onSubmitted={handleSubmitted} />
	{:else}
		<ViewReviews {reviews} currentUser={username} {wayId} onDeleted={onReviewAdded} />
	{/if}
</div>
