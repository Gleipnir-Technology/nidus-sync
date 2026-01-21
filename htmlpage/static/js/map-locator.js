var map = null;
// A map that can be used to locate a single point by setting its location explicitly
// or by allowing the user to move a marker.
class MapLocator extends HTMLElement {
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({mode: "open" });

		// Initial render
		this.render();

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

	_initializeMap() {
		console.log("Setting up the map...");
		const apiKey = this.getAttribute("api-key");
		const lat = Number(this.getAttribute('latitude') || 36.2);
		const lng = Number(this.getAttribute('longitude') || -119.2);
		const zoom = Number(this.getAttribute('zoom') || 15);

		const mapElement = this.shadowRoot.querySelector("#map");
		mapboxgl.accessToken = apiKey;
		map = new mapboxgl.Map({
			container: mapElement,
			center: {
				lat: lat,
				lng: lng,
			},
			style: 'mapbox://styles/mapbox/streets-v12', // style URL
			zoom: zoom,
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
			this.dispatchEvent(new CustomEvent('load'), {
				bubbles: true,
				composed: true, // Allows event to cross shadow DOM boundary
				detail: {
					map: this
				}
			});
		});
	}

	// Initial render of component
	render() {
		this.shadowRoot.innerHTML = `
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

	setMarker(coords) {
		console.log("Setting map marker", coords);
		this._map.jumpTo({
			center: coords,
			zoom: 14,
		});
		this._markers.forEach((marker) => marker.remove());

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
		this._markers = [marker];
	}
}

customElements.define('map-locator', MapLocator);
