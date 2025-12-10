import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import type { LeafletMap } from './LeafletMap';
import { findRoute } from './mapActions';

describe('findRoute', () => {
	let mapInstance: LeafletMap;
	let alertSpy: ReturnType<typeof vi.fn>;
	let consoleErrorSpy: ReturnType<typeof vi.spyOn>;
	let fetchSpy: ReturnType<typeof vi.spyOn>;

	beforeEach(() => {
		mapInstance = {
			getStartLatLng: vi.fn().mockReturnValue([1, 2] as [number, number]),
			getEndLatLng: vi.fn().mockReturnValue([3, 4] as [number, number]),
			removeRouteLayer: vi.fn(),
			loadRouteLayer: vi.fn()
		} as unknown as LeafletMap;

		alertSpy = vi.fn();
		globalThis.alert = alertSpy;

		consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
		fetchSpy = vi.spyOn(globalThis, 'fetch') as unknown as ReturnType<typeof vi.spyOn>;
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('fetches and applies route to map', async () => {
		const mockGeojson = { type: 'FeatureCollection', features: [] };
		fetchSpy.mockResolvedValue({
			ok: true,
			status: 200,
			json: vi.fn().mockResolvedValue(mockGeojson)
		} as unknown as Response);

		await findRoute({ mapInstance });

		expect(fetchSpy).toHaveBeenCalledWith(
			expect.stringContaining('/api/route?startLatitude=1&startLongitude=2')
		);
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
		fetchSpy.mockResolvedValue({
			ok: false,
			status: 500,
			json: vi.fn()
		} as unknown as Response);

		await findRoute({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			'Error fetching or displaying route: Failed to get route data'
		);
		expect(consoleErrorSpy).toHaveBeenCalled();
	});

	it('throws on 401 Unauthorized', async () => {
		fetchSpy.mockResolvedValue({
			ok: false,
			status: 401,
			json: vi.fn()
		} as unknown as Response);

		await expect(findRoute({ mapInstance })).rejects.toThrow('Unauthorized');
	});
});
