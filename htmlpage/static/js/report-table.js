// report-table.js

class ReportTable extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({ mode: 'open' });
		this._reports = [];
	}

	/**
	 * Set the reports data and render the table
	 */
	set reports(value) {
		this._reports = value;
		this.render();
	}

	/**
	 * Get the reports data
	 */
	get reports() {
		return this._reports;
	}

	connectedCallback() {
		this.render();
	}

	/**
	 * Format timestamp to relative time (e.g., "2 days ago")
	 */
	formatRelativeTime(timestamp) {
		const now = new Date();
		const date = new Date(timestamp);
		const diffInSeconds = Math.floor((now - date) / 1000);
		
		// Time units in seconds
		const minute = 60;
		const hour = minute * 60;
		const day = hour * 24;
		const week = day * 7;
		const month = day * 30;
		const year = day * 365;
		
		if (diffInSeconds < minute) {
			return 'just now';
		} else if (diffInSeconds < hour) {
			const minutes = Math.floor(diffInSeconds / minute);
			return `${minutes} ${minutes === 1 ? 'minute' : 'minutes'} ago`;
		} else if (diffInSeconds < day) {
			const hours = Math.floor(diffInSeconds / hour);
			return `${hours} ${hours === 1 ? 'hour' : 'hours'} ago`;
		} else if (diffInSeconds < week) {
			const days = Math.floor(diffInSeconds / day);
			return `${days} ${days === 1 ? 'day' : 'days'} ago`;
		} else if (diffInSeconds < month) {
			const weeks = Math.floor(diffInSeconds / week);
			return `${weeks} ${weeks === 1 ? 'week' : 'weeks'} ago`;
		} else if (diffInSeconds < year) {
			const months = Math.floor(diffInSeconds / month);
			return `${months} ${months === 1 ? 'month' : 'months'} ago`;
		} else {
			const years = Math.floor(diffInSeconds / year);
			return `${years} ${years === 1 ? 'year' : 'years'} ago`;
		}
	}

	/**
	 * Get badge color class based on report type
	 */
	getTypeClass(type) {
		switch(type) {
			case 'Nuisance':
				return 'bg-danger';
			case 'Quick':
				return 'bg-primary';
			case 'Green Pool':
				return 'bg-success';
			default:
				return 'bg-secondary';
		}
	}

	/**
	 * Get badge color class based on report status
	 */
	getStatusClass(status) {
		switch(status) {
			case 'Reported':
				return 'bg-warning text-dark';
			case 'Assigned':
				return 'bg-info text-dark';
			case 'On-Hold':
				return 'bg-secondary';
			case 'Complete':
				return 'bg-success';
			default:
				return 'bg-secondary';
		}
	}

	/**
	 * Format the report ID with hyphens
	 */
	formatId(id) {
		if (id.length === 12) {
			return `${id.substring(0, 4)}-${id.substring(4, 8)}-${id.substring(8)}`;
		}
		return id;
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
						<th scope="col">Report ID</th>
						<th scope="col">Reported</th>
						<th scope="col">Type</th>
						<th scope="col">Address</th>
						<th scope="col">Status</th>
					</tr>
				</thead>
				<tbody id="report-table-body">
		`;

		// Generate rows for each report
		if (this._reports.length > 0) {
			this._reports.forEach(report => {
				const typeClass = this.getTypeClass(report.type);
				const statusClass = this.getStatusClass(report.status);
				const formattedId = this.formatId(report.id);
				const relativeTime = this.formatRelativeTime(report.created);

				tableHTML += `
					<tr class="clickable-row" onclick="window.location='/status/${report.id}'">
						<td><strong>${formattedId}</strong></td>
						<td>${relativeTime}</td>
						<td><span class="badge ${typeClass} report-type-badge">${report.type}</span></td>
						<td>${report.address || 'N/A'}</td>
						<td><span class="badge ${statusClass}">${report.status}</span></td>
					</tr>
				`;
			});
		} else {
			tableHTML += `
				<tr><td colspan="5">No reports</td></tr>
			`;
		}

		tableHTML += `
				</tbody>
			</table>
		`;

		// Set the shadow DOM content
		this.shadowRoot.innerHTML = style + tableHTML;
	}
}

// Register the custom element
customElements.define('report-table', ReportTable);
