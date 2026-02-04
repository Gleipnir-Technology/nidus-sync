class AddressOrReportInput extends HTMLElement {
	// make element form-associated
	static formAssociated = true;

	constructor() {
		super();

		this.attachShadow({mode: "open" });
		this.internals = this.attachInternals();
		this.render();

		// Element references
		this._addresses = [];
		this._input = this.shadowRoot.querySelector("input");
		this._reports = [];
		this._suggestionsContainer = this.shadowRoot.querySelector(".suggestions-container");
		
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
		this._input.removeEventListener('input', this._handleInput);
	}

	// Lifecycle: watch these attributes for changes
	static get observedAttributes() {
		return ['placeholder', 'api-key'];
	}

	// Lifecycle: respond to attribute changes
	attributeChangedCallback(name, oldValue, newValue) {
		if (name === 'placeholder' && this._input) {
			this._input.placeholder = newValue;
		}
		
		if (name === 'api-key') {
			this._apiKey = newValue;
		}
	}
	
	// Properties API
	get value() {
		return this._input ? this._input.value : '';
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

		// Clear previous timer
		clearTimeout(this._debounceTimer);
		
		// Clear suggestions if input is less than 3 characters
		if (searchText.length < 3) {
			this._suggestionsContainer.innerHTML = '';
			return;
		}
		
		// Debounce API calls (wait 300ms after typing stops)
		this._debounceTimer = setTimeout(() => {
			this._handleSuggestions(searchText);
		}, 300);
		
	}

	async _fetchAddressSuggestions(text) {
		try {
			const url = `https://api.mapbox.com/search/geocode/v6/forward?q=${encodeURIComponent(text)}&access_token=${this._apiKey}`;
			
			const response = await fetch(url);
			const data = await response.json();
			return data.features;
		} catch (error) {
			console.error('Error fetching geocoding suggestions:', error);
		}
	}

	async _fetchReportSuggestions(text) {
		try {
			const url = `/report/suggest?r=${text}`
			const response = await fetch(url);
			const data = await response.json();
			return data.reports;
		} catch (error) {
			console.error("Error fetching report suggestions:", error);
		}
	}

	async _handleSuggestions(text) {
		await Promise.all([
			(async() => {
				this._addresses = await this._fetchAddressSuggestions(text);
			})(),
			(async() => {
				this._reports = await this._fetchReportSuggestions(text);
			})(),
		]);
		this._renderSuggestions(this._addresses, this._reports);
	}

	_renderSuggestions(addresses, reports) {
		console.log("Rendering suggestions", addresses, reports);
		const reportElements = reports.map((item, index) => {
			const formatted_id = _formatReportID(item.id);
			const type_display = _formatReportType(item.type);
			return `
				<div class="suggestion-item list-group-item"
					data-index="${index}"
					data-report-id="${item.id}"
					data-type="report">
						<div class="report-id">${formatted_id}</div>
						<div class="report-type">${type_display}</div>
				</div>`
		}).join("");
		const addressElements = addresses.map((item, index) => {
			if (item.properties.place_formatted != "") {
				return `
					<div class="suggestion-item list-group-item"
						data-index="${index}"
						data-lat="${item.geometry.coordinates[1]}"
						data-lng="${item.geometry.coordinates[0]}"
						data-type="address">
							<div class="main-address">${item.properties.name || item.properties.full_address}</div>
							<div class="place-info">${item.properties.place_formatted}</div>
					</div>`
			} else {
				return `
					<div class="suggestion-item list-group-item"
						data-index="${index}"
						data-lat="${item.coordinates.lat}"
						data-lng="${item.coordinates.lng}"
						data-type="address">
							<div class="main-address">${item.properties.name || item.properties.full_address}</div>
							<div class="place-info">${item.properties.place_formatted}</div>
					</div>`
			}
		}).join("");
		this._suggestionsContainer.innerHTML = reportElements + addressElements;
		// Add click listeners to suggestions
		this.shadowRoot.querySelectorAll('.suggestion-item').forEach(el => {
			el.addEventListener('click', e => {
				const type = el.dataset.type;
				if (type == "report") {
					const index = parseInt(el.dataset.index);
					const report = this._reports[index];
					this.value = _formatReportID(report.id);
					this._suggestionsContainer.innerHTML = "";
				} else if (type == "address") {
					const index = parseInt(el.dataset.index);
					const address = this._addresses[index];
					this.SetValue(suggestion);
					// Dispatch custom event
					this.dispatchEvent(new CustomEvent('address-selected', {
						bubbles: true,
						composed: true, // Allows event to cross shadow DOM boundary
						detail: {
							location: suggestion
						}
					}));
				}
			});
		});
	}

	// Initial render of component
	render() {
		const placeholder = this.getAttribute('placeholder') || 'Enter address';
		
		this.shadowRoot.innerHTML = `
			<style>
				@import url('/static/css/bootstrap.css');
				@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css');
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
					top: 48px;
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
			
			<div class="input-group">
				<span class="input-group-text"><i class="fas fa-search"></i></span>
				<input type="text" class="form-control form-control-lg" id="addressSearch" name="address" placeholder="${placeholder}">
				<div id="suggestions" class="suggestions-container list-group"></div>
			</div>
		`;
	}

	// Public methods
	clear() {
		if (this._input) {
			this._input.value = '';
			this._suggestionsContainer.innerHTML = '';
		}
	}

	SetValue(suggestion) {
		this.value = suggestion.properties.full_address;
		this._suggestionsContainer.innerHTML = '';
	}
}

function _formatReportID(id) {
	if (id.length === 12) {
		return `${id.substring(0, 4)}-${id.substring(4, 8)}-${id.substring(8)}`;
	}
	return id;
}

function _formatReportType(type) {
	if (type == "nuisance") {
		return "Mosquito Nuisance Report";
	} else if (type == "pool") {
		return "Standing Water Report";
	} else {
		return "Unknown Report Type";
	}
}

customElements.define('address-or-report-input', AddressOrReportInput);
