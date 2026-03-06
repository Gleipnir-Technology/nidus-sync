/**
 * Custom HTML element <time-relative> that displays relative time
 * Usage: <time-relative time="2024-01-01T12:00:00Z"></time-relative>
 */

class TimeRelative extends HTMLElement {
	constructor() {
		super();
		this.span = null;
	}

	static get observedAttributes() {
		return ["time"];
	}

	connectedCallback() {
		// Create the span element if it doesn't exist
		if (!this.span) {
			this.span = document.createElement("span");
			this.span.className = "time-relative";
			this.appendChild(this.span);
		}
		this.updateTime();
	}

	attributeChangedCallback(name, oldValue, newValue) {
		if (name === "time" && oldValue !== newValue) {
			this.updateTime();
		}
	}

	updateTime() {
		if (this.span) {
			const timeValue = this.getAttribute("time");
			if (timeValue) {
				this.span.textContent = this.formatRelativeTime(timeValue);
			}
		}
	}

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
			return "just now";
		} else if (diffInSeconds < hour) {
			const minutes = Math.floor(diffInSeconds / minute);
			return `${minutes} ${minutes === 1 ? "minute" : "minutes"} ago`;
		} else if (diffInSeconds < day) {
			const hours = Math.floor(diffInSeconds / hour);
			return `${hours} ${hours === 1 ? "hour" : "hours"} ago`;
		} else if (diffInSeconds < week) {
			const days = Math.floor(diffInSeconds / day);
			return `${days} ${days === 1 ? "day" : "days"} ago`;
		} else if (diffInSeconds < month) {
			const weeks = Math.floor(diffInSeconds / week);
			return `${weeks} ${weeks === 1 ? "week" : "weeks"} ago`;
		} else if (diffInSeconds < year) {
			const months = Math.floor(diffInSeconds / month);
			return `${months} ${months === 1 ? "month" : "months"} ago`;
		} else {
			const years = Math.floor(diffInSeconds / year);
			return `${years} ${years === 1 ? "year" : "years"} ago`;
		}
	}

	// Property getter and setter for JavaScript access
	get time() {
		return this.getAttribute("time");
	}

	set time(value) {
		if (value) {
			this.setAttribute("time", value);
		} else {
			this.removeAttribute("time");
		}
	}
}

// Register the custom element
customElements.define("time-relative", TimeRelative);
