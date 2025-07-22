import L, { Map as LeafletMapType, Marker, GeoJSON, TileLayer } from 'leaflet';

import type {
	GeoJSONOptions,
	PathOptions,
	LatLngExpression,
	LeafletMouseEvent,
	LatLngBoundsLiteral
} from 'leaflet';

export class LeafletMap {
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

	constructor() {
		this.map = L.map('map', {
			maxBounds: LeafletMap.EDMONTON_BOUNDS,
			maxBoundsViscosity: LeafletMap.MAX_BOUNDS_VISCOSITY,
			minZoom: LeafletMap.MIN_ZOOM
		}).setView(LeafletMap.MAP_START, LeafletMap.INITIAL_ZOOM);

		this.map.createPane('routePane');
		this.map.getPane('routePane')!.style.zIndex = '650';
		this.map.createPane('infoPane');
		this.map.getPane('infoPane')!.style.zIndex = '600';
	}

	addTileLayer(tileUrl: string, maxZoom: number, minZoom: number, attribution: string): TileLayer {
		return L.tileLayer(tileUrl, {
			maxZoom,
			minZoom,
			attribution
		}).addTo(this.map);
	}

	loadInfoLayer(geojson: GeoJSON.GeoJsonObject, style?: PathOptions): void {
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
			}
		};
		this.infoLayer = L.geoJSON(geojson, options).addTo(this.map);
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
		this.startMarker = L.marker(latlng, { draggable: true })
			.addTo(this.map)
			.bindPopup(popupText)
			.openPopup();
	}

	addEndMarker(latlng: LatLngExpression, popupText = 'End'): void {
		if (this.endMarker) this.map.removeLayer(this.endMarker);
		this.endMarker = L.marker(latlng, { draggable: true })
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
		this.routeLayer = L.geoJSON(geojson, {
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
	}

	removeRouteLayer(): void {
		if (this.routeLayer) {
			this.map.removeLayer(this.routeLayer);
			this.routeLayer = null;
		}
	}

	onMapClick(handler: (latlng: [number, number]) => void): void {
		this.map.on('click', (e: LeafletMouseEvent) => {
			handler([e.latlng.lat, e.latlng.lng]);
		});
	}
}
