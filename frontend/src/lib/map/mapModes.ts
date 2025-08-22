/**
 * mapModes.ts
 *
 * Purpose:
 * Manages the state of map interaction modes (selecting start and end points)
 * for the LeafletMap component.
 *
 * Types:
 * - MapModeState: object with boolean flags `selectStartActive` and `selectEndActive`
 *
 * Behavior:
 * - toggleSelectStart(state): toggles the `selectStartActive` flag
 *     - If activating start selection, ensures end selection is deactivated
 * - toggleSelectEnd(state): toggles the `selectEndActive` flag
 *     - If activating end selection, ensures start selection is deactivated
 *
 * Notes:
 * - Pure functions; return new state objects without mutating input
 * - Used by map click handlers and sidebar buttons to manage selection mode
 */

export type MapModeState = {
	selectStartActive: boolean;
	selectEndActive: boolean;
};

export function toggleSelectStart(state: MapModeState): MapModeState {
	return state.selectStartActive
		? { ...state, selectStartActive: false }
		: { selectStartActive: true, selectEndActive: false };
}

export function toggleSelectEnd(state: MapModeState): MapModeState {
	return state.selectEndActive
		? { ...state, selectEndActive: false }
		: { selectStartActive: false, selectEndActive: true };
}
