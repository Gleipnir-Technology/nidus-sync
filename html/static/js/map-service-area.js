// A test of maplibre-gl in a custom element
class MapServiceArea extends HTMLElement {
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
		const apiKey = this.getAttribute("api-key");
		const centroid = JSON.parse(this.getAttribute("centroid"));
		const csv_file = this.getAttribute("csv-file");
		const organization_id = this.getAttribute("organization-id");
		const lat = Number(this.getAttribute("latitude") || 36.2);
		const lng = Number(this.getAttribute("longitude") || -119.2);
		const mapElement = this.shadowRoot.querySelector("#map");
		const tegola = this.getAttribute("tegola");
		const xmin = parseFloat(this.getAttribute("xmin"));
		const ymin = parseFloat(this.getAttribute("ymin"));
		const xmax = parseFloat(this.getAttribute("xmax"));
		const ymax = parseFloat(this.getAttribute("ymax"));
		const bounds = [
			[xmin, ymin],
			[xmax, ymax],
		];
		console.log("fitting", bounds);
		this._map = new maplibregl.Map({
			container: mapElement,
			center: centroid.coordinates,
			style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json",
		}).fitBounds(bounds, {
			padding: { top: 10, bottom: 10, left: 10, right: 10 },
		});
		this._map.on("load", () => {
			this._map.addSource("tegola-nidus", {
				type: "vector",
				tiles: [`${tegola}maps/nidus/{z}/{x}/{y}?id=${organization_id}`],
			});
			this._map.addLayer({
				id: "service-area",
				source: "tegola-nidus",
				"source-layer": "service-area-bounds",
				type: "fill",
				paint: {
					"fill-opacity": 0.4,
					"fill-color": "#dc3545",
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

	SetLayoutProperty(layout, property, value) {
		return this._map.setLayoutProperty(layout, property, value);
	}
}

customElements.define("map-service-area", MapServiceArea);
