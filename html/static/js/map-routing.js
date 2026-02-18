// A test of maplibre-gl in a custom element
class MapRouting extends HTMLElement {
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
		const centroid = JSON.parse(this.getAttribute("centroid"));
		const organization_id = this.getAttribute("organization-id");
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

		/*
		this._map = new maplibregl.Map({
			center: centroid.coordinates,
			container: mapElement,
			style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json", // Style URL; see our documentation for more options
		}).fitBounds(bounds, {
			padding: { top: 10, bottom: 10, left: 10, right: 10 },
		});
		this._map.on("load", () => {
			this.dispatchEvent(new CustomEvent("load"), {
				bubbles: true,
				composed: true, // Allows event to cross shadow DOM boundary
				detail: {
					map: this,
				},
			});
		});
		*/
		this._map = new maplibregl.Map({
			center: {
				lat: 36.351947895503585,
				lng: -119.31857880996313,
			},
			container: mapElement,
			style: "https://tiles.stadiamaps.com/styles/alidade_smooth.json", // Style URL; see our documentation for more options
		}).fitBounds(
			[
				{ lat: 36.33870557056423, lng: -119.35466592321588 },
				{ lat: 36.36630172845781, lng: -119.28771302024407 },
			],
			{
				padding: { top: 10, bottom: 10, left: 10, right: 10 },
			},
		);
		const routeData = {
			type: "Feature",
			properties: {},
			geometry: {
				type: "LineString",
				coordinates: [
					[-119.31104, 36.3419],
					[-119.31005, 36.34185],
					[-119.30905, 36.34183],
					[-119.30815, 36.34181],
					[-119.30778, 36.34182],
					[-119.30755, 36.34184],
					[-119.30678, 36.34188],
					[-119.30656, 36.34188],
					[-119.30618, 36.34187],
					[-119.3056, 36.34187],
					[-119.3056, 36.34277],
					[-119.30561, 36.34345],
					[-119.3056, 36.34362],
					[-119.30562, 36.34523],
					[-119.30563, 36.34627],
					[-119.30563, 36.3473],
					[-119.30563, 36.3483],
					[-119.30566, 36.3501],
					[-119.30565, 36.35052],
					[-119.30566, 36.3508],
					[-119.30567, 36.35129],
					[-119.30567, 36.35191],
					[-119.30569, 36.35228],
					[-119.30573, 36.35276],
					[-119.30575, 36.35306],
					[-119.30574, 36.35338],
					[-119.30574, 36.35625],
					[-119.30574, 36.35641],
					[-119.30574, 36.35651],
					[-119.30572, 36.35806],
					[-119.30513, 36.35806],
					[-119.30353, 36.35805],
					[-119.30352, 36.35752],
					[-119.30393, 36.35753],
					[-119.30438, 36.35753],
					[-119.30438, 36.35753],
					[-119.3046, 36.35753],
					[-119.30512, 36.35753],
					[-119.3052, 36.35751],
					[-119.30524, 36.35746],
					[-119.30524, 36.35696],
					[-119.30521, 36.3569],
					[-119.30509, 36.35688],
					[-119.3046, 36.35688],
					[-119.30394, 36.35687],
					[-119.30308, 36.35687],
					[-119.3024, 36.35687],
					[-119.30181, 36.35687],
					[-119.30175, 36.35689],
					[-119.30173, 36.35695],
					[-119.30173, 36.35721],
					[-119.30133, 36.35721],
					[-119.30134, 36.3565],
					[-119.30191, 36.3565],
					[-119.30249, 36.3565],
					[-119.30345, 36.3565],
					[-119.30492, 36.35651],
					[-119.30509, 36.35651],
					[-119.30528, 36.35651],
					[-119.30574, 36.35651],
					[-119.30574, 36.35641],
					[-119.30574, 36.35625],
					[-119.30574, 36.35338],
					[-119.30575, 36.35306],
					[-119.30573, 36.35276],
					[-119.30569, 36.35228],
					[-119.30567, 36.35191],
					[-119.30567, 36.35129],
					[-119.30566, 36.3508],
					[-119.30565, 36.35052],
					[-119.30566, 36.3501],
					[-119.30597, 36.3501],
					[-119.30613, 36.35009],
					[-119.30629, 36.35008],
					[-119.30642, 36.35007],
					[-119.30688, 36.35001],
					[-119.30721, 36.34992],
					[-119.30754, 36.34984],
					[-119.30817, 36.34955],
					[-119.30851, 36.34946],
					[-119.30906, 36.34933],
					[-119.30917, 36.34932],
					[-119.30949, 36.34928],
					[-119.31007, 36.34928],
					[-119.31152, 36.34928],
					[-119.31195, 36.34928],
					[-119.3124, 36.34928],
					[-119.31337, 36.3493],
					[-119.31354, 36.3493],
					[-119.31374, 36.3493],
					[-119.31391, 36.3493],
					[-119.31417, 36.34932],
					[-119.31426, 36.34932],
					[-119.31456, 36.34933],
					[-119.31484, 36.34933],
					[-119.31505, 36.34933],
					[-119.31528, 36.34931],
					[-119.31654, 36.34921],
					[-119.31692, 36.3492],
					[-119.31708, 36.34921],
					[-119.31786, 36.34921],
					[-119.31867, 36.34918],
					[-119.31972, 36.34917],
					[-119.32087, 36.34918],
					[-119.32228, 36.34917],
					[-119.32246, 36.34917],
					[-119.32263, 36.34916],
					[-119.32313, 36.34915],
					[-119.32339, 36.34916],
					[-119.32375, 36.34918],
					[-119.324, 36.34917],
					[-119.3241, 36.34922],
					[-119.32555, 36.34923],
					[-119.32625, 36.34923],
					[-119.32706, 36.34922],
					[-119.32722, 36.34915],
					[-119.32777, 36.34917],
					[-119.32776, 36.34811],
					[-119.32776, 36.3475],
					[-119.32775, 36.34709],
					[-119.32772, 36.34709],
					[-119.32712, 36.34709],
					[-119.32713, 36.34759],
					[-119.32713, 36.3477],
					[-119.32708, 36.34776],
					[-119.327, 36.34782],
					[-119.327, 36.34782],
					[-119.32708, 36.34776],
					[-119.32713, 36.3477],
					[-119.32713, 36.34759],
					[-119.32712, 36.34709],
					[-119.32772, 36.34709],
					[-119.32775, 36.34709],
					[-119.32776, 36.3475],
					[-119.32776, 36.34811],
					[-119.32777, 36.34917],
					[-119.32824, 36.34917],
					[-119.32845, 36.34917],
					[-119.32885, 36.34917],
					[-119.33003, 36.34918],
					[-119.33057, 36.34918],
					[-119.33075, 36.34918],
					[-119.3309, 36.34918],
					[-119.33099, 36.34922],
					[-119.33116, 36.34922],
					[-119.33126, 36.34925],
					[-119.33195, 36.34926],
					[-119.33197, 36.34976],
					[-119.33198, 36.35],
					[-119.33199, 36.35024],
					[-119.33203, 36.35129],
					[-119.33201, 36.35191],
					[-119.33202, 36.35275],
					[-119.33202, 36.35279],
					[-119.33202, 36.353],
					[-119.33203, 36.35327],
					[-119.33204, 36.35457],
					[-119.33205, 36.35516],
					[-119.33205, 36.35532],
					[-119.33205, 36.3556],
					[-119.33205, 36.35601],
					[-119.33198, 36.35611],
					[-119.33197, 36.35633],
					[-119.33197, 36.35641],
					[-119.33197, 36.35657],
					[-119.33199, 36.35746],
					[-119.33199, 36.35756],
					[-119.33202, 36.35785],
					[-119.33203, 36.35815],
					[-119.33203, 36.35865],
					[-119.33203, 36.35903],
					[-119.3321, 36.35914],
					[-119.3321, 36.35923],
					[-119.33209, 36.35952],
					[-119.33211, 36.36154],
					[-119.33194, 36.36154],
					[-119.33114, 36.36153],
					[-119.33029, 36.36154],
					[-119.32824, 36.36153],
					[-119.32824, 36.36165],
					[-119.32826, 36.36241],
					[-119.32826, 36.36262],
					[-119.3283, 36.36284],
				],
			},
		};
		const stopData = {
			type: "Feature",
			geometry: {
				type: "MultiPoint",
				coordinates: [
					[-119.31104, 36.3419],
					[-119.30438, 36.35753],
					[-119.327, 36.34782],
					[-119.3283, 36.36284],
				],
			},
			properties: {},
		};

