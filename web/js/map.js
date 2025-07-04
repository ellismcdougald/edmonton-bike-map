// Leaflet map -- centred on Edmonton
const map = L.map("map").setView([53.5461, -113.4938], 12);
// Add OpenStreetMap tiles
L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 19,
  attribution: "© OpenStreetMap contributors",
}).addTo(map);
