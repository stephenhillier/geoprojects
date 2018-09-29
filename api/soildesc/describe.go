package soildesc

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/stephenhillier/geoprojects/backend/models"
)

// Describe is an HTTP Handler function that takes an unformatted soil description
// and returns a JSON response containing a more structured, consistent description format
func (api *Server) Describe(w http.ResponseWriter, req *http.Request) {
	desc := req.FormValue("desc")

	parsed := parseDescription(desc)

	created, err := api.db.CreateDescription(parsed)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}

	response, err := json.Marshal(created)
	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// parseDescription takes an input string, scans it for keywords and fills
// a models.Description struct type with a best guess for each category
// (primary, secondary soil etc)
//
// examples of input descriptions are "sandy gravel, very wet" or "water bearing silts".
// the output Description will contain consistent, standard terms for the primary soil type,
// secondary soil type, moisture content and consistency (loose, compact etc)
//
// TODO: this code started small, checking input against some limited cases
// Adding more cases and categories (e.g. moisture, consistency) has increased
// need for refactor.
func parseDescription(orig string) models.Description {
	d := models.Description{}
	d.Original = orig

	var singleWords []string

	for _, word := range strings.Split(orig, " ") {
		singleWords = append(singleWords, strings.Trim(word, ","))
	}

	// important description terms:
	// primary: gravel, sand, clay, silt
	// secondary: sandy, gravelly, silty, clayey, some gravel,
	// some sand, some silt, some clay, trace sand, trace gravel,
	// trace clay, trace silt

	baseType := make(map[string]string)
	baseType["gravelly"] = "gravel"
	baseType["gravels"] = "gravel"
	baseType["sandy"] = "sand"
	baseType["sands"] = "sand"
	baseType["silty"] = "silt"
	baseType["silts"] = "silt"
	baseType["clayey"] = "clay"
	baseType["clays"] = "clay"
	baseType["water bearing"] = "wet"
	baseType["water"] = "wet"

	terms := make(map[string][]string)

	// parsing a description works by brute force - words in the original description
	// are matched against the `terms` map.
	//
	// standard terminology is relatively limited, but this list could be stored
	// in a database in the future to allow adding more terms easily

	terms["primary"] = []string{"gravel", "sand", "clay", "silt"}
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
				} else if d.Secondary == "" {
					// some secondary soil types might come in the form "sand and gravel" (e.g. gravel will be secondary)
					// we can catch these while searching for primary terms
					d.Secondary = term
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

	return d

}
