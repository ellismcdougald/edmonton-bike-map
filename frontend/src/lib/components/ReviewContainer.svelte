<script lang="ts">
	import Review from './Review.svelte';
	import AddReviewPopup from './AddReviewPopup.svelte';

	let addReviewActive: boolean = false;
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
		/>
	{/if}

	<Review />
</div>

<style></style>
