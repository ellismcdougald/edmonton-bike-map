<script lang="ts">
	import Review from './Review.svelte';
	import type { Review as ReviewObj } from '$lib/types';

	let {
		reviews,
		currentUser,
		wayId,
		onDeleted
	}: {
		reviews: ReviewObj[];
		currentUser?: string | null;
		wayId?: number;
		onDeleted?: () => void;
	} = $props();
</script>

{#if reviews}
	<div>
		{#each reviews as review (review.createdAt + review.rating)}
			<Review {review} {currentUser} {wayId} {onDeleted} />
		{/each}
	</div>
{/if}
