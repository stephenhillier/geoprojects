package soildesc

import (
	"errors"
	"regexp"
	"strings"
)

// Description contains information about a parsed soil visual description
type Description struct {
	Original    string   `json:"original" db:"original"`
	Ordered     []string `json:"ordered" db:"ordered"`
	Primary     string   `json:"primary" db:"primary"`
	Secondary   string   `json:"secondary" db:"secondary"`
	Consistency string   `json:"consistency" db:"consistency"`
	Moisture    string   `json:"moisture" db:"moisture"`
}

// Token represents a single term (soil type, color, hardness etc)
// found in a description.
type Token struct {
	Original    string
	Modifier    string
	Position    int
	Capitalized bool
	Term        string
	Class       string
}

type soilToken struct {
	Original    string
	Modifier    string
	Capitalized bool
	Term        string
	Class       string
}

// ParseDescription takes an input string, scans it for keywords and fills
// a description struct type with a best guess for each category
// (primary, secondary soil etc)
//
// examples of input descriptions are "sandy gravel, very wet" or "water bearing silts".
// the output Description will contain consistent, standard terms for the primary soil type,
// secondary soil type, moisture content and consistency (loose, compact etc)
//
// TODO: this code started small, checking input against some limited cases
// Adding more cases and categories (e.g. moisture, consistency) has increased
// need for refactor.
func ParseDescription(orig string) (Description, error) {
	d := Description{}
	d.Original = orig

	singleWords := splitWords(orig)

	// important description terms:
	// primary: gravel, sand, clay, silt
	// secondary: sandy, gravelly, silty, clayey, some gravel,
	// some sand, some silt, some clay, trace sand, trace gravel,
	// trace clay, trace silt

	baseType := map[string]string{
		"gravelly":      "gravel",
		"gravels":       "gravel",
		"sandy":         "sand",
		"sands":         "sand",
		"silty":         "silt",
		"silts":         "silt",
		"clayey":        "clay",
		"clays":         "clay",
		"water bearing": "wet",
		"water":         "wet",
	}

	terms := make(map[string][]string)

	// parsing a description works by brute force - words in the original description
	// are matched against the `terms` map.
	//
	// standard terminology is relatively limited, but this list could be stored
	// in a database in the future to allow adding more terms easily

	terms["primary"] = []string{
		"gravel",
		"sand",
		"clay",
		"silt",
		"hardpan",
		"soil",
		"bedrock",
		"muskeg",
		"topsoil",
		"mudstone",
		"granite",
		"conglomerate",
		"granodiorite",
		"basalt",
		"sandstone",
		"shale",
		"boulders",
		"cobbles",
		"gravels",
		"mud",
		"till",
		"rock",
		"gneiss",
		"quartz",
		"quartzite",
		"limestone",
		"pebbles",
		"organics",
		"volcanics",
		"feldspar"}

	terms["secondary"] = []string{
		"sandy",
		"gravelly",
		"silty",
		"clayey",
		"some sand",
		"some gravel",
		"some silt",
		"some clay",
		"trace sand",
		"trace gravel",
		"trace silt",
		"trace clay",
	}

	// consistency terms (firmness/looseness of material)
	terms["consistency"] = []string{"loose", "soft", "firm", "compact", "hard", "dense"}

	// moisture content terms
	// some terms will be converted to a more "standard" one (e.g. "water bearing" will beocme "wet")
	// via the baseType map
	terms["moisture"] = []string{"very dry", "very wet", "water bearing", "water", "dry", "damp", "moist", "wet"}

	var prev string
	var soil string
	var moisture string

	for _, word := range singleWords {
		// determine primary constituent before moving on to other properties
	primary:
		for _, term := range terms["primary"] {

			// select first matching term and check that it is not part of "some gravel", "trace silt" etc.
			if (word == term || word == term+"s") && prev != "some" && prev != "trace" {
				if d.Primary == "" && prev != "and" && prev != "&" {
					d.Primary = term
					d.Ordered = append([]string{term}, d.Ordered...)
				} else if d.Secondary == "" {
					// some secondary soil types might come in the form "sand and gravel" (e.g. gravel will be secondary)
					// we can catch these while searching for primary terms
					d.Secondary = term
					d.Ordered = append(d.Ordered, term)
					break primary
				}
			}
		}
		prev = word
	}

	prev = "" // reset prev to an empty string before iterating again

	for _, word := range singleWords {
		// determine secondary constituent(s) for -ly terms (gravelly etc.)
		for _, term := range terms["secondary"] {
			if word == term || prev+" "+word == term {
				// words like "gravelly" need to be converted
				soil = baseType[word]
				// if soil is not in baseType map, default to current word
				if soil == "" {
					soil = word
				}

				if d.Secondary == "" {
					d.Secondary = soil
					d.Ordered = append(d.Ordered, soil)
				} else {
					d.Ordered = append(d.Ordered, soil)
				}
			}
		}

		if d.Consistency == "" {
		consistency:
			// determine consistency
			for _, term := range terms["consistency"] {
				if word == term {
					d.Consistency = word
					break consistency
				}
			}
		}

		if d.Moisture == "" {
		moisture:
			for _, term := range terms["moisture"] {
				if word == term || prev+" "+word == term {
					moisture = baseType[term]

					// if soil is not in baseType map, default to current word
					if moisture == "" {
						d.Moisture = term
					} else {
						d.Moisture = moisture
					}
					break moisture
				}
			}
		}

		prev = word

	}

	return d, nil
}

