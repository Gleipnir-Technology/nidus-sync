var map = null;
// A map that can be used to locate a single point by setting its location explicitly
// or by allowing the user to move a marker.
class MapAggregate extends HTMLElement {
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({mode: "open" });

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

	// Lifecycle: watch these attributes for changes
	static get observedAttributes() {
		return ["api-key", "latitude", "longitude", "organization-id", "tegola", "zoom"];
	}

	// Lifecycle: respond to attribute changes
	attributeChangedCallback(name, oldValue, newValue) {
		// Only handle if map exists and values actually changed
		if (!this._map || oldValue === newValue) return;
		
		if (name === 'api-key') {
			this._apiKey = newValue;
		}
		
		if (name === 'latitude' || name === 'longitude') {
			if (this.hasAttribute('latitude') && this.hasAttribute('longitude')) {
				const lat = Number(this.getAttribute('latitude'));
				const lng = Number(this.getAttribute('longitude'));
				this._map.setCenter([lat, lng]);
			}
		}

		if (name === 'organization-id') {
			this._organizationID = newValue;
		}

		if (name === 'tegola') {
			this._tegola = newValue;
		}

		if (name === 'zoom') {
			this._map.setZoom(Number(newValue));
		}
	}
	
	_initializeMap() {
		const apiKey = this.getAttribute("api-key");
		const lat = Number(this.getAttribute("latitude") || 36.2);
		const lng = Number(this.getAttribute("longitude") || -119.2);
		const organization_id = Number(this.getAttribute("organization-id") || 0);
		const tegola = this.getAttribute("tegola")
		const zoom = Number(this.getAttribute("zoom") || 15);

		mapboxgl.accessToken = apiKey;
		const mapElement = this.shadowRoot.querySelector("#map");
		map = new mapboxgl.Map({
			container: mapElement,
			center: {
				lat: lat,
				lng: lng,
			},
			style: 'mapbox://styles/mapbox/streets-v12', // style URL
			zoom: zoom,
		});
		map.on("load", () => {
			map.addSource('tegola', {
				'type': 'vector',
				'tiles': [
					`https://${tegola}/maps/nidus/{z}/{x}/{y}?organization_id=${organization_id}`
				]
			});
			map.addInteraction('nidus-mouseenter-interaction', {
				type: 'mouseenter',
				target: { layerId: 'nidus' },
				handler: () => {
					map.getCanvas().style.cursor = 'pointer';
				}
			});
			map.addInteraction('nidus-mouseleave-interaction', {
				type: 'mouseleave',
				target: { layerId: 'nidus' },
				handler: () => {
					map.getCanvas().style.cursor = '';
				}
			});
			map.addLayer({
				'id': 'mosquito_source',
				'type': 'fill',
				'filter': ['==', ['zoom'], ['+', 2, ['to-number', ['get', 'resolution']]]],
				'source': 'tegola',
				'source-layer': 'mosquito_source',
				'paint': {
					'fill-opacity': 0.4,
					'fill-color': '#dc3545'
				}
			});
			map.addLayer({
				'id': 'service_request',
				'type': 'fill',
				'filter': ['==', ['zoom'], ['+', 2, ['to-number', ['get', 'resolution']]]],
				'source': 'tegola',
				'source-layer': 'service_request',
				'paint': {
					'fill-opacity': 0.4,
					'fill-color': '#ffc107'
				}
			});
			map.addLayer({
				'id': 'trap',
				'type': 'fill',
				'filter': ['==', ['zoom'], ['+', 2, ['to-number', ['get', 'resolution']]]],
				'source': 'tegola',
				'source-layer': 'trap',
				'paint': {
					'fill-opacity': 0.4,
					'fill-color': '#0dcaf0'
				}
			});
			var self = this;
			map.addInteraction("tegola-click-interaction", {
				type: "click",
				target: { layerId: "mosquito_source" },
				handler: (e) => {
					const coordinates = e.feature.geometry.coordinates.slice();
					const properties = e.feature.properties;
					self.dispatchEvent(new CustomEvent("cell-click", {
						bubbles: true,
						composed: true, // Allows event to cross shadow DOM boundary
						detail: {
							cell: properties.cell
						}
					}));
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

customElements.define('map-aggregate', MapAggregate);
