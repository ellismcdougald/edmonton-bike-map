<script lang="ts">
	let { closePopup } = $props();

	let rating: number | null = null;
	let reviewText: string | null = null;
</script>

<div class="fixed inset-0 bg-black/25 flex items-center justify-center z-50">
	<form class="bg-white rounded p-6 shadow max-w-lg w-full max-h-[80vh] overflow-auto">
		<label class="block mb-2 font-semibold" for="rating">Rating</label>
		<input
			id="rating"
			type="number"
			min="1"
			max="10"
			placeholder="5"
			class="w-full border border-gray-300 rounded px-3 py-2 mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
			required
		/>

		<label class="block mb-2 font-semibold" for="reviewText">Review</label>
		<textarea
			id="reviewText"
			rows="4"
			placeholder="Write your review here..."
			class="w-full border border-gray-300 rounded px-3 py-2 mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
			required
		></textarea>

		<div class="flex justify-end gap-2">
			<button type="button" class="px-4 py-2 rounded border" onclick={closePopup}>Cancel</button>
			<button type="submit" class="px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700"
				>Submit</button
			>
		</div>
	</form>
</div>
