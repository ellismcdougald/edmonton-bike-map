export class LeafletMap {
    constructor() {
        this.startMarker = null; // L.Marker
        this.endMarker = null; // L.Marker
        this.routeLayer = null; // L.GeoJSON
        this.infoLayer = null; // L.GeoJSON
        this.map = L.map("map", {
            maxBounds: LeafletMap.EDMONTON_BOUNDS,
            maxBoundsViscosity: LeafletMap.MAX_BOUNDS_VISCOSITY,
            minZoom: LeafletMap.MIN_ZOOM,
        }).setView(LeafletMap.MAP_START, LeafletMap.INITIAL_ZOOM);
        this.map.createPane("routePane");
        this.map.getPane("routePane").style.zIndex = "650"; // display above overlayPane (zindex 400)
    }
    // Tile layer:
    addTileLayer(tileUrl, maxZoom, minZoom, attribution) {
        L.tileLayer(tileUrl, {
            maxZoom,
            minZoom,
            attribution,
        }).addTo(this.map);
    }
    // Info layer:
    loadInfoLayer(geojson, style) {
        /*
        if (this.infoLayer) {
          this.map.removeLayer(this.infoLayer);
          this.infoLayer = null;
        }
        this.infoLayer = L.geoJSON(geojson, {
          className: "info-layer",
          style: style || { color: "red", weight: 1, opacity: 0 },
          interactive: true,
          onEachFeature: (feature: any, layer: any) => {
            const name = feature.properties?.name || "Unnamed street";
            layer.bindPopup(`<strong>${name}</strong>`);
          },
        }).addTo(this.map);
        */
    }
    showInfoLayer() {
        if (this.infoLayer && !this.map.hasLayer(this.infoLayer)) {
            this.infoLayer.addTo(this.map);
        }
    }
    hideInfoLayer() {
        if (this.infoLayer && this.map.hasLayer(this.infoLayer)) {
            this.map.removeLayer(this.infoLayer);
        }
    }
    // Markers:
    addStartMarker(latlng, popupText = "Start") {
        if (this.startMarker)
            this.map.removeLayer(this.startMarker);
        this.startMarker = L.marker(latlng, { draggable: true })
            .addTo(this.map)
            .bindPopup(popupText)
            .openPopup();
    }
    addEndMarker(latlng, popupText = "End") {
        if (this.endMarker)
            this.map.removeLayer(this.endMarker);
        this.endMarker = L.marker(latlng, { draggable: true })
            .addTo(this.map)
            .bindPopup(popupText)
            .openPopup();
    }
    removeStartMarker() {
        if (this.startMarker) {
            this.map.removeLayer(this.startMarker);
            this.startMarker = null;
        }
    }
    removeEndMarker() {
        if (this.endMarker) {
            this.map.removeLayer(this.endMarker);
            this.endMarker = null;
        }
    }
    getStartLatLng() {
        if (!this.startMarker)
            return null;
        const { lat, lng } = this.startMarker.getLatLng();
        return [lat, lng];
    }
    getEndLatLng() {
        if (!this.endMarker)
            return null;
        const { lat, lng } = this.endMarker.getLatLng();
        return [lat, lng];
    }
    // Route layer:
    loadRouteLayer(geojson) {
        if (this.routeLayer) {
            this.map.removeLayer(this.routeLayer);
            this.routeLayer = null;
        }
        this.routeLayer = L.geoJSON(geojson, {
            pane: "routePane",
            style: { color: "blue", weight: 5 },
            interactive: true,
            onEachFeature: (feature, layer) => {
                var _a;
                const name = ((_a = feature.properties) === null || _a === void 0 ? void 0 : _a.name) || "Unnamed route segment";
                layer.bindPopup(`<strong>${name}</strong>`);
                const element = layer.getElement();
                if (element) {
                    element.classList.add("route-path");
                }
            },
        }).addTo(this.map);
        const bounds = this.routeLayer.getBounds();
        if (bounds.isValid()) {
            this.map.fitBounds(bounds);
        }
    }
    removeRouteLayer() {
        if (this.routeLayer) {
            this.map.removeLayer(this.routeLayer);
            this.routeLayer = null;
        }
    }
    onMapClick(handler) {
        this.map.on("click", (e) => {
            handler([e.latlng.lat, e.latlng.lng]);
        });
    }
}
LeafletMap.EDMONTON_BOUNDS = [
    [53.3951, -113.7167],
    [53.7169, -113.2437],
];
LeafletMap.MAP_START = [53.5461, -113.4938];
LeafletMap.MAX_BOUNDS_VISCOSITY = 1.0;
LeafletMap.MIN_ZOOM = 10;
LeafletMap.INITIAL_ZOOM = 12;
