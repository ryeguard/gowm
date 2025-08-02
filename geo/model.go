package geo

type DirectResponse struct {
	Data []DirectData
}

type DirectData struct {
	Name    string  `json:"name"`    // Name of the found location
	Lat     float64 `json:"lat"`     // Geographical coordinates of the found location (latitude)
	Lon     float64 `json:"lon"`     // Geographical coordinates of the found location (longitude)
	Country string  `json:"country"` // Country of the found location
	State   string  `json:"state"`   // (where available) State of the found location

	// Name of the found location in different languages.
	// The list of names can be different for different locations.
	// Note: The key seem to be ISO 639 language codes
	LocalNames map[string]string `json:"local_names"`
}