		// Add map controls
		this._map.addControl(new maplibregl.NavigationControl());
		// Wait for the map to load
		this._map.on("load", () => {
			this._map.addSource("route", {
				type: "geojson",
				data: routeData,
			});
			this._map.addSource("stops", {
				type: "geojson",
				data: stopData,
			});

			// Add a layer to display the route
			this._map.addLayer({
				id: "route",
				type: "line",
				source: "route",
				layout: {
					"line-join": "round",
					"line-cap": "round",
				},
				paint: {
					"line-color": "#3887be",
					"line-width": 5,
					"line-opacity": 0.75,
				},
			});

			this._map.addLayer({
				id: "stops",
				type: "circle",
				source: "stops",
				paint: {
					"circle-radius": 8,
					"circle-color": "#f00",
				},
			});
		});
	}

	// Initial render of component
	render() {
		this.shadowRoot.innerHTML = `
			<style>
				@import url('//unpkg.com/maplibre-gl@5.0.1/dist/maplibre-gl.css');
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
			draggable: true,
		})
			.setLngLat(coords)
			.addTo(map);
		marker.on("dragend", function (e) {
			const markerDraggedEvent = new CustomEvent("markerdragend", {
				detail: {
					marker: marker,
				},
			});
			mapContainer.dispatchEvent(markerDraggedEvent);
		});
		this._markers = [marker];
	}

	SetLayoutProperty(layout, property, value) {
		return this._map.setLayoutProperty(layout, property, value);
	}
}

customElements.define("map-routing", MapRouting);
