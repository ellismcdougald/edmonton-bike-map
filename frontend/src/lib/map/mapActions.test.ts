import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import type { LeafletMap } from './LeafletMap';
import { findRoute } from './mapActions';

describe('findRoute', () => {
	let mapInstance: LeafletMap;
	let alertSpy: ReturnType<typeof vi.fn>;
	let consoleErrorSpy: ReturnType<typeof vi.spyOn>;

	beforeEach(() => {
		// Mock mapInstance
		mapInstance = {
			getStartLatLng: vi.fn().mockReturnValue([1, 2] as [number, number]),
			getEndLatLng: vi.fn().mockReturnValue([3, 4] as [number, number]),
			removeRouteLayer: vi.fn(),
			loadRouteLayer: vi.fn()
		} as unknown as LeafletMap;

		// Mock global alert
		alertSpy = vi.fn();
		globalThis.alert = alertSpy;

		// Mock console.error
		consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('fetches and applies route to map', async () => {
		const mockGeojson = { type: 'FeatureCollection', features: [] };
		const mockFetch = vi.fn().mockResolvedValue({
			ok: true,
			json: vi.fn().mockResolvedValue(mockGeojson)
		} as unknown as Response);

		await findRoute({ mapInstance, fetchFn: mockFetch });

		expect(mockFetch).toHaveBeenCalled();
		expect(mapInstance.removeRouteLayer).toHaveBeenCalled();
		expect(mapInstance.loadRouteLayer).toHaveBeenCalledWith(mockGeojson);
	});

	it('alerts if start or end is missing', async () => {
		mapInstance.getStartLatLng = vi.fn().mockReturnValue(null);

		await findRoute({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			'Make sure you have selected both a start and an end point!'
		);
	});

	it('alerts and logs if fetch fails', async () => {
		const mockFetch = vi.fn().mockResolvedValue({
			ok: false,
			json: vi.fn()
		} as unknown as Response);

		await findRoute({ mapInstance, fetchFn: mockFetch });

		expect(alertSpy).toHaveBeenCalledWith(
			'Error fetching or displaying route: Failed to get route data'
		);
		expect(consoleErrorSpy).toHaveBeenCalled();
	});
});
