// A map that can be used to locate a single point by setting its location explicitly
// or by allowing the user to move a marker.
class MapLocator extends HTMLElement {
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({ mode: "open" });

		// Initial render
		this.render();

		// markers shown on the map. Should be none or 1, generally.
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
		console.log("Setting up the map...");
		const apiKey = this.getAttribute("api-key");
		const lat = Number(this.getAttribute("latitude") || 36.2);
		const lng = Number(this.getAttribute("longitude") || -119.2);
		const zoom = Number(this.getAttribute("zoom") || 15);

		const mapElement = this.shadowRoot.querySelector("#map");
		this._map = new maplibregl.Map({
			container: mapElement,
			center: {
				lat: lat,
				lng: lng,
			},
			style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json",
			zoom: zoom,
		});
		/*
		map.addControl(new maplibregl.GeolocateControl({
			positionOptions: {
				enableHighAccuracy: true
			},
			trackUserLocation: true,
			showUserHeading: true
		}));
		map.addControl(new maplibregl.NavigationControl());
		*/
		this._map.on("click", (e) => {
			e.preventDefault();
			console.log("internal click", e);
			this.dispatchEvent(
				new CustomEvent("click", {
					bubbles: true,
					composed: true, // Allows event to cross shadow DOM boundary
					detail: {
						lngLat: e.lngLat,
					},
				}),
			);
		});
		this._map.on("load", () => {
			console.log("map loaded");
			this.dispatchEvent(
				new CustomEvent("load", {
					bubbles: true,
					composed: true, // Allows event to cross shadow DOM boundary
					detail: {
						map: this,
					},
				}),
			);
		});
		this._map.on("zoomend", (e) => {
			this.dispatchEvent(
				new CustomEvent("zoomend", {
					bubbles: true,
					composed: true,
					detail: e,
				}),
			);
		});
	}

	// Initial render of component
	render() {
		this.shadowRoot.innerHTML = `
			<style>
				@import url("//unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css");
				#map {
					height: 100%;
					width:100%;
				}
				#map img {
					max-width: none;
					min-width: 0px;
					height: auto;
				}
			</style>
			
			<div id="map"></div>
		`;
	}

	GetZoom() {
		return this._map.getZoom();
	}

	JumpTo(args) {
		this._map.jumpTo(args);
	}

	PanTo(coords, options) {
		this._map.panTo(coords, options);
	}

	SetMarker(coords) {
		console.log("Setting map marker", coords);
		this._markers.forEach((marker) => marker.remove());

		const marker = new maplibregl.Marker({
			color: "#FF0000",
			draggable: true,
		})
			.setLngLat(coords)
			.addTo(this._map);
		marker.on("dragend", (e) => {
			const markerDraggedEvent = new CustomEvent("markerdragend", {
				detail: {
					marker: marker,
				},
			});
			this.dispatchEvent(markerDraggedEvent);
		});
		this._markers = [marker];
	}
}

customElements.define("map-locator", MapLocator);
