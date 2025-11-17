import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import type { LeafletMap } from './LeafletMap';
import { findRoute } from './mapActions';

describe('findRoute', () => {
	let mapInstance: LeafletMap;
	let alertSpy: ReturnType<typeof vi.fn>;
	let consoleErrorSpy: ReturnType<typeof vi.spyOn>;

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
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('fetches and applies route to map', async () => {
		globalThis.localStorage?.setItem?.('token', 'jwt-abc');
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
		// set token so fetch is attempted
		globalThis.localStorage?.setItem?.('token', 'jwt-abc');
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

	it('errors immediately when no auth token present', async () => {
		// ensure token is not present
		globalThis.localStorage?.removeItem?.('token');

		await findRoute({ mapInstance });

		expect(alertSpy).toHaveBeenCalledWith(
			'Error fetching or displaying route: Missing auth token: please log in before requesting a route'
		);
	});

	it('sends Authorization header when token exists', async () => {
		// set token
		globalThis.localStorage?.setItem?.('token', 'jwt-abc');

		const mockFetch = vi.fn().mockResolvedValue({
			ok: true,
			json: vi.fn().mockResolvedValue({ type: 'FeatureCollection', features: [] })
		} as unknown as Response);

		await findRoute({ mapInstance, fetchFn: mockFetch });

		expect(mockFetch).toHaveBeenCalled();
		const calledWith = mockFetch.mock.calls[0];
		// second arg is options
		expect(calledWith[1]).toBeDefined();
		expect(calledWith[1].headers.Authorization).toBe('Bearer jwt-abc');
	});
});
