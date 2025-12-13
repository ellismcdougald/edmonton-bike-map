import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getAdjacentWays } from './ways';

describe('ways API client', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		global.fetch = vi.fn();
	});

	describe('getAdjacentWays', () => {
		it('fetches adjacent ways successfully', async () => {
			const mockResponse = {
				type: 'FeatureCollection',
				features: [
					{
						type: 'Feature',
						geometry: {
							type: 'LineString',
							coordinates: [
								[1, 2],
								[3, 4]
							]
						},
						properties: { id: 10, name: 'Adjacent Way 1' }
					},
					{
						type: 'Feature',
						geometry: {
							type: 'LineString',
							coordinates: [
								[5, 6],
								[7, 8]
							]
						},
						properties: { id: 11, name: 'Adjacent Way 2' }
					}
				]
			};

			(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
				ok: true,
				json: async () => mockResponse
			});

			const result = await getAdjacentWays(42);

			expect(global.fetch).toHaveBeenCalledWith('/api/adjacent-ways?id=42');
			expect(result).toEqual(mockResponse);
			expect(result.features).toHaveLength(2);
		});

		it('returns empty FeatureCollection when no adjacent ways exist', async () => {
			const mockResponse = {
				type: 'FeatureCollection',
				features: []
			};

			(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
				ok: true,
				json: async () => mockResponse
			});

			const result = await getAdjacentWays(99);

			expect(result.features).toHaveLength(0);
		});

		it('throws error when fetch fails with 404', async () => {
			(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
				ok: false,
				statusText: 'Not Found'
			});

			await expect(getAdjacentWays(999)).rejects.toThrow(
				'Failed to fetch adjacent ways: Not Found'
			);
		});

		it('throws error when fetch fails with 500', async () => {
			(global.fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
				ok: false,
				statusText: 'Internal Server Error'
			});

			await expect(getAdjacentWays(42)).rejects.toThrow(
				'Failed to fetch adjacent ways: Internal Server Error'
			);
		});
	});
});
