// Map centred on Edmonton and bounded on Edmonton
const edmontonBounds = [
  [53.3951, -113.7167],
  [53.7169, -113.2437],
];

const map = L.map("map", {
  maxBounds: edmontonBounds,
  maxBoundsViscosity: 1.0,
  minZoom: 10,
}).setView([53.5461, -113.4938], 12);

// Add OpenStreetMap tiles
L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 19,
  minZoom: 10,
  attribution: "© OpenStreetMap contributors",
}).addTo(map);

// Get data
fetch("edmonton_bike_data_geo.json")
  .then((response) => {
    if (!response.ok) throw new Error("No response from file.");
    return response.json();
  })
  .then((geojson) => {
    L.geoJSON(geojson, {
      style: {
        color: "red",
        weight: 1,
        opacity: 0,
      },
      interactive: true,
      onEachFeature: (feature, layer) => {
        let name = feature.properties.name || "Unnamed street";
        layer.bindPopup(`<strong>${name}</strong>`);
      },
    }).addTo(map);
  })
  .catch((error) => {
    console.error("Could not load GeoJSON file.");
  });

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
        interactive: true,
        onEachFeature: (feature, layer) => {
          let name = feature.properties.name || "Unnamed street";
          layer.bindPopup(`<strong>${name}</strong>`);
        },
      }).addTo(map);
    } catch (err) {
      console.error("Error parsing or displaying GeoJSON:", err);
    }
  }
});

// Track start and end locations
let startMarker = null;
let endMarker = null;
map.on("click", function (e) {
  const { lat, lng } = e.latlng;

  if (!startMarker) {
    startMarker = L.marker([lat, lng], { draggable: true })
      .addTo(map)
      .bindPopup("Start")
      .openPopup();
  } else if (!endMarker) {
    endMarker = L.marker([lat, lng], { draggable: true })
      .addTo(map)
      .bindPopup("End")
      .openPopup();
  }
});
