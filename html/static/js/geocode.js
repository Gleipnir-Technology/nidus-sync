async function geocodeReverse(lngLat) {
	// curl "https://api.stadiamaps.com/geocoding/v2/reverse?point.lat=59.444351&point.lon=24.750645&api_key=YOUR-API-KEY"
	const url = `https://api.stadiamaps.com/geocoding/v2/reverse?point.lat=${lngLat.lat}&point.lon=${lngLat.lng}`;
	const response = await fetch(url);
	const data = await response.json();
	console.log("reverse geocoded to", data);
	if (data.features.length == 0) {
		console.warn("No results for reverse geocode");
		return;
	}
	return data;
}
