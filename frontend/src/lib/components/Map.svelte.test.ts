import { render, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Map from './Map.svelte';

vi.stubGlobal(
	'fetch',
	vi.fn().mockResolvedValue({
		ok: true,
		json: vi.fn().mockResolvedValue({ type: 'FeatureCollection', features: [] })
	})
);

vi.mock('$lib/map/mapActions', () => ({
	findRoute: vi.fn()
}));

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const mockMapInstance: any = {
	addTileLayer: vi.fn(),
	onMapClick: vi.fn(),
	loadInfoLayer: vi.fn(),
	showInfoLayer: vi.fn(),
	hideInfoLayer: vi.fn(),
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

	it('loads info layer on mount', async () => {
		render(Map);
		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());
	});

	it('places a start marker when map is clicked in start mode', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const startButton = container.querySelector('#selectStartButton') as HTMLButtonElement;
		await fireEvent.click(startButton);

		// trigger the stored click handler
		const clickHandler = mockMapInstance.onMapClick.mock.calls[0][0];
		clickHandler([51, -0.1]);

		expect(mockMapInstance.removeStartMarker).toHaveBeenCalled();
		expect(mockMapInstance.addStartMarker).toHaveBeenCalledWith([51, -0.1]);
		expect(mockMapInstance.showInfoLayer).toHaveBeenCalled();
	});

	it('places an end marker when map is clicked in end mode', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const endButton = container.querySelector('#selectEndButton') as HTMLButtonElement;
		await fireEvent.click(endButton);

		// trigger the stored click handler
		const clickHandler = mockMapInstance.onMapClick.mock.calls[0][0];
		clickHandler([51, -0.1]);

		expect(mockMapInstance.removeEndMarker).toHaveBeenCalled();
		expect(mockMapInstance.addEndMarker).toHaveBeenCalledWith([51, -0.1]);
		expect(mockMapInstance.showInfoLayer).toHaveBeenCalled();
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

		await fireEvent.click(startButton);
		expect(mockMapInstance.hideInfoLayer).toHaveBeenCalled();

		await fireEvent.click(startButton);
		expect(mockMapInstance.showInfoLayer).toHaveBeenCalled();

		await fireEvent.click(endButton);
		expect(mockMapInstance.hideInfoLayer).toHaveBeenCalled();

		await fireEvent.click(endButton);
		expect(mockMapInstance.showInfoLayer).toHaveBeenCalled();
	});

	it('resets the map when Reset button is clicked', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mockMapInstance.addTileLayer).toHaveBeenCalled());

		const resetButton = container.querySelector('#resetButton') as HTMLButtonElement;
		await fireEvent.click(resetButton);

		expect(mockMapInstance.reset).toHaveBeenCalled();
	});
});
