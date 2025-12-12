/* eslint-disable @typescript-eslint/no-explicit-any */
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

	const geoJsonCalls: any[] = [];
	const geoJSON = (geojson: any, options: any) => {
		const layer: any = {
			_options: options,
			addTo(map: any) {
				if (!map._layers) map._layers = [];
				map._layers.push(layer);
				return layer;
			},
			getBounds() {
				return {
					isValid() {
						return true;
					}
				};
			}
		};
		geoJsonCalls.push({ geojson, options, layer });
		return layer;
	};

	const tileLayer = (_url: string, _opts: any) => ({ addTo: (_map: any) => ({}) });

	const mapFactory = (_selector: string, _opts: any) => {
		const root = document.getElementById('map') || document.body;
		const m: any = {
			_rootEl: root,
			_eventHandlers: {} as any,
			_layers: [] as any[],
			createPane: (_name: string) => {},
			getPane: (_name: string) => ({ style: {} }),
			fitBounds: (_b: any) => {},
			getContainer: () => root,
			on: (event: string, handler: any) => {
				if (!m._eventHandlers[event]) m._eventHandlers[event] = [];
				m._eventHandlers[event].push(handler);
			},
			removeLayer: (layer: any) => {
				m._layers = m._layers.filter((l: any) => l !== layer);
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
		map: mapFactory,
		_geoJsonCalls: geoJsonCalls
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

describe('LeafletMap.loadSelectedWayLayer', () => {
	beforeEach(() => {
		let el = document.getElementById('map');
		if (!el) {
			el = document.createElement('div');
			el.id = 'map';
			document.body.appendChild(el);
		}
	});

	it('renders highlight with visible style and replaces existing layer', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const featureA: any = {
			type: 'Feature',
			properties: {},
			geometry: { type: 'LineString', coordinates: [] }
		};
		const featureB: any = {
			type: 'Feature',
			properties: {},
			geometry: { type: 'LineString', coordinates: [] }
		};

		lm.loadSelectedWayLayer(featureA);
		lm.loadSelectedWayLayer(featureB);

		const calls = (L as any)._geoJsonCalls;
		expect(calls.length).toBe(2);
		expect(calls[1].options.style.color).toBe('#00e5ff');
		expect(calls[1].options.style.weight).toBe(6);
		expect(calls[1].options.style.opacity).toBe(0.9);
	});
});
