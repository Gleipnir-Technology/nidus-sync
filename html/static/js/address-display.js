class AddressDisplay extends HTMLElement {
	constructor() {
		super();

		// Create a shadow DOM
		this.attachShadow({mode: "open" });

		// Initial render
		this.render();

		// Element references
		this._locationDisplay = this.shadowRoot.querySelector(".location-display");
		this._streetAddress = this.shadowRoot.querySelector(".street-address");
		this._postCode = this.shadowRoot.querySelector(".post-code");
		this._district = this.shadowRoot.querySelector(".district");
		this._region = this.shadowRoot.querySelector(".region");
		this._country = this.shadowRoot.querySelector(".country");
	}

	// Initial render of component
	render() {
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
			
			<div class="location-display" class="mt-4 d-none">
				<h5 class="mb-3">Location Details</h5>
				<div class="location-card p-3 mb-3">
					<div class="location-detail">
						<div class="detail-label">Street Address</div>
						<div class="street-address detail-value">-</div>
					</div>
					
					<div class="location-detail">
						<div class="detail-label">Post Code</div>
						<div class="post-code detail-value">-</div>
					</div>
					
					<div class="location-detail">
						<div class="detail-label">District/Place</div>
						<div class="district detail-value">-</div>
					</div>
					
					<div class="location-detail">
						<div class="detail-label">Region/State</div>
						<div class="region detail-value">-</div>
					</div>
					
					<div class="location-detail">
						<div class="detail-label">Country</div>
						<div class="country detail-value">-</div>
					</div>
				</div>
			</div>
		`;
	}

	// Public methods
	show(location) {
		console.log("Showing location", location);
		// Extract context data from properties
		const props = location.properties;
		const context = props.context || {};
		
		// Populate structured fields
		// Street Address - combine address, street, housenumber if available
		let addressStr = '';
		if (context.address) addressStr += context.address.address_number;
		if (context.street) {
			 if (addressStr) addressStr += ' ';
			 addressStr += context.street.name;
		}
		if (addressStr === '') {
			 addressStr = props.name || props.full_address || '-';
		}
		this._streetAddress.textContent = addressStr;
		
		// Post Code
		this._postCode.textContent = context.postcode.name || '-';
		
		// District (could be district, locality, or place)
		this._district.textContent = context.district.name || context.place.name || context.locality.name || '-';
		
		// Region (state, province, etc.)
		this._region.textContent = context.region.name || '-';
		
		// Country
		this._country.textContent = context.country.name || '-';
	}
}

customElements.define('address-display', AddressDisplay);
