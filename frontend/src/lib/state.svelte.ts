import type { WayFeature, WayFeatureGeoJSON } from './types';

export const wayState = $state<{
	selectedWay: WayFeature | null;
	adjacentWays: WayFeatureGeoJSON[];
	additionalSelectedWayIds: number[];
	onAdjacentWayClick: ((wayId: number) => void) | null;
	isAddReviewActive: boolean;
}>({
	selectedWay: null,
	adjacentWays: [],
	additionalSelectedWayIds: [],
	onAdjacentWayClick: null,
	isAddReviewActive: false
});