func splitWords(orig string) []string {

	var singleWords []string

	r := strings.NewReplacer(
		",", " ",
		"-", " ",
		":", " ",
		"\\", " ",
		"/", " ",
		";", " ",
		".", " ",
		"?", " ",
		"(", " ",
		")", " ",
	)

	orig = r.Replace(orig)

	spaces := regexp.MustCompile(`\s+`)
	orig = spaces.ReplaceAllString(orig, " ")

	for _, word := range strings.Split(orig, " ") {
		singleWords = append(singleWords, word)
	}
	return singleWords
}

// ParseSoilTerms parses a soil description and returns key soil terms in order of importance
func ParseSoilTerms(description string) []string {

	words := splitWords(description)

	// capitalization may be significant...but some data may arrive completely capitalized
	// if every letter is capitalized, we will turn off ranking by capitilization
	inputIsAllCaps := true

	tokens := []Token{}

	prev := ""

	for _, word := range words {
		token := Token{}
		token.Capitalized = false

		if isUpper(word) {
			token.Capitalized = true
		} else {
			inputIsAllCaps = false
		}

		soilToken, err := classifySoil(word, prev)

		prev = word

		// if classifySoil returned an error, stop processing this term
		if err != nil {
			continue
		}

		token.Original = soilToken.Original
		token.Term = soilToken.Term
		token.Modifier = soilToken.Modifier
		token.Class = soilToken.Class
		tokens = append(tokens, token)
	}

	var soilTerms []string

	for _, token := range sortSoils(tokens, !inputIsAllCaps) {
		soilTerms = append(soilTerms, token.Term)
	}

	return soilTerms

}

func isUpper(word string) bool {
	return word == strings.ToUpper(word) && word != strings.ToLower(word)
}

// classifySoil takes a word from a soil description (along with the previous word in the string)
// and returns a soilToken containing the matched term and modifier (if prev is a modifier such as "trace" or "some")
// and the class of material (soil or bedrock)
func classifySoil(word string, prev string) (soilToken, error) {
	word = strings.ToLower(word)
	prev = strings.ToLower(prev)

	if (prev != "some") && (prev != "trace") && (prev != "and") {
		prev = ""
	}

	// generate a list of the original term, plus some extra terms to try (plural, -y and -ey suffixes removed)
	wordList := [][]string{
		[]string{word, prev},
		[]string{strings.TrimSuffix(word, "s"), prev},
		[]string{strings.TrimSuffix(word, "y"), "y"},
		[]string{strings.TrimSuffix(word, "ey"), "y"},
		[]string{strings.TrimSuffix(word, "ly"), "y"},
	}

	// try each word variant in order, returning the first successful match against
	// a list of valid terms.
	// we use a switch with soil and bedrock cases split up in order to further
	// classify a term as one of those two classes
	for _, singleWord := range wordList {
		switch singleWord[0] {
		case
			// soil terms
			"gravel",
			"sand",
			"clay",
			"silt",
			"hardpan",
			"soil",
			"muskeg",
			"topsoil",
			"mud",
			"till",
			"organic",
			"boulder",
			"cobble":
			return soilToken{
				Original: word,
				Term:     singleWord[0],
				Modifier: singleWord[1],
				Class:    "soil",
			}, nil
		case
			// bedrock terms
			"bedrock",
			"mudstone",
			"granite",
			"conglomerate",
			"granodiorite",
			"basalt",
			"sandstone",
			"shale",
			"rock",
			"gneiss",
			"quartz",
			"quartzite",
			"limestone",
			"pebbles",
			"volcanics",
			"feldspar":
			return soilToken{
				Original: word,
				Term:     singleWord[0],
				Modifier: singleWord[1],
				Class:    "bedrock",
			}, nil
		}
	}

	return soilToken{}, errors.New("No matching soil type found")
}

// sortSoils takes a list of soil Tokens and re-sorts them to make sure that
// capitalized soils come first, followed by unmodified/unqualified soils, which
// are in turn followed by qualified soils ("trace clay" etc.).
// Input should be in the order of appearance in the original description.
func sortSoils(soils []Token, capitalization bool) []Token {
	newList := []Token{}

	// put soils into groups
	// groupA contains soils that were capitalized (except where all terms were capitalized)
	// groupB contains terms that were determined to be "primary" constituents
	// (were not modifed by "trace", "some", or an "-ey" suffix like "clayey")
	// groupC contains other terms (mainly terms that were modified/qualified or did not fit the other criteria)
	// This allows terms to be written slightly out of order (as is often the case like in "silty sand, trace gravel").
	groupA := []Token{}
	groupB := []Token{}
	groupC := []Token{}

	for _, token := range soils {
		if token.Capitalized && capitalization && (token.Modifier == "" || token.Modifier == "and") {
			groupA = append(groupA, token)
		} else if token.Modifier == "" || token.Modifier == "and" {
			groupB = append(groupB, token)
		} else {
			groupC = append(groupC, token)
		}
	}

	newList = append(newList, groupA...)
	newList = append(newList, groupB...)
	newList = append(newList, groupC...)
	return newList
}
