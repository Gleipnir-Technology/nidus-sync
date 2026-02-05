var map = null;
// A map that shows multiple single point locations.
// Points have additional detail popups.
class MapMultipoint extends HTMLElement {
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({mode: "open" });

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
		return ['api-key', 'latitude', 'longitude', 'zoom'];
	}

	// Lifecycle: respond to attribute changes
	attributeChangedCallback(name, oldValue, newValue) {
		// Only handle if map exists and values actually changed
		if (!this._map || oldValue === newValue) return;
		
		if (name === 'latitude' || name === 'longitude') {
			if (this.hasAttribute('latitude') && this.hasAttribute('longitude')) {
				const lat = Number(this.getAttribute('latitude'));
				const lng = Number(this.getAttribute('longitude'));
				this._map.setCenter([lat, lng]);
			}
		}

		if (name === 'zoom') {
			this._map.setZoom(Number(newValue));
		}
	}
	
	_initializeMap() {
		const apiKey = this.getAttribute("api-key");
		const lat = Number(this.getAttribute('latitude') || 36.2);
		const lng = Number(this.getAttribute('longitude') || -119.2);
		const organization_id = Number(this.getAttribute("organization-id") || 0);
		const tegola = this.getAttribute("tegola")
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
		/*this._map.addControl(new mapboxgl.GeolocateControl({
			positionOptions: {
				enableHighAccuracy: true
			},
			trackUserLocation: true,
			showUserHeading: true
		}));
		this._map.addControl(new mapboxgl.NavigationControl());
		*/
		this._map.on("load", () => {
			this.dispatchEvent(new CustomEvent('load'), {
				bubbles: true,
				composed: true, // Allows event to cross shadow DOM boundary
				detail: {
					map: this
				}
			});
		});
	}

	async _fetchAddressSuggestions(text) {
		try {
			const url = `https://api.mapbox.com/search/geocode/v6/forward?q=${encodeURIComponent(text)}&access_token=${this._apiKey}`;
			
			const response = await fetch(url);
			const data = await response.json();
			return data;
		} catch (error) {
			console.error('Error fetching geocoding suggestions:', error);
		}
	}

	_renderSuggestions(suggestions) {
		console.log("Rendering suggestions", suggestions);
		this._suggestions.innerHTML = suggestions.map((item, index) => {
			if (item.properties.place_formatted != "") {
				return `
					<div class="suggestion-item list-group-item"
						data-index="${index}"
						data-lat="${item.geometry.coordinates[1]}"
						data-lng="${item.geometry.coordinates[0]}">
							<div class="main-address">${item.properties.name || item.properties.full_address}</div>
							<div class="place-info">${item.properties.place_formatted}</div>
					</div>`
			} else {
				return `
					<div class="suggestion-item list-group-item"
						data-index="${index}"
						data-lat="${item.coordinates.lat}"
						data-lng="${item.coordinates.lng}">
							<div class="main-address">${item.properties.name || item.properties.full_address}</div>
							<div class="place-info">${item.properties.place_formatted}</div>
					</div>`
			}
		}).join('');
		
		// Add click listeners to suggestions
		this.shadowRoot.querySelectorAll('.suggestion-item').forEach(el => {
			el.addEventListener('click', e => {
				const index = parseInt(el.dataset.index);
				const suggestion = suggestions[index];
				this.value = suggestion.properties.full_address;
				this._suggestions.innerHTML = '';
				
				// Dispatch custom event
				this.dispatchEvent(new CustomEvent('address-selected', {
					bubbles: true,
					composed: true, // Allows event to cross shadow DOM boundary
					detail: {
						location: suggestion
					}
				}));
			});
		});
	}

	// Initial render of component
	render() {
		this.shadowRoot.innerHTML = `
			<style>
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

	SetLayoutProperty(layout, property, value) {
		return this._map.setLayoutProperty(layout, property, value);
	}

}

customElements.define('map-multipoint', MapMultipoint);
