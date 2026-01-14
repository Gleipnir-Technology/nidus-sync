var map = null;
var markers = [];

function mapAddMarker(coords) {
	const mapContainer = document.getElementById("map-container");
	const marker = new mapboxgl.Marker({
		color: "#FF0000",
		draggable: true
	}).setLngLat(coords).addTo(map);
	marker.on('dragend', function(e) {
		const markerDraggedEvent = new CustomEvent("markerdragend", {
			detail: {
				marker: marker
			}
		});
		mapContainer.dispatchEvent(markerDraggedEvent);
	});
	markers.push(marker);
}

function mapLoad(MAPBOX_ACCESS_TOKEN) {
	return new Promise((resolve, reject) => {
		console.log("Setting up the map...");
		mapboxgl.accessToken = MAPBOX_ACCESS_TOKEN;
		map = new mapboxgl.Map({
			container: "map",
			center: {
				lat: 36.2,
				lng: -119.2
			},
			style: 'mapbox://styles/mapbox/streets-v12', // style URL
			zoom: 15,
		});
		map.addControl(new mapboxgl.GeolocateControl({
			positionOptions: {
				enableHighAccuracy: true
			},
			trackUserLocation: true,
			showUserHeading: true
		}));
		map.addControl(new mapboxgl.NavigationControl());
		map.on("load", function() {
			console.log("Map loaded.");
			resolve(map);
		});
	});
}

function mapJumpTo(args) {
	map.jumpTo(args);
}

function mapSetMarker(coords) {
	console.log("Setting map marker", coords);
	map.jumpTo({
		center: coords,
		zoom: 14,
	});
	markers.forEach((marker) => marker.remove());
	mapAddMarker(coords);
}

