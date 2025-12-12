<script lang="ts">
	import Review from './Review.svelte';
	import type { Review as ReviewObj } from '$lib/types';

	let { reviews }: { reviews: ReviewObj[] } = $props();
</script>

{#if reviews}
	<div>
		{#each reviews as review (review.createdAt + review.rating)}
			<Review {review} />
		{/each}
	</div>
{/if}
