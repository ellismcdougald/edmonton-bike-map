import { LeafletMap } from "./Map.js";
export class App {
    constructor() {
        // Logic:
        this.selectStartActive = false;
        this.selectEndActive = false;
        // UI Elements:
        this.selectStartButton = null;
        this.selectEndButton = null;
        this.findRouteButton = null;
        this.resetButton = null;
        this.map = new LeafletMap();
        this.map.addTileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", 19, 10, "© OpenStreetMap contributors");
        this.initUI();
        this.loadBikeData();
        this.map.onMapClick(this.onMapClick.bind(this));
    }
    initUI() {
        var _a, _b, _c, _d;
        this.selectStartButton = document.getElementById("selectStartButton");
        (_a = this.selectStartButton) === null || _a === void 0 ? void 0 : _a.addEventListener("click", () => {
            this.toggleSelectStartActive();
        });
        this.selectEndButton = document.getElementById("selectEndButton");
        (_b = this.selectEndButton) === null || _b === void 0 ? void 0 : _b.addEventListener("click", () => {
            this.toggleSelectEndActive();
        });
        this.findRouteButton = document.getElementById("findRouteButton");
        (_c = this.findRouteButton) === null || _c === void 0 ? void 0 : _c.addEventListener("click", () => {
            this.findRoute();
        });
        this.resetButton = document.getElementById("resetButton");
        (_d = this.resetButton) === null || _d === void 0 ? void 0 : _d.addEventListener("click", () => {
            this.resetState();
        });
    }
    toggleSelectStartActive() {
        if (this.selectStartActive) {
            this.selectStartActive = false;
            this.map.showInfoLayer();
        }
        else {
            this.selectStartActive = true;
            this.selectEndActive = false;
            this.map.hideInfoLayer();
        }
        if (this.selectStartButton)
            this.selectStartButton.classList.toggle("active", this.selectStartActive);
    }
    toggleSelectEndActive() {
        if (this.selectEndActive) {
            this.selectEndActive = false;
            this.map.showInfoLayer();
        }
        else {
            this.selectEndActive = true;
            this.selectStartActive = false;
            this.map.hideInfoLayer();
        }
        if (this.selectEndButton)
            this.selectEndButton.classList.toggle("active", this.selectEndActive);
    }
    findRoute() {
        const startLatLng = this.map.getStartLatLng();
        const endLatLng = this.map.getEndLatLng();
        if (!startLatLng || !endLatLng) {
            alert("Make sure you have selected both a start and an end point!");
            return;
        }
        const params = new URLSearchParams({
            startLatitude: startLatLng[0].toString(),
            startLongitude: startLatLng[1].toString(),
            endLatitude: endLatLng[0].toString(),
            endLongitude: endLatLng[1].toString(),
        });
        fetch(`/api/route?${params.toString()}`)
            .then((res) => {
            if (!res.ok)
                throw new Error("Failed to get route data");
            return res.json();
        })
            .then((geojson) => {
            this.map.removeRouteLayer();
            this.map.loadRouteLayer(geojson);
        })
            .catch((err) => {
            alert("Error fetching or displaying route: " + err.message);
            console.error(err);
        });
    }
    resetState() {
        var _a, _b;
        this.map.removeStartMarker();
        this.map.removeEndMarker();
        this.map.removeRouteLayer();
        this.map.showInfoLayer();
        (_a = this.selectStartButton) === null || _a === void 0 ? void 0 : _a.classList.toggle("active", false);
        (_b = this.selectEndButton) === null || _b === void 0 ? void 0 : _b.classList.toggle("active", false);
        this.selectStartActive = false;
        this.selectEndActive = false;
    }
    loadBikeData() {
        fetch("edmonton_bike_data_geo.json")
            .then((res) => {
            if (!res.ok)
                throw new Error("No response from file.");
            return res.json();
        })
            .then((geojson) => {
            this.map.loadInfoLayer(geojson);
        })
            .catch((err) => {
            console.error("Failed to load bike data:", err);
        });
    }
    onMapClick(latlng) {
        const [lat, lng] = latlng;
        if (this.selectStartActive) {
            this.map.removeStartMarker();
            this.map.addStartMarker([lat, lng]);
            this.toggleSelectStartActive();
            this.map.showInfoLayer();
        }
        else if (this.selectEndActive) {
            this.map.removeEndMarker();
            this.map.addEndMarker([lat, lng]);
            this.toggleSelectEndActive();
            this.map.showInfoLayer();
        }
    }
}
