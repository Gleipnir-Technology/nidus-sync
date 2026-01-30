function getGeolocation(options) {
	return new Promise((resolve, reject) => {
		// Check if geolocation is supported by the browser
		if (!navigator.geolocation) {
			reject(new Error("Geolocation is not supported by your browser"));
			return;
		}
		
		// Default options if none provided
		const geolocationOptions = options || {
			enableHighAccuracy: true,
			timeout: 5000,
			maximumAge: 0
		};
		
		// Call the geolocation API
		navigator.geolocation.getCurrentPosition(
			position => resolve(position),
			error => reject(error),
			geolocationOptions
		);
	});
}
