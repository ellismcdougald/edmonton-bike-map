import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import type { LeafletMap } from './LeafletMap';
import { findRoute, findRoutes } from './mapActions';

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

	it('does nothing if mapInstance is null', async () => {
		await findRoute({ mapInstance: null });

		expect(fetchSpy).not.toHaveBeenCalled();
		expect(alertSpy).not.toHaveBeenCalled();
	});

	it('alerts if only start point is missing', async () => {
		mapInstance.getStartLatLng = vi.fn().mockReturnValue(null);
		mapInstance.getEndLatLng = vi.fn().mockReturnValue([3, 4]);

		await findRoute({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			'Make sure you have selected both a start and an end point!'
		);
		expect(fetchSpy).not.toHaveBeenCalled();
	});

	it('alerts if only end point is missing', async () => {
		mapInstance.getStartLatLng = vi.fn().mockReturnValue([1, 2]);
		mapInstance.getEndLatLng = vi.fn().mockReturnValue(null);

		await findRoute({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			'Make sure you have selected both a start and an end point!'
		);
		expect(fetchSpy).not.toHaveBeenCalled();
	});

	it('handles network error during fetch', async () => {
		const networkError = new Error('Network failure');
		fetchSpy.mockRejectedValue(networkError);

		await findRoute({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith('Error fetching or displaying route: Network failure');
		expect(consoleErrorSpy).toHaveBeenCalledWith('Error fetching route:', networkError);
		expect(consoleErrorSpy).toHaveBeenCalledWith(networkError);
	});

	it('handles non-Error thrown values', async () => {
		const nonError = 'Some string error';
		fetchSpy.mockRejectedValue(nonError);

		await findRoute({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith('Error fetching or displaying route: Some string error');
		expect(consoleErrorSpy).toHaveBeenCalledWith(nonError);
	});

	it('constructs correct query parameters with various coordinate values', async () => {
		mapInstance.getStartLatLng = vi.fn().mockReturnValue([53.5461, -113.4938]);
		mapInstance.getEndLatLng = vi.fn().mockReturnValue([53.5234, -113.5234]);

		const mockGeojson = { type: 'FeatureCollection', features: [] };
		fetchSpy.mockResolvedValue({
			ok: true,
			status: 200,
			json: vi.fn().mockResolvedValue(mockGeojson)
		} as unknown as Response);

		await findRoute({ mapInstance });

		expect(fetchSpy).toHaveBeenCalledWith(
			'/api/route?startLatitude=53.5461&startLongitude=-113.4938&endLatitude=53.5234&endLongitude=-113.5234'
		);
	});

	it('removes existing route before loading new one', async () => {
		const mockGeojson = { type: 'FeatureCollection', features: [] };
		fetchSpy.mockResolvedValue({
			ok: true,
			status: 200,
			json: vi.fn().mockResolvedValue(mockGeojson)
		} as unknown as Response);

		await findRoute({ mapInstance });

		// Check that removeRouteLayer is called before loadRouteLayer
		const removeCall = (mapInstance.removeRouteLayer as ReturnType<typeof vi.fn>).mock
			.invocationCallOrder[0];
		const loadCall = (mapInstance.loadRouteLayer as ReturnType<typeof vi.fn>).mock
			.invocationCallOrder[0];

		expect(removeCall).toBeLessThan(loadCall);
	});
});

describe('findRoutes', () => {
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

	it('fetches and applies multiple routes to map with default k value', async () => {
		const mockGeojson = {
			type: 'FeatureCollection',
			features: [
				{ type: 'Feature', properties: { route_index: 1 }, geometry: {} },
				{ type: 'Feature', properties: { route_index: 2 }, geometry: {} },
				{ type: 'Feature', properties: { route_index: 3 }, geometry: {} }
			]
		};
		fetchSpy.mockResolvedValue({
			ok: true,
			status: 200,
			json: vi.fn().mockResolvedValue(mockGeojson)
		} as unknown as Response);

		await findRoutes({ mapInstance });

		expect(fetchSpy).toHaveBeenCalledWith(
			expect.stringContaining('/api/routes?startLatitude=1&startLongitude=2')
		);
		expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('k=3'));
		expect(mapInstance.removeRouteLayer).toHaveBeenCalled();
		expect(mapInstance.loadRouteLayer).toHaveBeenCalledWith(mockGeojson);
	});

	it('fetches routes with custom k value', async () => {
		const mockGeojson = {
			type: 'FeatureCollection',
			features: [
				{ type: 'Feature', properties: { route_index: 1 }, geometry: {} },
				{ type: 'Feature', properties: { route_index: 2 }, geometry: {} },
				{ type: 'Feature', properties: { route_index: 3 }, geometry: {} },
				{ type: 'Feature', properties: { route_index: 4 }, geometry: {} },
				{ type: 'Feature', properties: { route_index: 5 }, geometry: {} }
			]
		};
		fetchSpy.mockResolvedValue({
			ok: true,
			status: 200,
			json: vi.fn().mockResolvedValue(mockGeojson)
		} as unknown as Response);

		await findRoutes({ mapInstance, k: 5 });

		expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('k=5'));
	});

	it('alerts if start or end is missing', async () => {
		mapInstance.getStartLatLng = vi.fn().mockReturnValue(null);

		await findRoutes({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			'Make sure you have selected both a start and an end point!'
		);
		expect(fetchSpy).not.toHaveBeenCalled();
	});

	it('does nothing if mapInstance is null', async () => {
		await findRoutes({ mapInstance: null });

		expect(fetchSpy).not.toHaveBeenCalled();
		expect(alertSpy).not.toHaveBeenCalled();
	});

	it('alerts and logs if fetch fails', async () => {
		fetchSpy.mockResolvedValue({
			ok: false,
			status: 500,
			json: vi.fn()
		} as unknown as Response);

		await findRoutes({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			'Error fetching or displaying routes: Failed to get routes data'
		);
		expect(consoleErrorSpy).toHaveBeenCalled();
	});

	it('throws on 401 Unauthorized', async () => {
		fetchSpy.mockResolvedValue({
			ok: false,
			status: 401,
			json: vi.fn()
		} as unknown as Response);

		await expect(findRoutes({ mapInstance })).rejects.toThrow('Unauthorized');
	});

	it('alerts if only start point is missing', async () => {
		mapInstance.getStartLatLng = vi.fn().mockReturnValue(null);
		mapInstance.getEndLatLng = vi.fn().mockReturnValue([3, 4]);

		await findRoutes({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			'Make sure you have selected both a start and an end point!'
		);
		expect(fetchSpy).not.toHaveBeenCalled();
	});

	it('alerts if only end point is missing', async () => {
		mapInstance.getStartLatLng = vi.fn().mockReturnValue([1, 2]);
		mapInstance.getEndLatLng = vi.fn().mockReturnValue(null);

		await findRoutes({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			'Make sure you have selected both a start and an end point!'
		);
		expect(fetchSpy).not.toHaveBeenCalled();
	});

	it('handles network error during fetch', async () => {
		const networkError = new Error('Network failure');
		fetchSpy.mockRejectedValue(networkError);

		await findRoutes({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith('Error fetching or displaying routes: Network failure');
		expect(consoleErrorSpy).toHaveBeenCalledWith('Error fetching routes:', networkError);
		expect(consoleErrorSpy).toHaveBeenCalledWith(networkError);
	});

	it('handles non-Error thrown values', async () => {
		const nonError = { code: 'NETWORK_ERROR' };
		fetchSpy.mockRejectedValue(nonError);

		await findRoutes({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			expect.stringContaining('Error fetching or displaying routes:')
		);
		expect(consoleErrorSpy).toHaveBeenCalledWith(nonError);
	});

	it('constructs correct query parameters with various coordinate values', async () => {
		mapInstance.getStartLatLng = vi.fn().mockReturnValue([53.5461, -113.4938]);
		mapInstance.getEndLatLng = vi.fn().mockReturnValue([53.5234, -113.5234]);

		const mockGeojson = { type: 'FeatureCollection', features: [] };
		fetchSpy.mockResolvedValue({
			ok: true,
			status: 200,
			json: vi.fn().mockResolvedValue(mockGeojson)
		} as unknown as Response);

		await findRoutes({ mapInstance, k: 2 });

		expect(fetchSpy).toHaveBeenCalledWith(
			'/api/routes?startLatitude=53.5461&startLongitude=-113.4938&endLatitude=53.5234&endLongitude=-113.5234&k=2'
		);
	});

	it('removes existing routes before loading new ones', async () => {
		const mockGeojson = { type: 'FeatureCollection', features: [] };
		fetchSpy.mockResolvedValue({
			ok: true,
			status: 200,
			json: vi.fn().mockResolvedValue(mockGeojson)
		} as unknown as Response);

		await findRoutes({ mapInstance });

		// Check that removeRouteLayer is called before loadRouteLayer
		const removeCall = (mapInstance.removeRouteLayer as ReturnType<typeof vi.fn>).mock
			.invocationCallOrder[0];
		const loadCall = (mapInstance.loadRouteLayer as ReturnType<typeof vi.fn>).mock
			.invocationCallOrder[0];

		expect(removeCall).toBeLessThan(loadCall);
	});

	it('handles k=1 to find single route', async () => {
		const mockGeojson = {
			type: 'FeatureCollection',
			features: [{ type: 'Feature', properties: { route_index: 1 }, geometry: {} }]
		};
		fetchSpy.mockResolvedValue({
			ok: true,
			status: 200,
			json: vi.fn().mockResolvedValue(mockGeojson)
		} as unknown as Response);

		await findRoutes({ mapInstance, k: 1 });

		expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('k=1'));
		expect(mapInstance.loadRouteLayer).toHaveBeenCalledWith(mockGeojson);
	});
});
