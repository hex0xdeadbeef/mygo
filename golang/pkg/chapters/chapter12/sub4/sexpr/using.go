package sexpr

import "log"

func PrintUsing() {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Color           bool
		Actors          map[string]string
		Oscars          []string
		Sequel          *string
	}

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actors: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
		},

		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
		},
	}

	encoded, err := Marshal(strangelove)
	if err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}

	Print(encoded)
}
