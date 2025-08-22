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
	let { review }: { review: ReviewObj } = $props();
</script>

<div class="mb-2 border border-gray-600 p-1 rounded-xs">
	<div class="flex gap-2 text-sm text-gray-600 mb-1">
		<span>User: {review.username}</span>
		<span>—</span>
		<span>{review.createdAt.substring(0, 10)}</span>
		<span>—</span>
		<span>{`Rating: ${review.rating} / 10`}</span>
	</div>
	<p>{review.comment}</p>
</div>

<style></style>
