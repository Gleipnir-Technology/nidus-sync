var map = null;
// A map that just shows a bunch of markers, it can't change them
class MapWithMarkers extends HTMLElement {
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({mode: "open" });

		// Initial render
		this.render();

		this._map = null;

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
		const lat = Number(this.getAttribute('latitude') || 36.2);
		const lng = Number(this.getAttribute('longitude') || -119.2);
		const zoom = Number(this.getAttribute('zoom') || 15);

		mapboxgl.accessToken = apiKey;
		const mapElement = this.shadowRoot.querySelector("#map");
		this._map = new mapboxgl.Map({
			container: mapElement,
			center: {
				lat: lat,
				lng: lng,
			},
			style: 'mapbox://styles/mapbox/streets-v12', // style URL
			zoom: zoom,
		});
		this._map.on("load", () => {
			console.log("map loaded");
			this.dispatchEvent(new CustomEvent('load'), {
				bubbles: true,
				composed: true, // Allows event to cross shadow DOM boundary
				detail: {
					map: this
				}
			});
		});
		this._markers = [];
	}

	// Initial render of component
	render() {
		this.shadowRoot.innerHTML = `
			<link href='https://api.tiles.mapbox.com/mapbox-gl-js/v0.53.0/mapbox-gl.css' rel='stylesheet' />
			<style>
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

	jumpTo(args) {
		this._map.jumpTo(args);
	}

	clearMarkers() {
		this._markers.forEach((marker) => marker.remove());
	}
	addMarker(coords, color) {
		console.log("Add marker", coords, color);
		const el = document.createElement("div");
		el.id = "marker";
		const marker = new mapboxgl.Marker({
			color: color,
			scale: 1.5,
		}).setLngLat(coords).addTo(this._map);
		this._markers.push(marker);
	}
}

customElements.define('map-with-markers', MapWithMarkers);
