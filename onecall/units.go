package onecall

type units int

//go:generate goenums units.go
const (
	unknownUnits units = iota // invalid
	standard
	metric
	imperial
)
