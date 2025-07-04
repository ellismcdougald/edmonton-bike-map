// Map centred on Edmonton
const map = L.map("map").setView([53.5461, -113.4938], 12);

// Add OpenStreetMap tiles
L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 19,
  attribution: "© OpenStreetMap contributors",
}).addTo(map);

// Plot route on map
document.body.addEventListener("htmx:afterOnLoad", (event) => {
  if (event.detail.target.id === "route-data") {
    console.log("1");
    try {
      const geojson = JSON.parse(event.detail.xhr.responseText);

      L.geoJSON(geojson, {
        style: {
          color: "blue",
          weight: 5,
        },
      }).addTo(map);
    } catch (err) {
      console.error("Error parsing or displaying GeoJSON:", err);
    }
  }
});
