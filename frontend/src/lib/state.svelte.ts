import type { WayFeature, WayFeatureGeoJSON } from './types';

export const wayState = $state<{
	selectedWay: WayFeature | null;
	adjacentWays: WayFeatureGeoJSON[];
}>({
	selectedWay: null,
	adjacentWays: []
});
