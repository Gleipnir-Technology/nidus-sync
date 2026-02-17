// A test of maplibre-gl in a custom element
class MapRouting extends HTMLElement {
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({ mode: "open" });

		// Initial render
		this.render();

		this._map = null;

		// markers shown on the map
		this._markers = [];
	}

	// Lifecycle: when element is added to the DOM
	connectedCallback() {
		// Initialize the map when the element is added to the DOM
		setTimeout(() => this._initializeMap(), 0);
	}

	disconnectedCallback() {
		if (this._map) {
			this._map.remove();
		}
	}

	_initializeMap() {
		const centroid = JSON.parse(this.getAttribute("centroid"));
		const organization_id = this.getAttribute("organization-id");
		const tegola = this.getAttribute("tegola");
		const xmin = parseFloat(this.getAttribute("xmin"));
		const ymin = parseFloat(this.getAttribute("ymin"));
		const xmax = parseFloat(this.getAttribute("xmax"));
		const ymax = parseFloat(this.getAttribute("ymax"));
		const bounds = [
			[xmin, ymin],
			[xmax, ymax],
		];

		const mapElement = this.shadowRoot.querySelector("#map");

		/*
		this._map = new maplibregl.Map({
			center: centroid.coordinates,
			container: mapElement,
			style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json", // Style URL; see our documentation for more options
		}).fitBounds(bounds, {
			padding: { top: 10, bottom: 10, left: 10, right: 10 },
		});
		this._map.on("load", () => {
			this.dispatchEvent(new CustomEvent("load"), {
				bubbles: true,
				composed: true, // Allows event to cross shadow DOM boundary
				detail: {
					map: this,
				},
			});
		});
		*/
		this._map = new maplibregl.Map({
			container: mapElement,
			style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json", // Style URL; see our documentation for more options
			center: [-122.4194, 37.7749], // San Francisco coordinates
			zoom: 12,
		});
		// Sample GeoJSON for a driving route in San Francisco
		const routeData = {
			type: "Feature",
			properties: {},
			geometry: {
				type: "LineString",
				coordinates: [
					[-122.4194, 37.7749], // Start: San Francisco downtown
					[-122.4211, 37.7837],
					[-122.4156, 37.7908],
					[-122.4089, 37.7973],
					[-122.3998, 37.8015],
					[-122.3919, 37.8057],
					[-122.3873, 37.8067], // End: Fisherman's Wharf area
				],
			},
		};

		// Add map controls
		this._map.addControl(new maplibregl.NavigationControl());
		// Wait for the map to load
		this._map.on("load", () => {
			// Add the route source
			this._map.addSource("route", {
				type: "geojson",
				data: routeData,
			});

			// Add a layer to display the route
			this._map.addLayer({
				id: "route",
				type: "line",
				source: "route",
				layout: {
					"line-join": "round",
					"line-cap": "round",
				},
				paint: {
					"line-color": "#3887be",
					"line-width": 5,
					"line-opacity": 0.75,
				},
			});

			// Add start point
			this._map.addSource("start-point", {
				type: "geojson",
				data: {
					type: "Feature",
					geometry: {
						type: "Point",
						coordinates: routeData.geometry.coordinates[0],
					},
					properties: {
						description: "Start",
					},
				},
			});

			this._map.addLayer({
				id: "start-point",
				type: "circle",
				source: "start-point",
				paint: {
					"circle-radius": 8,
					"circle-color": "#3bb2d0",
				},
			});

			// Add end point
			this._map.addSource("end-point", {
				type: "geojson",
				data: {
					type: "Feature",
					geometry: {
						type: "Point",
						coordinates:
							routeData.geometry.coordinates[
								routeData.geometry.coordinates.length - 1
							],
					},
					properties: {
						description: "End",
					},
				},
			});

			this._map.addLayer({
				id: "end-point",
				type: "circle",
				source: "end-point",
				paint: {
					"circle-radius": 8,
					"circle-color": "#f30",
				},
			});
		});
	}

	// Initial render of component
	render() {
		this.shadowRoot.innerHTML = `
			<style>
				@import url('//unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css');
				.map-container {
					background-color: #e9ecef;
					border-radius: 10px;
					box-shadow: 0 4px 6px rgba(0,0,0,0.05);
					height: 500px;
					display: flex;
					align-items: center;
					justify-content: center;
					margin-top: 20px;
				}
				#map {
					height: 500px;
					width:100%;
					margin-bottom: 10px;
				}
				#map img {
					max-width: none;
					min-width: 0px;
					height: auto;
				}
			</style>
			
			<div id="map-container" class="map-container">
				<div id="map"></div>
			</div>
		`;
	}

	addLayer(a) {
		return this._map.addLayer(a);
	}
	addSource(a, b) {
		return this._map.addSource(a, b);
	}
	jumpTo(args) {
		return this._map.jumpTo(args);
	}
	on(a, b) {
		return this._map.on(a, b);
	}
	once(a, b) {
		return this._map.once(a, b);
	}
	queryRenderedFeatures(a) {
		return this._map.queryRenderedFeatures(a);
	}

	setMarker(coords) {
		console.log("Setting map marker", coords);
		this._map.jumpTo({
			center: coords,
			zoom: 14,
		});
		this._markers.forEach((marker) => marker.remove());

		const marker = new mapboxgl.Marker({
			color: "#FF0000",
			draggable: true,
		})
			.setLngLat(coords)
			.addTo(map);
		marker.on("dragend", function (e) {
			const markerDraggedEvent = new CustomEvent("markerdragend", {
				detail: {
					marker: marker,
				},
			});
			mapContainer.dispatchEvent(markerDraggedEvent);
		});
		this._markers = [marker];
	}

	SetLayoutProperty(layout, property, value) {
		return this._map.setLayoutProperty(layout, property, value);
	}
}

customElements.define("map-routing", MapRouting);
