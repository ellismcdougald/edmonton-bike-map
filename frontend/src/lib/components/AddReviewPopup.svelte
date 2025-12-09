<script lang="ts">
	import { createReview } from '$lib/api/reviews';

	let { closePopup, wayId, onReviewAdded } = $props();

	let rating: number | null = $state(null);
	let comment: string | null = $state(null);

	let errorMsg: string = $state('');
	let isSubmitting: boolean = $state(false);

	async function handleSubmit(event: Event) {
		event.preventDefault();

		if (rating == null || rating == undefined) {
			errorMsg = 'Please provide a rating!';
			return;
		}

		isSubmitting = true;
		errorMsg = '';

		try {
			await createReview(wayId, rating, comment);

			rating = null;
			comment = null;

			if (onReviewAdded && typeof onReviewAdded === 'function') onReviewAdded?.();
			closePopup();
		} catch (err: unknown) {
			errorMsg = err instanceof Error ? err.message : String(err);
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div class="fixed inset-0 bg-black/25 flex items-center justify-center z-50" id="addReviewPopup">
	<form
		onsubmit={handleSubmit}
		class="bg-white rounded p-6 shadow max-w-lg w-full max-h-[80vh] overflow-auto"
	>
		<label class="block mb-2 font-semibold" for="rating">Rating</label>
		<input
			id="rating"
			type="number"
			min="1"
			max="10"
			placeholder="5"
			class="w-full border border-gray-300 rounded px-3 py-2 mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
			required
			bind:value={rating}
		/>

		<label class="block mb-2 font-semibold" for="comment">Review</label>
		<textarea
			id="comment"
			rows="4"
			placeholder="Write your review here..."
			class="w-full border border-gray-300 rounded px-3 py-2 mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
			bind:value={comment}
		></textarea>

		{#if errorMsg}
			<p class="text-red-600 mb-4">{errorMsg}</p>
		{/if}

		<div class="flex justify-end gap-2">
			<button type="button" id="cancelButton" class="px-4 py-2 rounded border" onclick={closePopup}
				>Cancel</button
			>
			<button
				type="submit"
				id="submitButton"
				class="px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700"
				disabled={isSubmitting}
			>
				Submit
			</button>
		</div>
	</form>
</div>
