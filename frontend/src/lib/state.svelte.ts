import type { WayFeature } from './types';

export const wayState = $state<{
	selectedWay: WayFeature | null;
}>({
	selectedWay: null
});
