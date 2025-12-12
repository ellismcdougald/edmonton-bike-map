import { render, waitFor } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import ViewReviews from './ViewReviews.svelte';
import type { Review } from '$lib/types';

describe('ViewReviews.svelte', () => {
	it('renders a list of reviews', async () => {
		const reviews: Review[] = [
			{ wayId: 1, username: 'alice', rating: 7, comment: 'Nice!', createdAt: '2025-09-01' },
			{ wayId: 1, username: 'bob', rating: 5, comment: 'Ok', createdAt: '2025-09-02' }
		];

		const { getByText } = render(ViewReviews, { props: { reviews } });

		await waitFor(() => {
			expect(getByText(/alice/)).toBeTruthy();
			expect(getByText(/bob/)).toBeTruthy();
			expect(getByText('Nice!')).toBeTruthy();
			expect(getByText('Ok')).toBeTruthy();
			expect(getByText('Rating: 7 / 10')).toBeTruthy();
			expect(getByText('Rating: 5 / 10')).toBeTruthy();
		});
	});
});
