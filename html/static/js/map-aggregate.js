var map = null;
// A map that can be used to locate a single point by setting its location explicitly
// or by allowing the user to move a marker.
class MapAggregate extends HTMLElement {
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({ mode: "open" });

		// Initial render
		this.render();
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
		const organization_id = Number(this.getAttribute("organization-id") || 0);
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
			style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json",
		}).fitBounds(bounds, {
			padding: { top: 10, bottom: 10, left: 10, right: 10 },
		});
		this._map.on("load", () => {
			this._map.addSource("tegola", {
				type: "vector",
				tiles: [
					`${tegola}maps/nidus/{z}/{x}/{y}?id=${organization_id}&organization_id=${organization_id}`,
				],
			});
			this._map.addLayer({
				id: "mosquito_source",
				type: "fill",
				filter: [
					"==",
					["zoom"],
					["+", 2, ["to-number", ["get", "resolution"]]],
				],
				source: "tegola",
				"source-layer": "mosquito_source",
				paint: {
					"fill-opacity": 0.4,
					"fill-color": "#dc3545",
				},
			});
			this._map.addLayer({
				id: "service_request",
				type: "fill",
				filter: [
					"==",
					["zoom"],
					["+", 2, ["to-number", ["get", "resolution"]]],
				],
				source: "tegola",
				"source-layer": "service_request",
				paint: {
					"fill-opacity": 0.4,
					"fill-color": "#ffc107",
				},
			});
			this._map.addLayer({
				id: "trap",
				type: "fill",
				filter: [
					"==",
					["zoom"],
					["+", 2, ["to-number", ["get", "resolution"]]],
				],
				source: "tegola",
				"source-layer": "trap",
				paint: {
					"fill-opacity": 0.4,
					"fill-color": "#0dcaf0",
				},
			});
			this._map.addLayer({
				id: "service-area",
				source: "tegola",
				"source-layer": "service-area-bounds",
				type: "line",
				paint: {
					"line-color": "#f00",
				},
			});
			this._map.on("mouseenter", "mosquito_source", (e) => {
				this._map.getCanvas().style.cursor = "pointer";
			});
			this._map.on("mouseleave", "mosquito_source", (e) => {
				this._map.getCanvas().style.cursor = "";
			});
			const _handleClick = (e) => {
				const feature = e.features[0];
				const coordinates = feature.geometry.coordinates.slice();
				const properties = feature.properties;
				this.dispatchEvent(
					new CustomEvent("cell-click", {
						bubbles: true,
						composed: true, // Allows event to cross shadow DOM boundary
						detail: {
							cell: properties.cell,
						},
					}),
				);
			};
			this._map.on("click", "mosquito_source", _handleClick);
			this._map.on("click", "service_request", _handleClick);
			this._map.on("click", "trap", _handleClick);
		});
	}

	// Initial render of component
	render() {
		this.shadowRoot.innerHTML = `
			<style>
				@import url("//unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css");
				.map-container {
					background-color: #e9ecef;
					border-radius: 10px;
					box-shadow: 0 4px 6px rgba(0,0,0,0.05);
					height: 500px;
					margin-top: 20px;
					position: relative;
				}
				#map {
					position: absolute;
					top: 0;
					bottom: 0;
					left: 0;
					right: 0;
					height: 100%;
					width: 100%
				}
			</style>
			
			<div id="map-container" class="map-container">
				<div id="map"></div>
			</div>
		`;
	}

	jumpTo(args) {
		this._map.jumpTo(args);
	}
}

customElements.define("map-aggregate", MapAggregate);
