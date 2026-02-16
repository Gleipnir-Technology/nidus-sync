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

	// Lifecycle: watch these attributes for changes
	static get observedAttributes() {
		return ["latitude", "longitude", "zoom"];
	}

	// Lifecycle: respond to attribute changes
	attributeChangedCallback(name, oldValue, newValue) {
		// Only handle if map exists and values actually changed
		if (!this._map || oldValue === newValue) return;

		if (name === "latitude" || name === "longitude") {
			if (this.hasAttribute("latitude") && this.hasAttribute("longitude")) {
				const lat = Number(this.getAttribute("latitude"));
				const lng = Number(this.getAttribute("longitude"));
				this._map.setCenter([lat, lng]);
			}
		}

		if (name === "zoom") {
			this._map.setZoom(Number(newValue));
		}
	}

	_initializeMap() {
		const apiKey = this.getAttribute("api-key");
		const lat = Number(this.getAttribute("latitude") || 36.2);
		const lng = Number(this.getAttribute("longitude") || -119.2);
		const mapElement = this.shadowRoot.querySelector("#map");
		const tegola = this.getAttribute("tegola");
		const zoom = Number(this.getAttribute("zoom") || 15);

		this._map = new maplibregl.Map({
			container: mapElement,
			center: {
				lat: lat,
				lng: lng,
			},
			style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json", // Style URL; see our documentation for more options
			zoom: zoom,
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
