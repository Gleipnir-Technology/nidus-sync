class TableSite extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({ mode: "open" });
		this._sites = [];
	}

	/**
	 * Set the sites data and render the table
	 */
	set sites(value) {
		this._sites = value;
		this.render();
	}

	/**
	 * Get the sites data
	 */
	get sites() {
		return this._sites;
	}

	connectedCallback() {
		this.render();
	}

	/**
	 * Get badge color class based on report status
	 */
	getConditionClass(status) {
		switch (status) {
			case "Reported":
				return "bg-warning text-dark";
			case "Assigned":
				return "bg-info text-dark";
			case "On-Hold":
				return "bg-secondary";
			case "Complete":
				return "bg-success";
			default:
				return "bg-secondary";
		}
	}

	render() {
		// Create the styles
		const style = `
			<style>
				:host {
					display: block;
				}
				.table {
					width: 100%;
					margin-bottom: 0;
					border-collapse: collapse;
				}
				.table-light {
					background-color: #f8f9fa;
				}
				.table-hover tbody tr:hover {
					background-color: rgba(0, 0, 0, 0.075);
				}
				th, td {
					padding: 0.75rem;
					border-bottom: 1px solid #dee2e6;
					text-align: left;
				}
				.clickable-row {
					cursor: pointer;
					transition: background-color 0.15s ease-in-out;
				}
				.clickable-row:hover {
					background-color: rgba(13, 110, 253, 0.1);
				}
				.badge {
					display: inline-block;
					padding: 0.35em 0.65em;
					font-size: 0.75em;
					font-weight: 700;
					line-height: 1;
					color: #fff;
					text-align: center;
					white-space: nowrap;
					vertical-align: baseline;
					border-radius: 0.25rem;
				}
				.bg-danger {
					background-color: #dc3545;
				}
				.bg-primary {
					background-color: #0d6efd;
				}
				.bg-success {
					background-color: #198754;
				}
				.bg-warning {
					background-color: #ffc107;
				}
				.bg-info {
					background-color: #0dcaf0;
				}
				.bg-secondary {
					background-color: #6c757d;
				}
				.report-type-badge {
					font-size: 0.85rem;
				}
				.text-dark {
					color: #212529 !important;
				}
			</style>
		`;

		// Create the table
		let tableHTML = `
			<table class="table table-hover mb-0">
				<thead class="table-light">
					<tr>
						<th scope="col">Site ID</th>
						<th scope="col">Condition</th>
						<th scope="col">Address</th>
					</tr>
				</thead>
				<tbody id="report-table-body">
		`;

		// Generate rows for each report
		if (this._sites.length > 0) {
			this._sites.forEach((site) => {
				tableHTML += `
					<tr class="clickable-row" data-site-id="${site.id}">
						<td><strong>${site.id}</strong></td>
						<td>${site.condition}</td>
						<td>${site.address}</td>
					</tr>
				`;
			});
		} else {
			tableHTML += `
				<tr><td colspan="3">No sites</td></tr>
			`;
		}

		tableHTML += `
				</tbody>
			</table>
		`;

		// Set the shadow DOM content
		this.shadowRoot.innerHTML = style + tableHTML;
		// Add click handlers for the rows
		this.shadowRoot.querySelectorAll("tr.clickable-row").forEach((el) => {
			el.addEventListener("click", (e) => {
				let element = e.target;
				while (element.nodeName != "TR") {
					element = element.parentElement;
				}
				this.dispatchEvent(
					new CustomEvent("row-clicked", {
						bubbles: true,
						composed: true, // Allows event to cross shadow DOM boundary
						detail: {
							reportId: element.dataset.reportId,
						},
					}),
				);
			});
		});
	}
}

// Register the custom element
customElements.define("table-site", TableSite);
