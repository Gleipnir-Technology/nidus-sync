var map = null;
// A map that shows multiple single point locations.
// Points have additional detail popups.
class MapMultipoint extends HTMLElement {
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({ mode: "open" });

		// Initial render
		this.render();

		// Keep track of any 'on' calls to add to the map as soon as we create it.
		this._preOns = [];
		this._map = null;
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
			bounds: bounds,
			container: mapElement,
			style: "https://tiles.stadiamaps.com/styles/osm_bright.json",
		});
		this._map.on("load", () => {
			if (organization_id != 0) {
				this._map.addSource("tegola", {
					type: "vector",
					tiles: [
						`${tegola}maps/nidus/{z}/{x}/{y}?id=${organization_id}&organization_id=${organization_id}`,
					],
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
			}
			this.dispatchEvent(new CustomEvent("load"), {
				bubbles: true,
				composed: true, // Allows event to cross shadow DOM boundary
				detail: {
					map: this,
				},
			});
		});
		for (const on of this._preOns) {
			this._map.on(on.a, on.b);
		}
	}

	// Initial render of component
	render() {
		this.shadowRoot.innerHTML = `
			<style>
				@import url("//unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css");
				#map {
					height: 100%;
					width: 100%;
				}
			</style>
			
			<div id="map"></div>
		`;
	}

	addLayer(a) {
		return this._map.addLayer(a);
	}
	addSource(a, b) {
		return this._map.addSource(a, b);
	}
	flyTo(a, b) {
		return this._map.flyTo(a, b);
	}
	jumpTo(args) {
		return this._map.jumpTo(args);
	}
	on(a, b) {
		if (this._map != null) {
			return this._map.on(a, b);
		} else {
			this._preOns.push({ a: a, b: b });
		}
	}
	once(a, b) {
		return this._map.once(a, b);
	}
	panTo(a, b) {
		return this._map.panTo(a, b);
	}
	queryRenderedFeatures(a) {
		return this._map.queryRenderedFeatures(a);
	}

	FitBounds(bounds, options) {
		return this._map.fitBounds(bounds, options);
	}
	SetLayoutProperty(layout, property, value) {
		return this._map.setLayoutProperty(layout, property, value);
	}
	SetMarkers(markers) {
		console.log("Setting map markers", markers);
		this._markers.forEach((marker) => marker.remove());
		this._markers = markers;
		for (let m of markers) {
			m.addTo(this._map);
		}
	}
}

customElements.define("map-multipoint", MapMultipoint);
