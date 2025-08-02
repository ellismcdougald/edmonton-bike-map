<!--
  AddReviewPopup component

  Props:
  - closePopup: function to close the popup/modal
  - wayId: ID of the currently selected way to associate the review with

  State:
  - rating: number or null, user rating input
  - reviewText: string or null, user review text input
  - errorMsg: string, error message shown on validation or submission failure
  - isSubmitting: boolean, disables form inputs and shows submission state

  Behavior:
  - Validates rating presence before submitting
  - Sends POST request to submit review tied to wayId
  - Resets inputs and closes popup on successful submission
  - Shows error message if submission fails
-->

<script lang="ts">
	let { closePopup, wayId } = $props();

	let rating: number | null = $state(null);
	let reviewText: string | null = $state(null);

	let errorMsg: string = $state('');
	let isSubmitting: boolean = $state(false);

	async function handleSubmit(event: Event) {
		event.preventDefault();
		if (rating == null) {
			errorMsg = 'Please provide a rating!';
      return
		}

		isSubmitting = true;
		try {
			const res = await fetch('http://localhost:8080/api/reviews', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ wayId, rating, reviewText })
			});
			if (!res.ok) {
				throw new Error(`Failed to submit review: ${res.statusText}`);
			}
			rating = null;
			reviewText = null;
			closePopup();
		} catch (err: unknown) {
      if (err instanceof Error) {
        errorMsg = err.message;
      } else {
        errorMsg = String(err);
      }
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div class="fixed inset-0 bg-black/25 flex items-center justify-center z-50">
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

		<label class="block mb-2 font-semibold" for="reviewText">Review</label>
		<textarea
			id="reviewText"
			rows="4"
			placeholder="Write your review here..."
			class="w-full border border-gray-300 rounded px-3 py-2 mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
			required
			bind:value={reviewText}
		></textarea>

		{#if errorMsg}
			<p class="text-red-600 mb-4">{errorMsg}</p>
		{/if}

		<div class="flex justify-end gap-2">
			<button type="button" class="px-4 py-2 rounded border" onclick={closePopup}>Cancel</button>
			<button
				type="submit"
				class="px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700"
				disabled={isSubmitting}
			>
				Submit
			</button>
		</div>
	</form>
</div>
