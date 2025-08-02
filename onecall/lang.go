package onecall

type lang int

//go:generate goenums lang.go
const (
	unknownLang lang = iota // invalid
	English                 // en
	Spanish                 // sp
	Swedish                 // sv
)
