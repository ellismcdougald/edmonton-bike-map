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
	findRoute: vi.fn(() => Promise.resolve())
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
	reset: vi.fn()
};
const MockLeafletMap = vi.fn(() => mockMapInstance);
vi.mock('$lib/map/loadLeaflet', () => ({
	loadLeaflet: () => Promise.resolve({ LeafletMap: MockLeafletMap })
}));

import { findRoute } from '$lib/map/mapActions';

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

	it('calls findRoute with the map instance when Find Route button is clicked', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const findRouteButton = container.querySelector('#findRouteButton') as HTMLButtonElement;
		await fireEvent.click(findRouteButton);

		expect(findRoute).toHaveBeenCalledTimes(1);
		expect(findRoute).toHaveBeenCalledWith({ mapInstance: mockMapInstance });
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
});
