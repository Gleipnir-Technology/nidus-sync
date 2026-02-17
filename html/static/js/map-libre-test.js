// A test of maplibre-gl in a custom element
class MapLibreTest extends HTMLElement {
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
	}

	// Initial render of component
	render() {
		this.shadowRoot.innerHTML = `
			<style>
				@import url('//unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css');
				.mapboxgl-ctrl-bottom-right {
					display: none;
				}
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

customElements.define("map-libre-test", MapLibreTest);
