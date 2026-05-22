package goguessocuntry

import (
	"regexp"
	"testing"
)

func TestIso2(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantName string
		wantISO3 string
	}{
		{name: "uppercase code", code: "US", wantName: "United States", wantISO3: "USA"},
		{name: "lowercase code", code: "jp", wantName: "Japan", wantISO3: "JPN"},
		{name: "territory code", code: "tw", wantName: "Taiwan", wantISO3: "TWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Iso2(tt.code)
			if err != nil {
				t.Fatalf("Iso2(%q) returned error: %v", tt.code, err)
			}
			if got.Name != tt.wantName || got.ISO3 != tt.wantISO3 {
				t.Fatalf("Iso2(%q) = %s/%s, want %s/%s", tt.code, got.Name, got.ISO3, tt.wantName, tt.wantISO3)
			}
		})
	}
}

func TestIso2NotMatched(t *testing.T) {
	got, err := Iso2("XX")
	if err == nil {
		t.Fatal("Iso2(\"XX\") returned nil error")
	}
	if got != (Country{}) {
		t.Fatalf("Iso2(\"XX\") = %#v, want zero Country", got)
	}
}

func TestIso3(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantName string
		wantISO2 string
	}{
		{name: "uppercase code", code: "CAN", wantName: "Canada", wantISO2: "CA"},
		{name: "lowercase code", code: "deu", wantName: "Germany", wantISO2: "DE"},
		{name: "historical code", code: "sun", wantName: "Soviet Union (former)", wantISO2: "SU"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Iso3(tt.code)
			if err != nil {
				t.Fatalf("Iso3(%q) returned error: %v", tt.code, err)
			}
			if got.Name != tt.wantName || got.ISO2 != tt.wantISO2 {
				t.Fatalf("Iso3(%q) = %s/%s, want %s/%s", tt.code, got.Name, got.ISO2, tt.wantName, tt.wantISO2)
			}
		})
	}
}

func TestIso3NotMatched(t *testing.T) {
	got, err := Iso3("XXX")
	if err == nil {
		t.Fatal("Iso3(\"XXX\") returned nil error")
	}
	if got != (Country{}) {
		t.Fatalf("Iso3(\"XXX\") = %#v, want zero Country", got)
	}
}

func TestGuessCountry(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantISO2 string
	}{
		{name: "iso2", input: "us", wantName: "United States", wantISO2: "US"},
		{name: "iso3", input: "jpn", wantName: "Japan", wantISO2: "JP"},
		{name: "canonical name ignores case", input: "south korea", wantName: "South Korea", wantISO2: "KR"},
		{name: "formal name", input: "United States of America", wantName: "United States", wantISO2: "US"},
		{name: "alternate name", input: "Swaziland", wantName: "Eswatini", wantISO2: "SZ"},
		{name: "countryguess alias", input: "britain", wantName: "United Kingdom", wantISO2: "GB"},
		{name: "abbreviation", input: "U.S.A.", wantName: "United States", wantISO2: "US"},
		{name: "fuzzy name", input: "Untd Stts", wantName: "United States", wantISO2: "US"},
		{name: "trimmed input", input: "  Canada  ", wantName: "Canada", wantISO2: "CA"},
		{name: "historical iso2", input: "SU", wantName: "Soviet Union (former)", wantISO2: "SU"},
		{name: "historical iso4", input: "SUHH", wantName: "Soviet Union (former)", wantISO2: "SU"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GuessCountry(tt.input)
			if err != nil {
				t.Fatalf("GuessCountry(%q) returned error: %v", tt.input, err)
			}
			if got.Name != tt.wantName || got.ISO2 != tt.wantISO2 {
				t.Fatalf("GuessCountry(%q) = %s/%s, want %s/%s", tt.input, got.Name, got.ISO2, tt.wantName, tt.wantISO2)
			}
		})
	}
}

func TestGuessCountryNotMatched(t *testing.T) {
	tests := []string{"", "   ", "not a real country"}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			got, err := GuessCountry(input)
			if err == nil {
				t.Fatalf("GuessCountry(%q) returned nil error", input)
			}
			if got != (Country{}) {
				t.Fatalf("GuessCountry(%q) = %#v, want zero Country", input, got)
			}
		})
	}
}

func TestMakeISO2LUT(t *testing.T) {
	countrylut, iso2lut := makeISO2LUT()

	if got := countrylut["United States"]; got != "us" {
		t.Fatalf("countrylut[United States] = %q, want %q", got, "us")
	}
	if got := iso2lut["us"]; got != "United States" {
		t.Fatalf("iso2lut[us] = %q, want %q", got, "United States")
	}
}

func TestCountryDataHasUsableCodes(t *testing.T) {
	for _, c := range goguesscountries {
		if c.Name == "" {
			t.Fatal("country has empty Name")
		}
		if c.ISO2 == "" && c.ISO3 == "" && c.ISO4 == "" {
			t.Fatalf("%s has no ISO code", c.Name)
		}
	}
}

func TestGuessCountryAttribute(t *testing.T) {
	got, err := GuessCountryAttribute("PoRtUgAl", "iso2")
	if err != nil {
		t.Fatalf("GuessCountryAttribute returned error: %v", err)
	}
	if got != "PT" {
		t.Fatalf("GuessCountryAttribute = %q, want %q", got, "PT")
	}
}

func TestCountryDataGetWithDefault(t *testing.T) {
	defaultCountry := Country{Name: "fallback", NameShort: "fallback"}
	got, ok := NewCountryData().Get("no such country", GuessOptions{Default: &defaultCountry})
	if ok {
		t.Fatal("Get returned ok for an unknown country")
	}
	if got.Name != "fallback" {
		t.Fatalf("Get default = %#v, want fallback country", got)
	}
}

func TestCountryDataCustomRegexMap(t *testing.T) {
	got, ok := NewCountryData().Get("Mongol Uls", GuessOptions{
		RegexMap: map[string]*regexp.Regexp{
			"MN": regexp.MustCompile(`(?i)^mongol\s+uls$`),
		},
	})
	if !ok {
		t.Fatal("custom regex lookup returned no match")
	}
	if got.Name != "Mongolia" {
		t.Fatalf("custom regex lookup = %q, want Mongolia", got.Name)
	}
}

func TestIso4(t *testing.T) {
	got, err := Iso4("CSHH")
	if err != nil {
		t.Fatalf("Iso4 returned error: %v", err)
	}
	if got.Name != "Czechoslovakia, Czechoslovak Socialist Republic" || got.ISO2 != "CS" || got.ISO3 != "CSK" {
		t.Fatalf("Iso4(CSHH) = %s/%s/%s", got.Name, got.ISO2, got.ISO3)
	}
}
