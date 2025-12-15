<!--
  Review.svelte

  Purpose:
  Displays a single user review with metadata (username, date, rating, and comment).

  Props:
  - review (Review): Review object containing username, createdAt, rating, and comment

  State:
  - none

  Behavior:
  - Formats createdAt to show only YYYY-MM-DD
  - Displays rating as "x / 10"

  Notes:
  - Purely presentational; no external dependencies besides $lib/types
-->

<script lang="ts">
	import type { Review as ReviewObj } from '$lib/types';
	import { deleteReview } from '$lib/api/client/reviews';

	let {
		review,
		currentUser,
		wayId,
		onDeleted
	}: {
		review: ReviewObj;
		currentUser?: string | null;
		wayId?: number;
		onDeleted?: () => void;
	} = $props();

	let deleting: boolean = $state(false);

	async function handleDelete() {
		if (!wayId) return;
		deleting = true;
		try {
			await deleteReview(wayId);
			if (onDeleted) onDeleted();
		} catch (err) {
			console.error('Failed to delete review', err);
		} finally {
			deleting = false;
		}
	}
</script>

<div class="mb-3 rounded-md border border-gray-200 bg-white p-3 shadow-sm">
	<div class="flex items-start gap-2 text-sm text-gray-600">
		<div class="flex flex-col gap-1">
			<div class="flex flex-wrap items-center gap-2">
				<span class="font-semibold text-gray-800">{review.username}</span>
				<span>•</span>
				<span>{review.createdAt.substring(0, 10)}</span>
			</div>
			<div class="text-gray-800">{`Rating: ${review.rating} / 10`}</div>
		</div>
		{#if currentUser && review.username === currentUser}
			<button
				type="button"
				class="ml-auto inline-flex items-center gap-1 rounded-full px-2 py-1 text-xs font-medium text-red-600 ring-1 ring-red-200 hover:bg-red-50 disabled:opacity-60"
				disabled={deleting}
				onclick={handleDelete}
			>
				<span aria-hidden="true">🗑</span>
				<span>{deleting ? 'Removing' : 'Delete'}</span>
			</button>
		{/if}
	</div>
	<p class="mt-2 text-sm text-gray-900">{review.comment}</p>
</div>

<style></style>
