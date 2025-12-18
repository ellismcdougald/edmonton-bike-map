/**
 * LeafletMap.ts
 *
 * Purpose:
 * Provides a wrapper around Leaflet for Edmonton Bike Map. Encapsulates map setup,
 * marker management, route and info layers, and click handling.
 *
 * State:
 * - map (Leaflet Map instance)
 * - startMarker (Marker | null): draggable start location marker
 * - endMarker (Marker | null): draggable end location marker
 * - routeLayer (GeoJSON | null): current route layer
 * - infoLayer (GeoJSON | null): layer showing all ways/streets
 *
 * Constants:
 * - EDMONTON_BOUNDS: map bounding box
 * - MAP_START: default map center
 * - MAX_BOUNDS_VISCOSITY, MIN_ZOOM, INITIAL_ZOOM: map settings
 *
 * Behavior:
 * - Initializes Leaflet map with custom panes for route and info layers
 * - Adds tile layers with addTileLayer()
 * - Loads info layer (streets) with loadInfoLayer(), with optional click callback
 * - Shows/hides info layer
 * - Adds/removes draggable start/end markers and retrieves their coordinates
 * - Loads/removes route layer with fit-to-bounds
 * - Provides onMapClick() for handling user clicks on the map
 *
 * Notes:
 * - Depends on 'leaflet' types and Leaflet itself (passed in constructor)
 * - Info layer popups show street/way names; route layer popups show route segment names
 * - Designed for Edmonton-specific map bounds and initial view
 */

import type { Map as LeafletMapType, Marker, GeoJSON, TileLayer } from 'leaflet';

import type {
	GeoJSONOptions,
	PathOptions,
	LatLngExpression,
	LeafletMouseEvent,
	LatLngBoundsLiteral
} from 'leaflet';

import type { WayFeature } from '../types';

export class LeafletMap {
	private L: typeof import('leaflet');

	static readonly EDMONTON_BOUNDS: LatLngBoundsLiteral = [
		[53.3951, -113.7167],
		[53.7169, -113.2437]
	];
	static readonly MAP_START: LatLngExpression = [53.5461, -113.4938];
	static readonly MAX_BOUNDS_VISCOSITY = 1.0;
	static readonly MIN_ZOOM = 10;
	static readonly INITIAL_ZOOM = 12;

	public map: LeafletMapType;
	private startMarker: Marker | null = null;
	private endMarker: Marker | null = null;
	private routeLayer: GeoJSON | null = null;
	private otherRoutesLayer: GeoJSON | null = null;
	private infoLayer: GeoJSON | null = null;
	private selectedWayLayer: GeoJSON | null = null;
	private adjacentWaysLayer: GeoJSON | null = null;
	private additionalSelectedWaysLayer: GeoJSON | null = null;
	private distanceControl: L.Control | null = null;
	private timeControl: L.Control | null = null;

	constructor(L: typeof import('leaflet')) {
		this.L = L;
		this.map = L.map('map', {
			maxBounds: LeafletMap.EDMONTON_BOUNDS,
			maxBoundsViscosity: LeafletMap.MAX_BOUNDS_VISCOSITY,
			minZoom: LeafletMap.MIN_ZOOM
		}).setView(LeafletMap.MAP_START, LeafletMap.INITIAL_ZOOM);

		this.map.createPane('routePane');
		this.map.getPane('routePane')!.style.zIndex = '650';
		this.map.createPane('infoPane');
		this.map.getPane('infoPane')!.style.zIndex = '600';

		// Set up cursor behavior: default normally, grabbing when dragging
		const mapContainer = this.map.getContainer();
		mapContainer.style.cursor = 'default';

		this.map.on('mousedown', () => {
			mapContainer.style.cursor = 'grabbing';
		});

		this.map.on('mouseup', () => {
			mapContainer.style.cursor = 'default';
		});
	}

	addTileLayer(tileUrl: string, maxZoom: number, minZoom: number, attribution: string): TileLayer {
		return this.L.tileLayer(tileUrl, {
			maxZoom,
			minZoom,
			attribution
		}).addTo(this.map);
	}

	loadInfoLayer(
		geojson: GeoJSON.GeoJsonObject,
		style?: PathOptions,
		onClick?: (way: WayFeature) => void
	): void {
		if (this.infoLayer) {
			this.map.removeLayer(this.infoLayer);
			this.infoLayer = null;
		}
		const options: GeoJSONOptions = {
			pane: 'infoPane',
			style: style || { color: 'red', weight: 1, opacity: 0 },
			interactive: true,
			onEachFeature: (feature, layer) => {
				const name = feature.properties?.name || 'Unnamed street';
				layer.bindPopup(`<strong>${name}</strong>`);

				if (onClick) {
					layer.on('click', () => {
						const wayFeature: WayFeature = {
							id: feature.properties?.id,
							tags: feature.properties || {}
						};
						onClick(wayFeature);
					});
				}
			}
		};
		this.infoLayer = this.L.geoJSON(geojson, options).addTo(this.map);
	}

