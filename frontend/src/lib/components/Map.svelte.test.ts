import { fireEvent, render, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Map from './Map.svelte';
import { findRoute } from '$lib/map/mapActions';

vi.mock('$lib/map/LeafletMap', () => {
	let instance: any;

	const LeafletMap = vi.fn().mockImplementation(() => {
		instance = {
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
		return instance;
	});

	const __getInstance = () => instance;

	return { LeafletMap, __getInstance };
});

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

describe('Map.svelte', () => {
	let getInstance: () => any;

	beforeEach(async () => {
		vi.clearAllMocks();
		const mocked = await vi.importMock<typeof import('$lib/map/LeafletMap')>('$lib/map/LeafletMap');
		// @ts-expect-error exists only in mock
		getInstance = mocked.__getInstance;
	});

	it('loads info layer on mount', async () => {
		render(Map);

		await waitFor(() => {
			expect(getInstance().loadInfoLayer).toHaveBeenCalled();
		});
	});
});

describe('Select Start / Select End', () => {
	let getInstance: () => any;

	beforeEach(async () => {
		vi.clearAllMocks();
		const mocked = await vi.importMock<typeof import('$lib/map/LeafletMap')>('$lib/map/LeafletMap');
		// @ts-expect-error exists only in mock
		getInstance = mocked.__getInstance;
	});

	it('places a start marker when map is clicked in start mode', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(getInstance().addTileLayer).toHaveBeenCalled());

		const startButton = container.querySelector('#selectStartButton') as HTMLButtonElement;
		await fireEvent.click(startButton);

		const mapInstance = getInstance();
		const mapClickHandler = mapInstance.onMapClick.mock.calls[0][0];

		mapClickHandler([51, -0.1]);

		expect(mapInstance.removeStartMarker).toHaveBeenCalled();
		expect(mapInstance.addStartMarker).toHaveBeenCalledWith([51, -0.1]);
		expect(mapInstance.showInfoLayer).toHaveBeenCalled();
	});

	it('places an end marker when map is clicked in end mode', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(getInstance()).toBeTruthy());

		const endButton = container.querySelector('#selectEndButton') as HTMLButtonElement;
		await fireEvent.click(endButton);

		const mapInstance = getInstance();
		const mapClickHandler = mapInstance.onMapClick.mock.calls[0][0];

		mapClickHandler([51, -0.1]);

		expect(mapInstance.removeEndMarker).toHaveBeenCalled();
		expect(mapInstance.addEndMarker).toHaveBeenCalledWith([51, -0.1]);
		expect(mapInstance.showInfoLayer).toHaveBeenCalled();
	});
});

describe('Find Route', () => {
	let getInstance: () => any;

	beforeEach(async () => {
		vi.clearAllMocks();
		const mocked = await vi.importMock<typeof import('$lib/map/LeafletMap')>('$lib/map/LeafletMap');
		// @ts-expect-error exists only in mock
		getInstance = mocked.__getInstance;
	});

	it('calls findRoute with the map instance when Find Route is clicked', async () => {
		const { container } = render(Map);

		await waitFor(() => expect(getInstance()).toBeTruthy());
		const mapInstance = getInstance();

		const findRouteButton = container.querySelector('#findRouteButton') as HTMLButtonElement;
		await fireEvent.click(findRouteButton);

		expect(findRoute).toHaveBeenCalledTimes(1);
		expect(findRoute).toHaveBeenCalledWith({ mapInstance });
	});
});

describe('Reset', () => {
	beforeEach(async () => {
		vi.clearAllMocks();
		const mocked = await vi.importMock<typeof import('$lib/map/LeafletMap')>('$lib/map/LeafletMap');
		// @ts-expect-error exists only in mock
		getInstance = mocked.__getInstance;
	});

	it('renders the Reset button and it is clickable', async () => {
		const { container } = render(Map);

		const resetButton = container.querySelector('#resetButton') as HTMLButtonElement;
		expect(resetButton).toBeTruthy();

		await fireEvent.click(resetButton);

		// TODO: add assertions once reset behaviour is implemented
	});
});
