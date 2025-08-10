package onecall

type weatherCondition int // Code[int64],Group[string],Description[string]

//go:generate goenums weather_condition.go
const (
	// Group 2xx: Thunderstorm

	thunderstormWithLightRain    weatherCondition = iota // 200,Thunderstorm,"thunderstorm with light rain"
	thunderstormWithRain                                 // 201,Thunderstorm,"thunderstorm with rain"
	thunderstormWithHeavyRain                            // 202,Thunderstorm,"thunderstorm with heavy rain"
	lightThunderstorm                                    // 210,Thunderstorm,"light thunderstorm"
	thunderstorm                                         // 211,Thunderstorm,"thunderstorm"
	heavyThunderstorm                                    // 212,Thunderstorm,"heavy thunderstorm"
	raggedThunderstorm                                   // 221,Thunderstorm,"ragged thunderstorm"
	thunderstormWithLightDrizzle                         // 230,Thunderstorm,"thunderstorm with light drizzle"
	thunderstormWithDrizzle                              // 231,Thunderstorm,"thunderstorm with drizzle"
	thunderstormWithHeavyDrizzle                         // 232,Thunderstorm,"thunderstorm with drizzle"

	// Group 3xx: Drizzle

	lightIntensityDrizzle     // 300,Drizzle,"light intensity drizzle"
	drizzle                   // 301,Drizzle,"drizzle"
	heavyIntensityDrizzle     // 310,Drizzle,"heavy intensity drizzle"
	lightIntensityDrizzleRain // 311,Drizzle,"light intensity drizzle rain"
	drizzleRain               // 312,Drizzle,"drizzle rain"
	showerRainAndDrizzle      // 313,Drizzle,"shower rain and drizzle"
	heavyShowerRainAndDrizzle // 314,Drizzle,"heavy shower rain and drizzle"
	showerDrizzle             // 321,Drizzle,"shower drizzle"

	// Group 5xx: Rain

	lightRain                // 500,Rain,"light rain"
	moderateRain             // 501,Rain,"moderate rain"
	heavyIntensityRain       // 502,Rain,"heavy intensity rain"
	veryHeavyRain            // 503,Rain,"very heavy rain"
	extremeRain              // 504,Rain,"extreme rain"
	freezingRain             // 511,Rain,"freezing rain"
	lightIntensityShowerRain // 520,Rain,"light intensity shower rain"
	showerRain               // 521,Rain,"shower rain"
	heavyIntensityShowerRain // 522,Rain,"heavy intensity shower rain"
	raggedShowerRain         // 531,Rain,"ragged shower rain"

	// Group 6xx: Snow

	lightSnow        // 600,Snow,"light snow"
	snow             // 601,Snow,"snow"
	heavySnow        // 602,Snow,"heavy snow"
	sleet            // 611,Snow,"sleet"
	lightShowerSleet // 612,Snow,"light shower sleet"
	showerSleet      // 613,Snow,"shower sleet"
	lightRainAndSnow // 615,Snow,"light rain and snow"
	rainAndSnow      // 616,Snow,"rain and snow"
	lightShowerSnow  // 620,Snow,"light shower snow"
	showerSnow       // 621,Snow,"shower snow"
	heavyShowerSnow  // 622,Snow,"heavy shower snow"

	// Group 7xx: Atmosphere

	mist           // 701,Mist,"mist"
	smoke          // 711,Smoke,"smoke"
	haze           // 721,Haze,"haze"
	sandDustWhirls // 731,Dust,"sand/dust whirls"
	fog            // 741,Fog,"fog"
	sand           // 751,Sand,"sand"
	dust           // 761,Dust,"dust"
	volcanicAsh    // 762,Ash,"volcanic ash"
	squalls        // 771,Squall,"squalls"
	tornado        // 781,Tornado,"tornado"

	// Group 800: Clear

	clear // 800,Clear,"clear sky"

	// Group 80x: Clouds

	fewClouds       // 801,Clouds,"few clouds"
	scatteredClouds // 802,Clouds,"scattered clouds"
	brokenClouds    // 803,Clouds,"broken clouds"
	overcastClouds  // 804,Clouds,"overcast clouds"
)
