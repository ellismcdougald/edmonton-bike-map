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

vi.mock('$lib/map/LeafletMap', () => {
	const instance: any = {
		addTileLayer: vi.fn(),
		onMapClick: vi.fn(),
		loadInfoLayer: vi.fn(),
		showInfoLayer: vi.fn(),
		hideInfoLayer: vi.fn(),
		map: { remove: vi.fn() },
		addStartMarker: vi.fn(),
		removeStartMarker: vi.fn(),
		addEndMarker: vi.fn(),
		removeEndMarker: vi.fn()
	};

	const LeafletMap = vi.fn(() => instance);

	const loadLeaflet = vi.fn().mockResolvedValue({ LeafletMap });

	return { loadLeaflet };
});

import { findRoute } from '$lib/map/mapActions';

describe('Map.svelte', () => {
	let mapInstance: any;

	beforeEach(async () => {
		vi.clearAllMocks();
		// Import the mocked module so we can access instance
		const mocked = await import('$lib/map/LeafletMap');
		const { LeafletMap } = await mocked.loadLeaflet();
		mapInstance = new LeafletMap();
	});

	it('loads info layer on mount', async () => {
		render(Map);
		await waitFor(() => expect(mapInstance.addTileLayer).toHaveBeenCalled());
	});

	it('places a start marker when map is clicked in start mode', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mapInstance.addTileLayer).toHaveBeenCalled());

		const startButton = container.querySelector('#selectStartButton') as HTMLButtonElement;
		await fireEvent.click(startButton);

		// trigger the stored click handler
		const clickHandler = mapInstance.onMapClick.mock.calls[0][0];
		clickHandler([51, -0.1]);

		expect(mapInstance.removeStartMarker).toHaveBeenCalled();
		expect(mapInstance.addStartMarker).toHaveBeenCalledWith([51, -0.1]);
		expect(mapInstance.showInfoLayer).toHaveBeenCalled();
	});

	it('places an end marker when map is clicked in end mode', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mapInstance.addTileLayer).toHaveBeenCalled());

		const endButton = container.querySelector('#selectEndButton') as HTMLButtonElement;
		await fireEvent.click(endButton);

		// trigger the stored click handler
		const clickHandler = mapInstance.onMapClick.mock.calls[0][0];
		clickHandler([51, -0.1]);

		expect(mapInstance.removeEndMarker).toHaveBeenCalled();
		expect(mapInstance.addEndMarker).toHaveBeenCalledWith([51, -0.1]);
		expect(mapInstance.showInfoLayer).toHaveBeenCalled();
	});

	it('calls findRoute with the map instance when Find Route button is clicked', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mapInstance.addTileLayer).toHaveBeenCalled());

		const findRouteButton = container.querySelector('#findRouteButton') as HTMLButtonElement;
		await fireEvent.click(findRouteButton);

		expect(findRoute).toHaveBeenCalledTimes(1);
		expect(findRoute).toHaveBeenCalledWith({ mapInstance });
	});

	it('toggles select start/end modes correctly', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(mapInstance.addTileLayer).toHaveBeenCalled());

		const startButton = container.querySelector('#selectStartButton') as HTMLButtonElement;
		const endButton = container.querySelector('#selectEndButton') as HTMLButtonElement;

		expect(startButton.classList.contains('active')).toBe(false);
		expect(endButton.classList.contains('active')).toBe(false);

		await fireEvent.click(startButton);
		expect(mapInstance.hideInfoLayer).toHaveBeenCalled();

		await fireEvent.click(startButton);
		expect(mapInstance.showInfoLayer).toHaveBeenCalled();

		await fireEvent.click(endButton);
		expect(mapInstance.hideInfoLayer).toHaveBeenCalled();

		await fireEvent.click(endButton);
		expect(mapInstance.showInfoLayer).toHaveBeenCalled();
	});
});
