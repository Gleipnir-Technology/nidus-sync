class AddressInput extends HTMLElement {
	// make element form-associated
	static formAssociated = true;

	constructor() {
		super();

		this.attachShadow({ mode: "open" });
		this.internals = this.attachInternals();
		this.render();

		// Element references
		this._input = this.shadowRoot.querySelector("input");
		this._suggestions = this.shadowRoot.querySelector(".suggestions-container");

		// Bind methods
		this._handleInput = this._handleInput.bind(this);

		// Debounce timer
		this._debounceTimer = null;

		// The suggestion data
		this._suggestionData = null;
	}

	// Lifecycle: when element is added to the DOM
	connectedCallback() {
		this._input.addEventListener("input", this._handleInput);
	}

	// Lifecycle: when element is removed from the DOM
	disconnectedCallback() {
		this._input.removeEventListener("input", this._handleInput);
	}

	// Lifecycle: watch these attributes for changes
	static get observedAttributes() {
		return ["placeholder", "api-key"];
	}

	// Lifecycle: respond to attribute changes
	attributeChangedCallback(name, oldValue, newValue) {
		if (name === "placeholder" && this._input) {
			this._input.placeholder = newValue;
		}

		if (name === "api-key") {
			this._apiKey = newValue;
		}
	}

	// Properties API
	get value() {
		return this._input ? this._input.value : "";
	}

	set value(val) {
		if (this._input) {
			this._input.value = val;
			const entries = new FormData();
			entries.append("address", val);
			this.internals.setFormValue(entries);
		}
	}

	// Private methods
	_handleInput(event) {
		const searchText = event.target.value.trim();

		// Set the form input value if they submit the form without choosing an option
		this.value = event.target.value;

		// Clear previous timer
		clearTimeout(this._debounceTimer);

		// Clear suggestions if input is less than 3 characters
		if (searchText.length < 3) {
			this._suggestions.innerHTML = "";
			return;
		}

		// Debounce API calls (wait 300ms after typing stops)
		this._debounceTimer = setTimeout(() => {
			this._fetchAddressSuggestions(searchText).then((response) => {
				this._renderSuggestions(response.features);
			});
		}, 300);
	}
	async _handleClick(gid) {
		try {
			const url = `https://api.stadiamaps.com/geocoding/v2/place_details?ids=${gid}`;
			const response = await fetch(url);
			const data = await response.json();
			const suggestion = data.features[0];
			this.SetValue(suggestion);

			// Dispatch custom event for clients of this library
			this.dispatchEvent(
				new CustomEvent("address-selected", {
					bubbles: true,
					composed: true, // Allows event to cross shadow DOM boundary
					detail: {
						location: suggestion,
					},
				}),
			);
		} catch (error) {
			console.error("Error fetching geocode of suggestion:", error);
		}
	}

	async _fetchAddressSuggestions(text) {
		try {
			//const url = `https://api.mapbox.com/search/geocode/v6/forward?q=${encodeURIComponent(text)}&access_token=${this._apiKey}`;
			const url = `https://api.stadiamaps.com/geocoding/v2/autocomplete?text=${encodeURIComponent(text)}&focus.point.lat=35&focus.point.lon=-115`;

			const response = await fetch(url);
			const data = await response.json();
			return data;
		} catch (error) {
			console.error("Error fetching geocoding suggestions:", error);
		}
	}

	_renderSuggestions(suggestions) {
		console.log("Rendering suggestions", suggestions);
		this._suggestions.innerHTML = suggestions
			.map((item, index) => {
				return `
				<div class="suggestion-item list-group-item"
					data-gid="${item.properties.gid}">
						<div class="main-address">${item.properties.name}</div>
						<div class="place-info">${item.properties.coarse_location}</div>
				</div>`;
			})
			.join("");

		// Add click listeners to suggestions
		this.shadowRoot.querySelectorAll(".suggestion-item").forEach((el) => {
			el.addEventListener("click", (e) => {
				this._handleClick(el.dataset.gid);
			});
		});
	}

	// Initial render of component
	render() {
		const placeholder = this.getAttribute("placeholder") || "Enter address";

		this.shadowRoot.innerHTML = `
			<style>
				@import url('/static/vendor/css/bootstrap.min.css');
				.detail-label {
					font-size: 0.8rem;
					text-transform: uppercase;
					color: #6c757d;
					margin-bottom: 2px;
					font-weight: 600;
				}
				.detail-value {
					font-weight: 500;
				}
				.main-address {
					font-weight: 500;
				}
				.place-info {
					font-size: 0.85rem;
					color: #6c757d;
					margin-top: 2px;
				}
				.suggestions-container {
					position: absolute;
					width: 100%;
					max-height: 300px;
					overflow-y: auto;
					z-index: 1000;
					box-shadow: 0 4px 8px rgba(0,0,0,0.1);
				}
				.suggestion-item {
					cursor: pointer;
					padding: 10px 12px;
					border-bottom: 1px solid #f0f0f0;
				}
				.suggestion-item:hover {
					background-color: #f8f9fa;
				}
			</style>
			
			<label for="addressInput" class="form-label">Enter address</label>
			<input class="form-control" id="address" maxlength="200" name="address" placeholder="${placeholder}" type="text">
			<div id="suggestions" class="suggestions-container list-group"></div>
		`;
	}

	// Public methods
	clear() {
		if (this._input) {
			this._input.value = "";
			this._suggestions.innerHTML = "";
		}
	}

	SetValue(suggestion) {
		this.value = suggestion.properties.formatted_address_line;
		this._suggestions.innerHTML = "";
	}
}

customElements.define("address-input", AddressInput);
