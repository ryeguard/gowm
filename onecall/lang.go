package onecall

type lang int

//go:generate goenums lang.go
const (
	unknownLang        lang = iota // invalid
	Albanian                       // sq
	Afrikaans                      // af
	Arabic                         // ar
	Azerbaijani                    // az
	Basque                         // eu
	Belarusian                     // be
	Bulgarian                      // bg
	Catalan                        // ca
	ChineseSimplified              // zh_cn
	ChineseTraditional             // zh_tw
	Croatian                       // hr
	Czech                          // cz
	Danish                         // da
	Dutch                          // nl
	English                        // en
	Finnish                        // fi
	French                         // fr
	Galician                       // gl
	German                         // de
	Greek                          // el
	Hebrew                         // he
	Hindi                          // hi
	Hungarian                      // hu
	Icelandic                      // is
	Indonesian                     // id
	Italian                        // it
	Japanese                       // ja
	Korean                         // kr
	Kurmanji                       // ku
	Latvian                        // la
	Lithuanian                     // lt
	Macedonian                     // mk
	Norwegian                      // no
	Persian                        // fa
	Polish                         // pl
	Portuguese                     // pt
	PortugueseBrazil               // pt_br
	Romanian                       // ro
	Russian                        // ru
	Serbian                        // sr
	Slovak                         // sk
	Slovenian                      // sl
	Spanish                        // sp
	Swedish                        // sv
	Thai                           // th
	Turkish                        // tr
	Ukrainian                      // ua
	Vietnamese                     // vi
	Zulu                           // zu
)
