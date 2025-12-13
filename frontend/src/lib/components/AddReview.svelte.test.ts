import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AddReview from './AddReview.svelte';
import * as reviewsApi from '$lib/api/client/reviews';
import * as waysApi from '$lib/api/client/ways';
import type { WayFeature } from '$lib/types';

// Provide a hoisted mock for wayState without referencing test-scoped vars
vi.mock('$lib/state.svelte', () => {
	const wayState: {
		selectedWay: WayFeature | null;
		adjacentWays: unknown[];
		onAdjacentWayClick: ((wayId: number) => void) | null;
	} = {
		selectedWay: null,
		adjacentWays: [],
		onAdjacentWayClick: null
	};
	return { wayState };
});
import { wayState } from '$lib/state.svelte';

describe('AddReview.svelte', () => {
	const wayId = 42;

	beforeEach(() => {
		vi.clearAllMocks();
		// Mock getAdjacentWays to return empty array
		vi.spyOn(waysApi, 'getAdjacentWays').mockResolvedValue({
			type: 'FeatureCollection',
			features: []
		});
	});

	it('renders Included Ways pill using selectedWay name when available', async () => {
		// set selectedWay via the mocked state
		wayState.selectedWay = { id: wayId, tags: { name: 'Sample Way' } } as WayFeature;

		const { container, getByText } = render(AddReview, { props: { wayId } });
		await waitFor(() => {
			expect(getByText('Included Ways')).toBeTruthy();
			const pill = container.querySelector('#includedWays');
			expect(pill?.textContent).toContain('Sample Way');
		});
	});

	it('submits rating/comment, calls API and onSubmitted', async () => {
		const createReviewSpy = vi.spyOn(reviewsApi, 'createReview').mockResolvedValue();
		const onSubmitted = vi.fn();

		// provide a basic selectedWay
		wayState.selectedWay = { id: wayId, tags: {} } as WayFeature;

		const { container } = render(AddReview, { props: { wayId, onSubmitted } });

		const ratingInput = container.querySelector<HTMLInputElement>('#rating')!;
		const commentInput = container.querySelector<HTMLTextAreaElement>('#comment')!;
		const submitButton = container.querySelector<HTMLButtonElement>('#submitButton')!;

		ratingInput.value = '9';
		await fireEvent.input(ratingInput);
		commentInput.value = 'Love this way!';
		await fireEvent.input(commentInput);

		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(createReviewSpy).toHaveBeenCalledWith([wayId], 9, 'Love this way!');
			expect(onSubmitted).toHaveBeenCalled();
		});
	});

	it('fetches adjacent ways on mount and updates state', async () => {
		const mockAdjacentWays = {
			type: 'FeatureCollection' as const,
			features: [
				{
					type: 'Feature' as const,
					geometry: {
						type: 'LineString' as const,
						coordinates: [
							[1, 2],
							[3, 4]
						]
					},
					properties: { id: 10, name: 'Adjacent Way 1' }
				},
				{
					type: 'Feature' as const,
					geometry: {
						type: 'LineString' as const,
						coordinates: [
							[5, 6],
							[7, 8]
						]
					},
					properties: { id: 11, name: 'Adjacent Way 2' }
				}
			]
		};

		const getAdjacentWaysSpy = vi
			.spyOn(waysApi, 'getAdjacentWays')
			.mockResolvedValue(mockAdjacentWays);
		wayState.selectedWay = { id: wayId, tags: {} } as WayFeature;

		render(AddReview, { props: { wayId } });

		await waitFor(() => {
			expect(getAdjacentWaysSpy).toHaveBeenCalledWith(wayId);
			expect(wayState.adjacentWays).toEqual(mockAdjacentWays.features);
		});
	});

	it('clears adjacent ways state on unmount', async () => {
		wayState.selectedWay = { id: wayId, tags: {} } as WayFeature;
		wayState.adjacentWays = [
			{
				type: 'Feature' as const,
				geometry: { type: 'LineString' as const, coordinates: [[1, 2]] },
				properties: { id: 10 }
			}
		];

		const { unmount } = render(AddReview, { props: { wayId } });

		expect(wayState.adjacentWays.length).toBeGreaterThan(0);

		unmount();

		expect(wayState.adjacentWays).toEqual([]);
	});

	it('refetches adjacent ways when wayId changes', async () => {
		const mockAdjacentWays1 = {
			type: 'FeatureCollection' as const,
			features: [
				{
					type: 'Feature' as const,
					geometry: { type: 'LineString' as const, coordinates: [[1, 2]] },
					properties: { id: 10 }
				}
			]
		};

		const mockAdjacentWays2 = {
			type: 'FeatureCollection' as const,
			features: [
				{
					type: 'Feature' as const,
					geometry: { type: 'LineString' as const, coordinates: [[3, 4]] },
					properties: { id: 20 }
				}
			]
		};

		const getAdjacentWaysSpy = vi
			.spyOn(waysApi, 'getAdjacentWays')
			.mockResolvedValueOnce(mockAdjacentWays1)
			.mockResolvedValueOnce(mockAdjacentWays2);

		wayState.selectedWay = { id: wayId, tags: {} } as WayFeature;

		const { rerender } = render(AddReview, { props: { wayId: 42 } });

		await waitFor(() => {
			expect(getAdjacentWaysSpy).toHaveBeenCalledWith(42);
		});

		// Change wayId
		await rerender({ wayId: 99 });

		await waitFor(() => {
			expect(getAdjacentWaysSpy).toHaveBeenCalledWith(99);
			expect(getAdjacentWaysSpy).toHaveBeenCalledTimes(2);
		});
	});
});
