// A map that can show ArcGIS map tiles
class MapArcgisTile extends HTMLElement {
	static observedAttributes = ["latitude", "longitude"];
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({ mode: "open" });

		// Initial render
		this.render();

		this._map = null;
		this._markers = [];
	}

	attributeChangedCallback(name, old_value, new_value) {
		//console.log("map-arcgis-tile: attribute changed", name, old_value, new_value);
		if ((name == "latitude" || name == "longitude") && this._map != null) {
			const latitude = parseFloat(this.getAttribute("latitude"));
			const longitude = parseFloat(this.getAttribute("longitude"));
			this._map.jumpTo({
				center: [longitude, latitude],
				zoom: 19,
			});
		}
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
		const arcgis_access_token = this.getAttribute("arcgis-access-token");
		const latitude = parseFloat(this.getAttribute("latitude"));
		const longitude = parseFloat(this.getAttribute("longitude"));
		const organization_id = Number(this.getAttribute("organization-id") || 0);
		const tegola = this.getAttribute("tegola");

		const mapElement = this.shadowRoot.querySelector("#map");
		this._map = new maplibregl.Map({
			center: [longitude, latitude],
			container: mapElement,
			style: "https://tiles.stadiamaps.com/styles/osm_bright.json",
			zoom: 20,
		});
		console.log("ArcGIS token", arcgis_access_token);
		const basemap_style = maplibreArcGIS.BasemapStyle.applyStyle(this._map, {
			style: "arcgis/light-gray",
			token: arcgis_access_token,
		});
		this._map.on("load", () => {
			console.log("map-arcgis-tile loaded");
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
			if (arcgis_access_token != "") {
				this._map.addSource("flyover", {
					type: "raster",
					tiles: [
						"https://tiles.arcgis.com/tiles/pV7SH1EgRc6tpxlJ/arcgis/rest/services/TrimmedFlyover2025/MapServer/tile/{z}/{y}/{x}?token=" +
							arcgis_access_token,
					],
				});
				console.log("added arcgis tile source");
				this._map.addLayer({
					id: "flyover-layer",
					source: "flyover",
					type: "raster",
				});
				console.log("added arcgis tile layer");
			}
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
		this._map.on("click", (e) => {
			this.dispatchEvent(
				new CustomEvent("map-click", {
					bubbles: true,
					composed: true,
					detail: {
						lng: e.lngLat.lng,
						lat: e.lngLat.lat,
						map: this,
						point: e.point,
					},
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

	FitBounds(bounds, options) {
		return this._map.fitBounds(bounds, options);
	}
	SetLayoutProperty(layout, property, value) {
		return this._map.setLayoutProperty(layout, property, value);
	}
	SetMarkers(markers) {
		console.log("Setting map markers", markers);
		this._markers.forEach((marker) => marker.remove());
		this._markers = markers.map((m) => {
			return new maplibregl.Marker({
				color: "#FF0000",
				draggable: false,
			})
				.setLngLat([m.longitude, m.latitude])
				.addTo(this._map);
		});
	}
}

customElements.define("map-arcgis-tile", MapArcgisTile);
