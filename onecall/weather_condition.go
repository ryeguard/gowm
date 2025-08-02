package onecall

type weatherCondition int // Code[int64],Group[string],Description[string],Icon[string]

//go:generate goenums weather_condition.go
const (
	// Group 2xx: Thunderstorm

	thunderstormWithLightRain    weatherCondition = iota // 200,Thunderstorm,"thunderstorm with light rain",11d
	thunderstormWithRain                                 // 201,Thunderstorm,"thunderstorm with rain",11d
	thunderstormWithHeavyRain                            // 202,Thunderstorm,"thunderstorm with heavy rain",11d
	lightThunderstorm                                    // 210,Thunderstorm,"light thunderstorm",11d
	thunderstorm                                         // 211,Thunderstorm,"thunderstorm",11d
	heavyThunderstorm                                    // 212,Thunderstorm,"heavy thunderstorm",11d
	raggedThunderstorm                                   // 221,Thunderstorm,"ragged thunderstorm",11d
	thunderstormWithLightDrizzle                         // 230,Thunderstorm,"thunderstorm with light drizzle",11d
	thunderstormWithDrizzle                              // 231,Thunderstorm,"thunderstorm with drizzle",11d
	thunderstormWithHeavyDrizzle                         // 232,Thunderstorm,"thunderstorm with drizzle",11d

	// Group 3xx: Drizzle

	lightIntensityDrizzle     // 300,Drizzle,"light intensity drizzle",09d
	drizzle                   // 301,Drizzle,"drizzle",09d
	heavyIntensityDrizzle     // 310,Drizzle,"heavy intensity drizzle",09d
	lightIntensityDrizzleRain // 311,Drizzle,"light intensity drizzle rain",09d
	drizzleRain               // 312,Drizzle,"drizzle rain",09d
	showerRainAndDrizzle      // 313,Drizzle,"shower rain and drizzle",09d
	heavyShowerRainAndDrizzle // 314,Drizzle,"heavy shower rain and drizzle",09d
	showerDrizzle             // 321,Drizzle,"shower drizzle",09d

	// Group 5xx: Rain

	lightRain                // 500,Rain,"light rain",10d
	moderateRain             // 501,Rain,"moderate rain",10d
	heavyIntensityRain       // 502,Rain,"heavy intensity rain",10d
	veryHeavyRain            // 503,Rain,"very heavy rain",10d
	extremeRain              // 504,Rain,"extreme rain",10d
	freezingRain             // 511,Rain,"freezing rain",13d
	lightIntensityShowerRain // 520,Rain,"light intensity shower rain",09d
	showerRain               // 521,Rain,"shower rain",09d
	heavyIntensityShowerRain // 522,Rain,"heavy intensity shower rain",09d
	raggedShowerRain         // 531,Rain,"ragged shower rain",09d

	// Group 6xx: Snow

	lightSnow        // 600,Snow,"light snow",13d
	snow             // 601,Snow,"snow",13d
	heavySnow        // 602,Snow,"heavy snow",13d
	sleet            // 611,Snow,"sleet",13d
	lightShowerSleet // 612,Snow,"light shower sleet",13d
	showerSleet      // 613,Snow,"shower sleet",13d
	lightRainAndSnow // 615,Snow,"light rain and snow",13d
	rainAndSnow      // 616,Snow,"rain and snow",13d
	lightShowerSnow  // 620,Snow,"light shower snow",13d
	showerSnow       // 621,Snow,"shower snow",13d
	heavyShowerSnow  // 622,Snow,"heavy shower snow",13d

	// Group 7xx: Atmosphere

	mist           // 701,Mist,"mist",50d
	smoke          // 711,Smoke,"smoke",50d
	haze           // 721,Haze,"haze",50d
	sandDustWhirls // 731,Dust,"sand/dust whirls",50d
	fog            // 741,Fog,"fog",50d
	sand           // 751,Sand,"sand",50d
	dust           // 761,Dust,"dust",50d
	volcanicAsh    // 762,Ash,"volcanic ash",50d
	squalls        // 771,Squall,"squalls",50d
	tornado        // 781,Tornado,"tornado",50d

	// Group 800: Clear

	clear // 800,Clear,"clear sky",01d

	// Group 80x: Clouds

	fewClouds       // 801,Clouds,"few clouds: 11-25%",02d
	scatteredClouds // 802,Clouds,"scattered clouds: 25-50%",03d
	brokenClouds    // 803,Clouds,"broken clouds: 51-84%",04d
	overcastClouds  // 804,Clouds,"overcast clouds: 85-100%",04d
)
