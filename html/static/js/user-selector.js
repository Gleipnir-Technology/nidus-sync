class UserSelector extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({ mode: "open" });
		this.selectedUser = null;
		this.debounceTimer = null;
	}

	connectedCallback() {
		this.render();
		this.setupEventListeners();
	}

	render() {
		this.shadowRoot.innerHTML = `
			<link href="/static/css/bootstrap.css" rel="stylesheet">
			<style>
				:host {
					display: block;
					position: relative;
				}
				
				.suggestions-dropdown {
					position: absolute;
					top: 100%;
					left: 0;
					right: 0;
					z-index: 1000;
					max-height: 300px;
					overflow-y: auto;
					display: none;
				}
				
				.suggestions-dropdown.show {
					display: block;
				}
				
				.suggestion-item {
					cursor: pointer;
					border-bottom: 1px solid #dee2e6;
				}
				
				.suggestion-item:last-child {
					border-bottom: none;
				}
				
				.suggestion-item:hover {
					background-color: #f8f9fa;
				}
				
				.user-display-name {
					font-weight: 500;
					color: #212529;
				}
				
				.user-username {
					font-size: 0.875rem;
					color: #6c757d;
				}
				
				.user-org {
					font-size: 0.875rem;
					color: #6c757d;
				}
				
				.loading {
					text-align: center;
					padding: 0.75rem;
					color: #6c757d;
				}
			</style>
			
			<div class="user-selector-container">
				<input 
					type="text" 
					class="form-control" 
					placeholder="Type to search users (min. 4 characters)..."
					id="userInput"
					autocomplete="off"
				/>
				
				<div class="suggestions-dropdown card shadow-sm" id="suggestionsDropdown">
					<div class="list-group list-group-flush" id="suggestionsList">
						<!-- Suggestions will be inserted here -->
					</div>
				</div>
			</div>
		`;
	}

	setupEventListeners() {
		const input = this.shadowRoot.getElementById("userInput");
		const dropdown = this.shadowRoot.getElementById("suggestionsDropdown");

		input.addEventListener("input", (e) => this.handleInput(e));
		input.addEventListener("focus", (e) => {
			if (e.target.value.length >= 4) {
				this.handleInput(e);
			}
		});

		// Close dropdown when clicking outside
		document.addEventListener("click", (e) => {
			if (!this.contains(e.target)) {
				this.hideSuggestions();
			}
		});
	}

	handleInput(e) {
		const query = e.target.value;

		// Clear previous timer
		clearTimeout(this.debounceTimer);

		if (query.length < 4) {
			this.hideSuggestions();
			return;
		}

		// Debounce API calls
		this.debounceTimer = setTimeout(() => {
			this.fetchSuggestions(query);
		}, 300);
	}

	async fetchSuggestions(query) {
		const suggestionsList = this.shadowRoot.getElementById("suggestionsList");
		const dropdown = this.shadowRoot.getElementById("suggestionsDropdown");

		// Show loading state
		suggestionsList.innerHTML = '<div class="loading">Loading...</div>';
		dropdown.classList.add("show");

		try {
			const response = await fetch(
				`/api/user/suggestion?query=${encodeURIComponent(query)}`,
			);

			if (!response.ok) {
				throw new Error("Network response was not ok");
			}

			const data = await response.json();
			this.displaySuggestions(data.users);
		} catch (error) {
			console.error("Error fetching suggestions:", error);
			suggestionsList.innerHTML = `
				<div class="alert alert-danger m-2" role="alert">
					Error loading suggestions. Please try again.
				</div>
			`;
		}
	}

	displaySuggestions(users) {
		const suggestionsList = this.shadowRoot.getElementById("suggestionsList");
		const dropdown = this.shadowRoot.getElementById("suggestionsDropdown");

		if (!users || users.length === 0) {
			suggestionsList.innerHTML = `
				<div class="loading">No users found</div>
			`;
			return;
		}

		suggestionsList.innerHTML = users
			.map(
				(user) => `
			<div class="list-group-item list-group-item-action suggestion-item" data-user='${JSON.stringify(user)}'>
				<div class="d-flex w-100 justify-content-between align-items-start">
					<div class="flex-grow-1">
						<div class="user-display-name">${this.escapeHtml(user.display_name)}</div>
						<div class="user-username">@${this.escapeHtml(user.username)}</div>
					</div>
					<div class="text-end">
						<span class="badge bg-secondary user-org">${this.escapeHtml(user.organization.name)}</span>
					</div>
				</div>
			</div>
		`,
			)
			.join("");

		// Add click handlers to suggestion items
		suggestionsList.querySelectorAll(".suggestion-item").forEach((item) => {
			item.addEventListener("click", (e) => {
				const userData = JSON.parse(e.currentTarget.getAttribute("data-user"));
				this.selectUser(userData);
			});
		});

		dropdown.classList.add("show");
	}

	selectUser(user) {
		this.selectedUser = user;
		const input = this.shadowRoot.getElementById("userInput");
		input.value = user.displayName || user.display_name;
		this.hideSuggestions();

		// Dispatch custom event
		this.dispatchEvent(
			new CustomEvent("user-selected", {
				detail: { user },
				bubbles: true,
				composed: true,
			}),
		);
	}

	hideSuggestions() {
		const dropdown = this.shadowRoot.getElementById("suggestionsDropdown");
		dropdown.classList.remove("show");
	}

	escapeHtml(text) {
		const map = {
			"&": "&amp;",
			"<": "&lt;",
			">": "&gt;",
			'"': "&quot;",
			"'": "&#039;",
		};
		return text.replace(/[&<>"']/g, (m) => map[m]);
	}

	// Public method to get selected user
	getSelectedUser() {
		return this.selectedUser;
	}

	// Public method to clear selection
	clear() {
		this.selectedUser = null;
		const input = this.shadowRoot.getElementById("userInput");
		input.value = "";
		this.hideSuggestions();
	}
}

// Register the custom element
customElements.define("user-selector", UserSelector);
