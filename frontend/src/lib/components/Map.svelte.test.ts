import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Map from './Map.svelte';

// Mock SvelteKit navigation
vi.mock('$app/navigation', () => ({
	goto: vi.fn(),
	beforeNavigate: vi.fn(() => () => {}), // returns unsubscribe fn
	afterNavigate: vi.fn(() => () => {})
}));

// Mock SvelteKit stores
vi.mock('$app/stores', () => {
	const subscribe = vi.fn(() => () => {}); // dummy unsubscribe
	return {
		page: { subscribe },
		navigating: { subscribe }
		// add more stores if needed
	};
});

vi.stubGlobal(
	'fetch',
	vi.fn().mockResolvedValue({
		ok: true,
		json: vi.fn().mockResolvedValue({ type: 'FeatureCollection', features: [] })
	})
);

vi.mock('$app/navigation', () => ({
	goto: vi.fn()
}));

vi.mock('$lib/map/mapActions', () => ({
	findRoutes: vi.fn(() => Promise.resolve())
}));

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const mockMapInstance: any = {
	addTileLayer: vi.fn(),
	onMapClick: vi.fn(),
	map: { remove: vi.fn() },
	addStartMarker: vi.fn(),
	removeStartMarker: vi.fn(),
	addEndMarker: vi.fn(),
	removeEndMarker: vi.fn(),
	reset: vi.fn(),
	loadSelectedWayLayer: vi.fn(),
	removeSelectedWayLayer: vi.fn(),
	loadAdjacentWaysLayer: vi.fn(),
	removeAdjacentWaysLayer: vi.fn(),
	// New layer for multi-way selection
	loadAdditionalSelectedWaysLayer: vi.fn(),
	removeAdditionalSelectedWaysLayer: vi.fn()
};
const MockLeafletMap = vi.fn(() => mockMapInstance);
vi.mock('$lib/map/loadLeaflet', () => ({
	loadLeaflet: () => Promise.resolve({ LeafletMap: MockLeafletMap })
}));

import { findRoutes } from '$lib/map/mapActions';

