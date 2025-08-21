package onecall

type part int

//go:generate goenums part.go
const (
	unknownPart part = iota // invalid
	current
	minutely
	hourly
	daily
	alerts
)
