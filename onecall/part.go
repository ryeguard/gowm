package onecall

type exclude int

//go:generate goenums exclude.go
const (
	unknownExclude exclude = iota // invalid
	current
	minutely
	hourly
	daily
	alerts
)
