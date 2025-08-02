package onecall

type unit int

//go:generate goenums unit.go
const (
	unknownUnit unit = iota // invalid
	standard
	metric
	imperial
)
