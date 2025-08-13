import { describe, it, expect } from 'vitest';
import { toggleSelectStart, toggleSelectEnd, type MapModeState } from './mapModes';

describe('MapModeState toggles', () => {
	it('toggleSelectStart should activate start and deactivate end', () => {
		const state: MapModeState = { selectStartActive: false, selectEndActive: true };
		const newState = toggleSelectStart(state);
		expect(newState).toEqual({ selectStartActive: true, selectEndActive: false });
	});

	it('toggleSelectStart should deactivate start if already active', () => {
		const state: MapModeState = { selectStartActive: true, selectEndActive: false };
		const newState = toggleSelectStart(state);
		expect(newState).toEqual({ selectStartActive: false, selectEndActive: false });
	});

	it('toggleSelectEnd should activate end and deactivate start', () => {
		const state: MapModeState = { selectStartActive: true, selectEndActive: false };
		const newState = toggleSelectEnd(state);
		expect(newState).toEqual({ selectStartActive: false, selectEndActive: true });
	});

	it('toggleSelectEnd should deactivate end if already active', () => {
		const state: MapModeState = { selectStartActive: false, selectEndActive: true };
		const newState = toggleSelectEnd(state);
		expect(newState).toEqual({ selectStartActive: false, selectEndActive: false });
	});
});