	showInfoLayer(): void {
		if (this.infoLayer && !this.map.hasLayer(this.infoLayer)) {
			this.infoLayer.addTo(this.map);
		}
	}

	hideInfoLayer(): void {
		if (this.infoLayer && this.map.hasLayer(this.infoLayer)) {
			this.map.removeLayer(this.infoLayer);
		}
	}

	addStartMarker(latlng: LatLngExpression, popupText = 'Start'): void {
		if (this.startMarker) this.map.removeLayer(this.startMarker);
		this.startMarker = this.L.marker(latlng, { draggable: true })
			.addTo(this.map)
			.bindPopup(popupText)
			.openPopup();
	}

	addEndMarker(latlng: LatLngExpression, popupText = 'End'): void {
		if (this.endMarker) this.map.removeLayer(this.endMarker);
		this.endMarker = this.L.marker(latlng, { draggable: true })
			.addTo(this.map)
			.bindPopup(popupText)
			.openPopup();
	}

	removeStartMarker(): void {
		if (this.startMarker) {
			this.map.removeLayer(this.startMarker);
			this.startMarker = null;
		}
	}

	removeEndMarker(): void {
		if (this.endMarker) {
			this.map.removeLayer(this.endMarker);
			this.endMarker = null;
		}
	}

	getStartLatLng(): [number, number] | null {
		if (!this.startMarker) return null;
		const { lat, lng } = this.startMarker.getLatLng();
		return [lat, lng];
	}

	getEndLatLng(): [number, number] | null {
		if (!this.endMarker) return null;
		const { lat, lng } = this.endMarker.getLatLng();
		return [lat, lng];
	}

	loadRouteLayer(geojson: GeoJSON.GeoJsonObject): void {
		if (this.routeLayer) {
			this.map.removeLayer(this.routeLayer);
			this.routeLayer = null;
		}
		this.endMarker?.closePopup();

		const featureCollection = geojson as GeoJSON.FeatureCollection;

		// Split geojson into shortest feature and other features
		const shortestFeature = featureCollection.features?.find(
			(f) => f.properties?.['route_index'] === 1 || f.properties?.['route_index'] === '1'
		);
		const otherFeatures = featureCollection.features?.filter(
			(f) => f.properties?.['route_index'] !== 1 && f.properties?.['route_index'] !== '1'
		);

		// Helper to create popup content
		const createPopupContent = (feature: GeoJSON.Feature): string => {
			const routeIndex = feature.properties?.['route_index'];
			const distance = feature.properties?.['distance_km'];
			const time = feature.properties?.['time_minutes'];

			let popupContent = '<strong>Route';
			if (routeIndex !== undefined) {
				if (routeIndex === 1) {
					popupContent = '<strong>Best Route';
				} else {
					popupContent += ` ${routeIndex}`;
				}
			}
			popupContent += '</strong><br/>';

			if (distance !== undefined) {
				const distNum = typeof distance === 'number' ? distance : parseFloat(distance);
				popupContent += `Distance: ${distNum.toFixed(2)} km<br/>`;
			}

			if (time !== undefined) {
				const timeNum = typeof time === 'number' ? Math.round(time) : Math.round(parseFloat(time));
				popupContent += `Time: ${timeNum} min`;
			}

			return popupContent;
		};

		// Helper to create GeoJSON layer with specified style
		const createRouteLayer = (
			data: GeoJSON.FeatureCollection | null,
			style: PathOptions
		): GeoJSON | null => {
			if (!data) return null;
			return this.L.geoJSON(data, {
				pane: 'routePane',
				style,
				interactive: true,
				onEachFeature: (feature, layer) => {
					layer.bindPopup(createPopupContent(feature));
				}
			});
		};

		// Create layers for shortest and other routes
		const shortestLayer = createRouteLayer(
			shortestFeature
				? ({ type: 'FeatureCollection', features: [shortestFeature] } as GeoJSON.FeatureCollection)
				: null,
			{ color: 'blue', weight: 5 }
		);
		const othersLayer = createRouteLayer(
			otherFeatures && otherFeatures.length > 0
				? ({ type: 'FeatureCollection', features: otherFeatures } as GeoJSON.FeatureCollection)
				: null,
			{ color: 'green', weight: 5 }
		);

		// Add layers to map (add others first so shortest renders on top)
		if (othersLayer) {
			this.otherRoutesLayer = othersLayer.addTo(this.map);
		}
		if (shortestLayer) {
			this.routeLayer = shortestLayer.addTo(this.map);
		}

		// Fit bounds
		if (this.routeLayer) {
			const bounds = this.routeLayer.getBounds();
			if (bounds.isValid()) {
				this.map.fitBounds(bounds);
			}
		}

		// Find shortest route for distance/time display
		const shortestRouteFeature =
			featureCollection.features?.find(
				(f) => f.properties?.['route_index'] === 1 || f.properties?.['route_index'] === '1'
			) ||
			featureCollection.features?.[0] ||
			((geojson as GeoJSON.Feature).properties ? (geojson as GeoJSON.Feature) : null);

		// Helper to parse numeric value
		const parseNumericValue = (value: unknown): number => {
			if (typeof value === 'number') return value;
			if (typeof value === 'string') {
				const parsed = parseFloat(value);
				return Number.isFinite(parsed) ? parsed : 0;
			}
			return 0;
		};

		const distanceNum = parseNumericValue(shortestRouteFeature?.properties?.['distance_km']);
		const rawTime = shortestRouteFeature?.properties?.['time_minutes'];
		let timeMin = parseNumericValue(rawTime);

		// Fallback estimate if time not provided
		if (timeMin === 0 && rawTime === undefined) {
			const avgSpeedKmh = 20;
			timeMin = Math.round((distanceNum / avgSpeedKmh) * 60);
		} else if (rawTime !== undefined) {
			timeMin = Math.round(timeMin);
		}

		// Add distance control
		this.distanceControl = new (this.L.Control.extend({
			onAdd: () => {
				const div = this.L.DomUtil.create('div', 'distance-control') as HTMLDivElement;
				div.innerHTML = `Distance: ${distanceNum.toFixed(2)} km`;
				return div;
			}
		}))({ position: 'topright' });
		this.distanceControl.addTo(this.map);

		// Add time control
		this.timeControl = new (this.L.Control.extend({
			onAdd: () => {
				const div = this.L.DomUtil.create('div', 'time-control') as HTMLDivElement;
				div.innerHTML = `Estimated time: ${timeMin} min`;
				return div;
			}
		}))({ position: 'topright' });
		this.timeControl.addTo(this.map);
	}

