import type { WayFeature, WayFeatureGeoJSON } from './types';

export const wayState = $state<{
	selectedWay: WayFeature | null;
	adjacentWays: WayFeatureGeoJSON[];
	onAdjacentWayClick: ((wayId: number) => void) | null;
}>({
	selectedWay: null,
	adjacentWays: [],
	onAdjacentWayClick: null
});
