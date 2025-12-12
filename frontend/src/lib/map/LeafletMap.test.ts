/* eslint-disable @typescript-eslint/no-explicit-any, @typescript-eslint/no-unused-vars */
import { describe, it, expect, beforeEach } from 'vitest';
import { LeafletMap } from './LeafletMap';

function makeFakeL() {
	// Minimal fake Leaflet-like API used by LeafletMap
	const DomUtil = {
		create: (tag: string, _className?: string) => {
			const el = document.createElement(tag);
			if (_className) el.className = _className;
			return el;
		}
	};

	const Control = {
		extend: (proto: any) => {
			return class {
				proto = proto;
				opts: any;
				_el: HTMLElement | null = null;
				constructor(opts: any) {
					this.opts = opts;
				}
				addTo(map: any) {
					if (this.proto && typeof this.proto.onAdd === 'function') {
						this._el = this.proto.onAdd();
						// attach to DOM so tests can query
						if (map && map._rootEl) map._rootEl.appendChild(this._el);
					}
					return this;
				}
				remove() {
					if (this._el && this._el.parentNode) this._el.parentNode.removeChild(this._el);
				}
			};
		}
	};

	const geoJSON = (geojson: any, _options: any) => {
		return {
			addTo(map: any) {
				// simple object representing a GeoJSON layer
				return {
					getBounds() {
						return {
							isValid() {
								return true;
							}
						};
					}
				};
			}
		};
	};

	const tileLayer = (_url: string, _opts: any) => ({ addTo: (_map: any) => ({}) });

	const mapFactory = (_selector: string, _opts: any) => {
		const root = document.getElementById('map') || document.body;
		const m: any = {
			_rootEl: root,
			_eventHandlers: {} as any,
			createPane: (_name: string) => {},
			getPane: (_name: string) => ({ style: {} }),
			fitBounds: (_b: any) => {},
			getContainer: () => root,
			on: (event: string, handler: any) => {
				if (!m._eventHandlers[event]) m._eventHandlers[event] = [];
				m._eventHandlers[event].push(handler);
			},
			// setView should return the map object so chaining works
			setView: (_v: any, _z: any) => m
		};
		return m;
	};

	return {
		DomUtil,
		Control,
		geoJSON,
		tileLayer,
		map: mapFactory
	} as any;
}

describe('LeafletMap.loadRouteLayer', () => {
	beforeEach(() => {
		// ensure a map container exists
		let el = document.getElementById('map');
		if (!el) {
			el = document.createElement('div');
			el.id = 'map';
			document.body.appendChild(el);
		}
	});

	it('uses backend time_minutes when present', () => {
		const L = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'Feature',
			properties: {
				distance_km: 5,
				time_minutes: 15
			},
			geometry: { type: 'LineString', coordinates: [] }
		};

		lm.loadRouteLayer(geojson as any);

		const dc = (lm as any).distanceControl as any;
		const tc = (lm as any).timeControl as any;
		expect(dc).toBeDefined();
		expect(tc).toBeDefined();
		expect(dc._el).toBeInstanceOf(HTMLElement);
		expect(tc._el).toBeInstanceOf(HTMLElement);
		expect(dc._el.innerHTML).toContain('Distance: 5.00 km');
		expect(tc._el.innerHTML).toContain('Estimated time: 15 min');
	});

	it('falls back to estimated time when time_minutes missing', () => {
		const L = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'Feature',
			properties: {
				distance_km: 10
			},
			geometry: { type: 'LineString', coordinates: [] }
		};

		lm.loadRouteLayer(geojson as any);

		const tc = (lm as any).timeControl as any;
		expect(tc).toBeDefined();
		expect(tc._el).toBeInstanceOf(HTMLElement);
		// default avg speed used in fallback is 20 km/h -> 10km = 30min
		expect(tc._el.innerHTML).toContain('Estimated time: 30 min');
	});
});
