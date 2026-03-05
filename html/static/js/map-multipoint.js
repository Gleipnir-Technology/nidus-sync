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

		this._map = null;
		// markers shown on the map. Should be none or 1, generally.
		this._markers = null;
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
		return ["api-key", "latitude", "longitude", "zoom"];
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
		const organization_id = Number(this.getAttribute("organization-id") || 0);
		const tegola = this.getAttribute("tegola");
		const zoom = Number(this.getAttribute("zoom") || 15);

		const mapElement = this.shadowRoot.querySelector("#map");
		this._map = new maplibregl.Map({
			container: mapElement,
			center: {
				lat: lat,
				lng: lng,
			},
			style: "https://tiles.stadiamaps.com/styles/osm_bright.json",
			zoom: zoom,
		});
		/*this._map.addControl(new maplibregl.GeolocateControl({
			positionOptions: {
				enableHighAccuracy: true
			},
			trackUserLocation: true,
			showUserHeading: true
		}));
		this._map.addControl(new maplibregl.NavigationControl());
		*/
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
				#map {
					height: 100%;
					width:100%;
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

customElements.define("map-multipoint", MapMultipoint);