	removeRouteLayer(): void {
		if (this.routeLayer) {
			this.map.removeLayer(this.routeLayer);
			this.routeLayer = null;
			this.distanceControl?.remove();
			this.distanceControl = null;
			this.timeControl?.remove();
			this.timeControl = null;
		}

		if (this.otherRoutesLayer) {
			this.map.removeLayer(this.otherRoutesLayer);
			this.otherRoutesLayer = null;
		}
	}

	loadSelectedWayLayer(geojson: GeoJSON.GeoJsonObject): void {
		if (this.selectedWayLayer) {
			this.map.removeLayer(this.selectedWayLayer);
			this.selectedWayLayer = null;
		}

		this.selectedWayLayer = this.L.geoJSON(geojson, {
			pane: 'routePane',
			style: { color: '#00e5ff', weight: 6, opacity: 0.9 }
		}).addTo(this.map);
	}

	removeSelectedWayLayer(): void {
		if (this.selectedWayLayer) {
			this.map.removeLayer(this.selectedWayLayer);
			this.selectedWayLayer = null;
		}
	}

	loadAdjacentWaysLayer(
		geojson: GeoJSON.GeoJsonObject,
		onWayClick?: (wayId: number) => void
	): void {
		if (this.adjacentWaysLayer) {
			this.map.removeLayer(this.adjacentWaysLayer);
			this.adjacentWaysLayer = null;
		}

		this.adjacentWaysLayer = this.L.geoJSON(geojson, {
			pane: 'routePane',
			style: { color: '#ff8c00', weight: 4, opacity: 0.7, className: 'adjacent-way' },
			onEachFeature: (feature, layer) => {
				if (onWayClick && feature.properties?.id) {
					layer.on('click', (e) => {
						this.L.DomEvent.stopPropagation(e);
						onWayClick(Number(feature.properties.id));
					});
				}
			}
		}).addTo(this.map);
	}

	removeAdjacentWaysLayer(): void {
		if (this.adjacentWaysLayer) {
			this.map.removeLayer(this.adjacentWaysLayer);
			this.adjacentWaysLayer = null;
		}
	}

	loadAdditionalSelectedWaysLayer(geojson: GeoJSON.GeoJsonObject): void {
		if (this.additionalSelectedWaysLayer) {
			this.map.removeLayer(this.additionalSelectedWaysLayer);
			this.additionalSelectedWaysLayer = null;
		}

		this.additionalSelectedWaysLayer = this.L.geoJSON(geojson, {
			pane: 'routePane',
			style: { color: '#22c55e', weight: 4, opacity: 0.8, className: 'additional-selected-way' }
		}).addTo(this.map);
	}

	removeAdditionalSelectedWaysLayer(): void {
		if (this.additionalSelectedWaysLayer) {
			this.map.removeLayer(this.additionalSelectedWaysLayer);
			this.additionalSelectedWaysLayer = null;
		}
	}

	onMapClick(handler: (latlng: [number, number]) => void): void {
		this.map.on('click', (e: LeafletMouseEvent) => {
			handler([e.latlng.lat, e.latlng.lng]);
		});
	}

	reset(): void {
		this.removeStartMarker();
		this.removeEndMarker();
		this.removeRouteLayer();
		this.removeSelectedWayLayer();
		this.removeAdjacentWaysLayer();
		this.showInfoLayer();
	}
}
