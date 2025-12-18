/* eslint-disable @typescript-eslint/no-explicit-any */
import { describe, it, expect, beforeEach, vi } from 'vitest';
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

	const DomEvent = {
		stopPropagation: (_e: any) => {}
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
			_geojson: geojson,
			_options: options,
			_features: [],
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

		// Process features to capture event handlers
		if (geojson && geojson.features) {
			geojson.features.forEach((feature: any) => {
				const mockLayer = {
					on: (event: string, handler: any) => {
						if (!mockLayer._events) mockLayer._events = {};
						if (!mockLayer._events[event]) mockLayer._events[event] = [];
						mockLayer._events[event].push(handler);
					},
					bindPopup: (content: string) => {
						mockLayer._popupContent = content;
						return mockLayer;
					},
					_events: {} as any,
					_popupContent: ''
				};

				if (options.onEachFeature) {
					options.onEachFeature(feature, mockLayer);
				}

				layer._features.push({ feature, layer: mockLayer });
			});
		}

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
			hasLayer: (layer: any) => m._layers.includes(layer),
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
		DomEvent,
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

	it('splits routes into shortest (blue) and other (green) layers', () => {
		const L = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { route_index: 1, distance_km: 5, time_minutes: 15 },
					geometry: { type: 'LineString', coordinates: [] }
				},
				{
					type: 'Feature',
					properties: { route_index: 2, distance_km: 6, time_minutes: 18 },
					geometry: { type: 'LineString', coordinates: [] }
				},
				{
					type: 'Feature',
					properties: { route_index: 3, distance_km: 7, time_minutes: 20 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadRouteLayer(geojson);

		const calls = (L as any)._geoJsonCalls;
		expect(calls.length).toBe(2);

		// First call should be shortest route (blue)
		expect(calls[0].geojson.features.length).toBe(1);
		expect(calls[0].geojson.features[0].properties.route_index).toBe(1);
		expect(calls[0].options.style.color).toBe('blue');
		expect(calls[0].options.style.weight).toBe(5);

		// Second call should be other routes (green)
		expect(calls[1].geojson.features.length).toBe(2);
		expect(calls[1].options.style.color).toBe('green');
		expect(calls[1].options.style.weight).toBe(5);
	});

	it('creates correct popup content for shortest route', () => {
		const L = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { route_index: 1, distance_km: 5.25, time_minutes: 15 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadRouteLayer(geojson);

		const calls = (L as any)._geoJsonCalls;
		const featureLayer = calls[0].layer._features[0].layer;
		expect(featureLayer._popupContent).toContain('<strong>Best Route</strong>');
		expect(featureLayer._popupContent).toContain('Distance: 5.25 km');
		expect(featureLayer._popupContent).toContain('Time: 15 min');
	});

	it('creates correct popup content for alternate routes', () => {
		const L = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { route_index: 1, distance_km: 5, time_minutes: 15 },
					geometry: { type: 'LineString', coordinates: [] }
				},
				{
					type: 'Feature',
					properties: { route_index: 2, distance_km: 6.75, time_minutes: 18 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadRouteLayer(geojson);

		const calls = (L as any)._geoJsonCalls;
		const alternateFeature = calls[1].layer._features[0].layer;
		expect(alternateFeature._popupContent).toContain('<strong>Route 2</strong>');
		expect(alternateFeature._popupContent).toContain('Distance: 6.75 km');
		expect(alternateFeature._popupContent).toContain('Time: 18 min');
	});

	it('uses backend time_minutes for controls when present', () => {
		const L = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { route_index: 0, distance_km: 5, time_minutes: 15 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadRouteLayer(geojson);

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
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { route_index: 0, distance_km: 10 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadRouteLayer(geojson);

		const tc = (lm as any).timeControl as any;
		expect(tc).toBeDefined();
		expect(tc._el).toBeInstanceOf(HTMLElement);
		// default avg speed used in fallback is 20 km/h -> 10km = 30min
		expect(tc._el.innerHTML).toContain('Estimated time: 30 min');
	});

	it('handles string values for distance and time', () => {
		const L = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { route_index: 1, distance_km: '5.5', time_minutes: '16' },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadRouteLayer(geojson);

		const calls = (L as any)._geoJsonCalls;
		const featureLayer = calls[0].layer._features[0].layer;
		expect(featureLayer._popupContent).toContain('Distance: 5.50 km');
		expect(featureLayer._popupContent).toContain('Time: 16 min');
	});

	it('removes previous route layer when loading new one', () => {
		const L = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson1: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { route_index: 1, distance_km: 5, time_minutes: 15 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		const geojson2: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { route_index: 1, distance_km: 8, time_minutes: 20 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadRouteLayer(geojson1);
		const firstLayerCount = (lm.map as any)._layers.length;

		lm.loadRouteLayer(geojson2);
		const secondLayerCount = (lm.map as any)._layers.length;

		// Should have same number of layers (old removed, new added)
		expect(secondLayerCount).toBe(firstLayerCount);
	});
});

describe('LeafletMap.removeRouteLayer', () => {
	beforeEach(() => {
		let el = document.getElementById('map');
		if (!el) {
			el = document.createElement('div');
			el.id = 'map';
			document.body.appendChild(el);
		}
	});

	it('removes route layer and controls', () => {
		const L = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { route_index: 1, distance_km: 5, time_minutes: 15 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadRouteLayer(geojson);
		expect((lm as any).routeLayer).toBeDefined();
		expect((lm as any).distanceControl).toBeDefined();
		expect((lm as any).timeControl).toBeDefined();

		lm.removeRouteLayer();

		expect((lm as any).routeLayer).toBeNull();
		expect((lm as any).distanceControl).toBeNull();
		expect((lm as any).timeControl).toBeNull();
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

	it('renders highlight with cyan style', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const feature: any = {
			type: 'Feature',
			properties: { id: 123 },
			geometry: { type: 'LineString', coordinates: [] }
		};

		lm.loadSelectedWayLayer(feature);

		const calls = (L as any)._geoJsonCalls;
		expect(calls.length).toBe(1);
		expect(calls[0].options.style.color).toBe('#00e5ff');
		expect(calls[0].options.style.weight).toBe(6);
		expect(calls[0].options.style.opacity).toBe(0.9);
	});

	it('replaces existing selected way layer', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const featureA: any = {
			type: 'Feature',
			properties: { id: 123 },
			geometry: { type: 'LineString', coordinates: [] }
		};
		const featureB: any = {
			type: 'Feature',
			properties: { id: 456 },
			geometry: { type: 'LineString', coordinates: [] }
		};

		lm.loadSelectedWayLayer(featureA);
		const firstLayer = (lm as any).selectedWayLayer;

		lm.loadSelectedWayLayer(featureB);
		const secondLayer = (lm as any).selectedWayLayer;

		expect(firstLayer).not.toBe(secondLayer);
		expect((lm.map as any)._layers).not.toContain(firstLayer);
	});
});

describe('LeafletMap.removeSelectedWayLayer', () => {
	beforeEach(() => {
		let el = document.getElementById('map');
		if (!el) {
			el = document.createElement('div');
			el.id = 'map';
			document.body.appendChild(el);
		}
	});

	it('removes selected way layer from map', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const feature: any = {
			type: 'Feature',
			properties: { id: 123 },
			geometry: { type: 'LineString', coordinates: [] }
		};

		lm.loadSelectedWayLayer(feature);
		expect((lm as any).selectedWayLayer).toBeDefined();

		lm.removeSelectedWayLayer();
		expect((lm as any).selectedWayLayer).toBeNull();
	});
});

describe('LeafletMap.loadAdjacentWaysLayer', () => {
	beforeEach(() => {
		let el = document.getElementById('map');
		if (!el) {
			el = document.createElement('div');
			el.id = 'map';
			document.body.appendChild(el);
		}
	});

	it('renders adjacent ways with orange style', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { id: 123 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadAdjacentWaysLayer(geojson);

		const calls = (L as any)._geoJsonCalls;
		expect(calls.length).toBe(1);
		expect(calls[0].options.style.color).toBe('#ff8c00');
		expect(calls[0].options.style.weight).toBe(4);
		expect(calls[0].options.style.opacity).toBe(0.7);
	});

	it('attaches click handlers to adjacent ways', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);
		const onWayClick = vi.fn();

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { id: 123 },
					geometry: { type: 'LineString', coordinates: [] }
				},
				{
					type: 'Feature',
					properties: { id: 456 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadAdjacentWaysLayer(geojson, onWayClick);

		const calls = (L as any)._geoJsonCalls;
		const layer = calls[0].layer;

		// Simulate clicking the first feature
		const clickEvent = { originalEvent: {} };
		layer._features[0].layer._events.click[0](clickEvent);

		expect(onWayClick).toHaveBeenCalledWith(123);

		// Simulate clicking the second feature
		layer._features[1].layer._events.click[0](clickEvent);
		expect(onWayClick).toHaveBeenCalledWith(456);
		expect(onWayClick).toHaveBeenCalledTimes(2);
	});

	it('replaces existing adjacent ways layer', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson1: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { id: 123 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		const geojson2: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { id: 456 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadAdjacentWaysLayer(geojson1);
		const firstLayer = (lm as any).adjacentWaysLayer;

		lm.loadAdjacentWaysLayer(geojson2);
		const secondLayer = (lm as any).adjacentWaysLayer;

		expect(firstLayer).not.toBe(secondLayer);
		expect((lm.map as any)._layers).not.toContain(firstLayer);
	});
});

describe('LeafletMap.removeAdjacentWaysLayer', () => {
	beforeEach(() => {
		let el = document.getElementById('map');
		if (!el) {
			el = document.createElement('div');
			el.id = 'map';
			document.body.appendChild(el);
		}
	});

	it('removes adjacent ways layer from map', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { id: 123 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadAdjacentWaysLayer(geojson);
		expect((lm as any).adjacentWaysLayer).toBeDefined();

		lm.removeAdjacentWaysLayer();
		expect((lm as any).adjacentWaysLayer).toBeNull();
	});
});

describe('LeafletMap.loadAdditionalSelectedWaysLayer', () => {
	beforeEach(() => {
		let el = document.getElementById('map');
		if (!el) {
			el = document.createElement('div');
			el.id = 'map';
			document.body.appendChild(el);
		}
	});

	it('renders additional selected ways with green style', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { id: 789 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadAdditionalSelectedWaysLayer(geojson);

		const calls = (L as any)._geoJsonCalls;
		expect(calls.length).toBe(1);
		expect(calls[0].options.style.color).toBe('#22c55e');
		expect(calls[0].options.style.weight).toBe(4);
		expect(calls[0].options.style.opacity).toBe(0.8);
	});

	it('replaces existing additional selected ways layer', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson1: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { id: 123 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		const geojson2: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { id: 456 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadAdditionalSelectedWaysLayer(geojson1);
		const firstLayer = (lm as any).additionalSelectedWaysLayer;

		lm.loadAdditionalSelectedWaysLayer(geojson2);
		const secondLayer = (lm as any).additionalSelectedWaysLayer;

		expect(firstLayer).not.toBe(secondLayer);
		expect((lm.map as any)._layers).not.toContain(firstLayer);
	});
});

describe('LeafletMap.removeAdditionalSelectedWaysLayer', () => {
	beforeEach(() => {
		let el = document.getElementById('map');
		if (!el) {
			el = document.createElement('div');
			el.id = 'map';
			document.body.appendChild(el);
		}
	});

	it('removes additional selected ways layer from map', () => {
		const L: any = makeFakeL();
		const lm = new LeafletMap(L);

		const geojson: any = {
			type: 'FeatureCollection',
			features: [
				{
					type: 'Feature',
					properties: { id: 789 },
					geometry: { type: 'LineString', coordinates: [] }
				}
			]
		};

		lm.loadAdditionalSelectedWaysLayer(geojson);
		expect((lm as any).additionalSelectedWaysLayer).toBeDefined();

		lm.removeAdditionalSelectedWaysLayer();
		expect((lm as any).additionalSelectedWaysLayer).toBeNull();
	});
});
