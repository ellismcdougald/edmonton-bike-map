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
	private infoLayer: GeoJSON | null = null;
	private selectedWayLayer: GeoJSON | null = null;
	private adjacentWaysLayer: GeoJSON | null = null;
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

		this.routeLayer = this.L.geoJSON(geojson, {
			pane: 'routePane',
			style: { color: 'blue', weight: 5 },
			interactive: true,
			onEachFeature: (feature, layer) => {
				const name = feature.properties?.name || 'Unnamed route segment';
				layer.bindPopup(`<strong>${name}</strong>`);
			}
		}).addTo(this.map);

		const bounds = this.routeLayer.getBounds();
		if (bounds.isValid()) {
			this.map.fitBounds(bounds);
		}

		const feature = geojson as GeoJSON.Feature;
		// read distance and time from backend properties if available
		const rawDistance = feature.properties?.['distance_km'];
		let distanceNum = 0;
		if (typeof rawDistance === 'number') {
			distanceNum = rawDistance;
		} else if (typeof rawDistance === 'string') {
			const parsed = parseFloat(rawDistance as string);
			distanceNum = Number.isFinite(parsed) ? parsed : 0;
		}

		// Add distance control (safe formatting)
		this.distanceControl = new (this.L.Control.extend({
			onAdd: () => {
				const div = this.L.DomUtil.create('div', 'distance-control') as HTMLDivElement;
				div.innerHTML = `Distance: ${distanceNum.toFixed(2)} km`;
				return div;
			}
		}))({ position: 'topright' });
		this.distanceControl.addTo(this.map);

		// Determine time in minutes: prefer backend-provided `time_minutes`, fall back to estimate
		const rawTime = feature.properties?.['time_minutes'];
		let timeMin = 0;
		if (typeof rawTime === 'number') {
			timeMin = Math.round(rawTime as number);
		} else if (typeof rawTime === 'string') {
			const parsed = parseFloat(rawTime as string);
			timeMin = Number.isFinite(parsed) ? Math.round(parsed) : 0;
		} else {
			// fallback estimate using a sensible default speed
			const avgSpeedKmh = 20; // commuter average speed fallback
			const timeH = distanceNum / avgSpeedKmh;
			timeMin = Math.round(timeH * 60);
		}

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
			style: { color: '#ff8c00', weight: 4, opacity: 0.7 },
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
