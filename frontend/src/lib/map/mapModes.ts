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