describe('Map.svelte', () => {
	beforeEach(async () => {
		vi.clearAllMocks();
	});

	it('loads the map on mount', async () => {
		render(Map);
		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());
		await waitFor(() => expect(mockMapInstance.onMapClick).toHaveBeenCalled());
	});

	it('places a start marker when map is clicked in start mode', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const startButton = container.querySelector('#selectStartButton') as HTMLButtonElement;
		await fireEvent.click(startButton);

		// trigger the stored click handler
		const clickHandler = mockMapInstance.onMapClick.mock.calls[0][0];
		await clickHandler([51, -0.1]);

		expect(mockMapInstance.removeStartMarker).toHaveBeenCalled();
		expect(mockMapInstance.addStartMarker).toHaveBeenCalledWith([51, -0.1]);
	});

	it('places an end marker when map is clicked in end mode', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const endButton = container.querySelector('#selectEndButton') as HTMLButtonElement;
		await fireEvent.click(endButton);

		// trigger the stored click handler
		const clickHandler = mockMapInstance.onMapClick.mock.calls[0][0];
		await clickHandler([51, -0.1]);

		expect(mockMapInstance.removeEndMarker).toHaveBeenCalled();
		expect(mockMapInstance.addEndMarker).toHaveBeenCalledWith([51, -0.1]);
	});

	it('calls findRoutes with the map instance and k=3 when Find Route button is clicked', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const findRouteButton = container.querySelector('#findRouteButton') as HTMLButtonElement;
		await fireEvent.click(findRouteButton);

		expect(findRoutes).toHaveBeenCalledTimes(1);
		expect(findRoutes).toHaveBeenCalledWith({ mapInstance: mockMapInstance, k: 3 });
	});

	it('toggles select start/end modes correctly', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const startButton = container.querySelector('#selectStartButton') as HTMLButtonElement;
		const endButton = container.querySelector('#selectEndButton') as HTMLButtonElement;

		expect(startButton.classList.contains('active')).toBe(false);
		expect(endButton.classList.contains('active')).toBe(false);

		// Click start button once to activate
		await fireEvent.click(startButton);
		expect(startButton.classList.contains('active')).toBe(true);

		// Click start button again to deactivate
		await fireEvent.click(startButton);
		expect(startButton.classList.contains('active')).toBe(false);

		// Click end button once to activate
		await fireEvent.click(endButton);
		expect(endButton.classList.contains('active')).toBe(true);

		// Click end button again to deactivate
		await fireEvent.click(endButton);
		expect(endButton.classList.contains('active')).toBe(false);
	});

	it('resets the map when Reset button is clicked', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const resetButton = container.querySelector('#resetButton') as HTMLButtonElement;
		await fireEvent.click(resetButton);

		expect(mockMapInstance.reset).toHaveBeenCalled();
	});

	it('fetches nearest way when map is clicked in normal mode', async () => {
		const mockFetch = vi.fn().mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ id: 12345, tags: { name: 'Test Street' } })
		});
		vi.stubGlobal('fetch', mockFetch);

		render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		// trigger the stored click handler in normal mode (not selecting markers)
		const clickHandler = mockMapInstance.onMapClick.mock.calls[0][0];
		await clickHandler([53.5461, -113.4938]);

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith(
				expect.stringContaining('/api/nearest-way?lat=53.5461&lng=-113.4938'),
				expect.objectContaining({ signal: expect.any(AbortSignal) })
			);
		});
	});

	it('loads adjacent ways layer when wayState.adjacentWays is populated', async () => {
		const { wayState } = await import('$lib/state.svelte');

		render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		// Set adjacent ways in state
		wayState.adjacentWays = [
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
		];

		await waitFor(() => {
			expect(mockMapInstance.loadAdjacentWaysLayer).toHaveBeenCalled();
			const callArg = mockMapInstance.loadAdjacentWaysLayer.mock.calls[0][0];
			expect(callArg.type).toBe('FeatureCollection');
			expect(callArg.features).toHaveLength(2);
			expect(callArg.features[0].properties.id).toBe(10);
			expect(callArg.features[1].properties.id).toBe(11);
		});

		// Cleanup
		wayState.adjacentWays = [];
	});

	it('handles Unauthorized error from findRoutes by redirecting to login', async () => {
		const { goto } = await import('$app/navigation');
		const mockFetch = vi.fn().mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({})
		});
		vi.stubGlobal('fetch', mockFetch);

		// Mock findRoutes to throw Unauthorized error
		const mockFindRoutes = vi.mocked(findRoutes);
		mockFindRoutes.mockRejectedValueOnce(new Error('Unauthorized'));

		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const findRouteButton = container.querySelector('#findRouteButton') as HTMLButtonElement;
		await fireEvent.click(findRouteButton);

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith('/logout', { method: 'POST' });
			expect(goto).toHaveBeenCalledWith('/login');
		});
	});

	it('silently handles non-Unauthorized errors from findRoutes', async () => {
		const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

		// Mock findRoutes to throw a non-Unauthorized error
		const mockFindRoutes = vi.mocked(findRoutes);
		mockFindRoutes.mockRejectedValueOnce(new Error('Network error'));

		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const findRouteButton = container.querySelector('#findRouteButton') as HTMLButtonElement;
		await fireEvent.click(findRouteButton);

		await waitFor(() => {
			expect(mockFindRoutes).toHaveBeenCalled();
		});

		// Error should be handled silently (no redirect)
		const { goto } = await import('$app/navigation');
		expect(goto).not.toHaveBeenCalled();

		consoleErrorSpy.mockRestore();
	});

	it('removes adjacent ways layer when wayState.adjacentWays is cleared', async () => {
		const { wayState } = await import('$lib/state.svelte');

		render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		// Set adjacent ways
		wayState.adjacentWays = [
			{
				type: 'Feature' as const,
				geometry: {
					type: 'LineString' as const,
					coordinates: [
						[1, 2],
						[3, 4]
					]
				},
				properties: { id: 10 }
			}
		];

		await waitFor(() => {
			expect(mockMapInstance.loadAdjacentWaysLayer).toHaveBeenCalled();
		});

		vi.clearAllMocks();

		// Clear adjacent ways
		wayState.adjacentWays = [];

		await waitFor(() => {
			expect(mockMapInstance.removeAdjacentWaysLayer).toHaveBeenCalled();
		});
	});

	it('loads additional selected ways layer when wayState.additionalSelectedWayIds is populated', async () => {
		const { wayState } = await import('$lib/state.svelte');

		render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		// Set selected way
		wayState.selectedWay = { id: 100, tags: { name: 'Main Street' } };

		// Set adjacent ways
		wayState.adjacentWays = [
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
		];

		// Mark way 10 as additionally selected
		wayState.additionalSelectedWayIds = [10];

		await waitFor(() => {
			expect(mockMapInstance.loadAdditionalSelectedWaysLayer).toHaveBeenCalled();
			const callArg = mockMapInstance.loadAdditionalSelectedWaysLayer.mock.calls[0][0];
			expect(callArg.type).toBe('FeatureCollection');
			expect(callArg.features).toHaveLength(1);
			expect(callArg.features[0].properties.id).toBe(10);
		});

		// Cleanup
		wayState.adjacentWays = [];
		wayState.additionalSelectedWayIds = [];
		wayState.selectedWay = null;
	});

	it('filters out additionally selected ways from adjacent ways layer', async () => {
		const { wayState } = await import('$lib/state.svelte');

		render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		// Set selected way
		wayState.selectedWay = { id: 100, tags: { name: 'Main Street' } };

		// Set adjacent ways
		wayState.adjacentWays = [
			{
				type: 'Feature' as const,
				geometry: {
					type: 'LineString' as const,
					coordinates: [
						[1, 2],
						[3, 4]
					]
				},
				properties: { id: 10, name: 'Way 1' }
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
				properties: { id: 11, name: 'Way 2' }
			},
			{
				type: 'Feature' as const,
				geometry: {
					type: 'LineString' as const,
					coordinates: [
						[9, 10],
						[11, 12]
					]
				},
				properties: { id: 12, name: 'Way 3' }
			}
		];

		// Mark way 10 as additionally selected
		wayState.additionalSelectedWayIds = [10];

		await waitFor(() => {
			expect(mockMapInstance.loadAdjacentWaysLayer).toHaveBeenCalled();
			// Only ways 11 and 12 should be in the adjacent ways layer (10 is selected)
			const callArg = mockMapInstance.loadAdjacentWaysLayer.mock.calls[0][0];
			expect(callArg.features).toHaveLength(2);
			expect(callArg.features[0].properties.id).toBe(11);
			expect(callArg.features[1].properties.id).toBe(12);
		});

		// Cleanup
		wayState.adjacentWays = [];
		wayState.additionalSelectedWayIds = [];
		wayState.selectedWay = null;
	});

	it('excludes the original selected way from adjacent ways layer', async () => {
		const { wayState } = await import('$lib/state.svelte');

		render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		// Set selected way
		wayState.selectedWay = { id: 100, tags: { name: 'Main Street' } };

		// Set adjacent ways including the selected way
		wayState.adjacentWays = [
			{
				type: 'Feature' as const,
				geometry: {
					type: 'LineString' as const,
					coordinates: [
						[1, 2],
						[3, 4]
					]
				},
				properties: { id: 100, name: 'Main Street' } // same as selected way
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
				properties: { id: 10, name: 'Adjacent Way' }
			}
		];

		await waitFor(() => {
			expect(mockMapInstance.loadAdjacentWaysLayer).toHaveBeenCalled();
			// Only way 10 should be in the adjacent ways layer (100 is the selected way)
			const callArg = mockMapInstance.loadAdjacentWaysLayer.mock.calls[0][0];
			expect(callArg.features).toHaveLength(1);
			expect(callArg.features[0].properties.id).toBe(10);
		});

		// Cleanup
		wayState.adjacentWays = [];
		wayState.selectedWay = null;
	});

	it('removes additional selected ways layer when additionalSelectedWayIds is cleared', async () => {
		const { wayState } = await import('$lib/state.svelte');

		render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		// Set selected way
		wayState.selectedWay = { id: 100, tags: { name: 'Main Street' } };

		// Set adjacent ways
		wayState.adjacentWays = [
			{
				type: 'Feature' as const,
				geometry: {
					type: 'LineString' as const,
					coordinates: [
						[1, 2],
						[3, 4]
					]
				},
				properties: { id: 10 }
			}
		];

		// Mark way 10 as additionally selected
		wayState.additionalSelectedWayIds = [10];

		await waitFor(() => {
			expect(mockMapInstance.loadAdditionalSelectedWaysLayer).toHaveBeenCalled();
		});

		vi.clearAllMocks();

		// Clear additional selected ways
		wayState.additionalSelectedWayIds = [];

		await waitFor(() => {
			expect(mockMapInstance.removeAdditionalSelectedWaysLayer).toHaveBeenCalled();
		});

		// Cleanup
		wayState.adjacentWays = [];
		wayState.selectedWay = null;
	});

	it('passes click handler to adjacent ways layer', async () => {
		const { wayState } = await import('$lib/state.svelte');

		render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		// Set up click handler
		const mockClickHandler = vi.fn();
		wayState.onAdjacentWayClick = mockClickHandler;

		// Set adjacent ways
		wayState.adjacentWays = [
			{
				type: 'Feature' as const,
				geometry: {
					type: 'LineString' as const,
					coordinates: [
						[1, 2],
						[3, 4]
					]
				},
				properties: { id: 10 }
			}
		];

		await waitFor(() => {
			expect(mockMapInstance.loadAdjacentWaysLayer).toHaveBeenCalled();
			const calls = mockMapInstance.loadAdjacentWaysLayer.mock.calls;
			// Second argument should be the click handler
			expect(calls[0][1]).toBe(mockClickHandler);
		});

		// Cleanup
		wayState.adjacentWays = [];
		wayState.onAdjacentWayClick = null;
	});
});
