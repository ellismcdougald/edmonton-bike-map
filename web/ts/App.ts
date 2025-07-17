import { LeafletMap } from "./Map.js";

export class App {
  // Map:
  private map: LeafletMap;
  // Logic:
  private selectStartActive: boolean = false;
  private selectEndActive: boolean = false;
  // UI Elements:
  private selectStartButton: HTMLElement | null = null;
  private selectEndButton: HTMLElement | null = null;
  private findRouteButton: HTMLElement | null = null;
  private resetButton: HTMLElement | null = null;

  constructor() {
    this.map = new LeafletMap();
    this.map.addTileLayer(
      "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
      19,
      10,
      "© OpenStreetMap contributors"
    );
    this.initUI();
    this.loadBikeData();

    this.map.onMapClick(this.onMapClick.bind(this));
  }

  private initUI(): void {
    this.selectStartButton = document.getElementById("selectStartButton");
    this.selectStartButton?.addEventListener("click", () => {
      this.toggleSelectStartActive();
    });

    this.selectEndButton = document.getElementById("selectEndButton");
    this.selectEndButton?.addEventListener("click", () => {
      this.toggleSelectEndActive();
    });

    this.findRouteButton = document.getElementById("findRouteButton");
    this.findRouteButton?.addEventListener("click", () => {
      this.findRoute();
    });

    this.resetButton = document.getElementById("resetButton");
    this.resetButton?.addEventListener("click", () => {
      this.resetState();
    });
  }

  private toggleSelectStartActive(): void {
    if (this.selectStartActive) {
      this.selectStartActive = false;
      this.map.showInfoLayer();
    } else {
      this.selectStartActive = true;
      this.selectEndActive = false;
      this.map.hideInfoLayer();
    }
    if (this.selectStartButton)
      this.selectStartButton.classList.toggle("active", this.selectStartActive);
  }
  private toggleSelectEndActive(): void {
    if (this.selectEndActive) {
      this.selectEndActive = false;
      this.map.showInfoLayer();
    } else {
      this.selectEndActive = true;
      this.selectStartActive = false;
      this.map.hideInfoLayer();
    }
    if (this.selectEndButton)
      this.selectEndButton.classList.toggle("active", this.selectEndActive);
  }

  private findRoute(): void {
    const startLatLng: [number, number] | null = this.map.getStartLatLng();
    const endLatLng: [number, number] | null = this.map.getEndLatLng();

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
        if (!res.ok) throw new Error("Failed to get route data");
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

  private resetState(): void {
    this.map.removeStartMarker();
    this.map.removeEndMarker();
    this.map.removeRouteLayer();
    this.map.showInfoLayer();
    this.selectStartButton?.classList.toggle("active", false);
    this.selectEndButton?.classList.toggle("active", false);
    this.selectStartActive = false;
    this.selectEndActive = false;
  }

  private loadBikeData(): void {
    fetch("edmonton_bike_data_geo.json")
      .then((res) => {
        if (!res.ok) throw new Error("No response from file.");
        return res.json();
      })
      .then((geojson) => {
        this.map.loadInfoLayer(geojson);
      })
      .catch((err) => {
        console.error("Failed to load bike data:", err);
      });
  }

  private onMapClick(latlng: [number, number]): void {
    const [lat, lng] = latlng;

    if (this.selectStartActive) {
      this.map.removeStartMarker();
      this.map.addStartMarker([lat, lng]);
      this.toggleSelectStartActive();
      this.map.showInfoLayer();
    } else if (this.selectEndActive) {
      this.map.removeEndMarker();
      this.map.addEndMarker([lat, lng]);
      this.toggleSelectEndActive();
      this.map.showInfoLayer();
    }
  }
}
