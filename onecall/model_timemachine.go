package onecall

type TimeMachineResponseRaw struct {
	oneCallResponseCommon
	Data []TimeMachineDataRaw `json:"data"`
}

type TimeMachineDataRaw struct {
	Dt         int64
	Sunrise    int64
	Sunset     int64
	Temp       *float64
	FeelsLike  *float64
	Pressure   *int
	Humidity   *int
	DewPoint   *float64
	UVI        *float64
	Clouds     *int
	Visibility *int
	WindSpeed  *float64
	WindDeg    *int
}
