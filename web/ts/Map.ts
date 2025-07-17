declare const L: any; // L is a leaflet instance -- not using leaflet types for now

export class LeafletMap {
  static readonly EDMONTON_BOUNDS: [number, number][] = [
    [53.3951, -113.7167],
    [53.7169, -113.2437],
  ];
  static readonly MAP_START: [number, number] = [53.5461, -113.4938];
  static readonly MAX_BOUNDS_VISCOSITY: number = 1.0;
  static readonly MIN_ZOOM: number = 10;
  static readonly INITIAL_ZOOM: number = 12;

  // Map and Details:
  public map: any; // type L.map
  private startMarker: any | null = null; // L.Marker
  private endMarker: any | null = null; // L.Marker
  private routeLayer: any | null = null; // L.GeoJSON
  private infoLayer: any | null = null; // L.GeoJSON

  constructor() {
    this.map = L.map("map", {
      maxBounds: LeafletMap.EDMONTON_BOUNDS,
      maxBoundsViscosity: LeafletMap.MAX_BOUNDS_VISCOSITY,
      minZoom: LeafletMap.MIN_ZOOM,
    }).setView(LeafletMap.MAP_START, LeafletMap.INITIAL_ZOOM);
  }

  // Tile layer:
  addTileLayer(
    tileUrl: string,
    maxZoom: number,
    minZoom: number,
    attribution: string
  ): void {
    L.tileLayer(tileUrl, {
      maxZoom,
      minZoom,
      attribution,
    }).addTo(this.map);
  }

  // Info layer:
  loadInfoLayer(geojson: any, style?: any): void {
    if (this.infoLayer) {
      this.map.removeLayer(this.infoLayer);
      this.infoLayer = null;
    }
    this.infoLayer = L.geoJSON(geojson, {
      style: style || { color: "red", weight: 1, opacity: 0 },
      interactive: true,
      onEachFeature: (feature: any, layer: any) => {
        const name = feature.properties?.name || "Unnamed street";
        layer.bindPopup(`<strong>${name}</strong>`);
      },
    }).addTo(this.map);
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

  // Markers:
  addStartMarker(latlng: [number, number], popupText = "Start"): void {
    if (this.startMarker) this.map.removeLayer(this.startMarker);
    this.startMarker = L.marker(latlng, { draggable: true })
      .addTo(this.map)
      .bindPopup(popupText)
      .openPopup();
  }

  addEndMarker(latlng: [number, number], popupText = "End"): void {
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

  // Route layer:
  loadRouteLayer(geojson: any): void {
    if (this.routeLayer) {
      this.map.removeLayer(this.routeLayer);
      this.routeLayer = null;
    }
    this.routeLayer = L.geoJSON(geojson, {
      style: { color: "blue", weight: 5 },
      interactive: true,
      onEachFeature: (feature: any, layer: any) => {
        const name = feature.properties?.name || "Unnamed route segment";
        layer.bindPopup(`<strong>${name}</strong>`);
      },
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
    this.map.on("click", (e: any) => {
      handler([e.latlng.lat, e.latlng.lng]);
    });
  }
}
