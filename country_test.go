package goguessocuntry

import "testing"

func TestGuessCountryISOAndNames(t *testing.T) {
	iso2, err := GuessCountry("gb")
	if err != nil {
		t.Fatalf("expected ISO2 lookup to succeed: %v", err)
	}
	if iso2.Name != "United Kingdom" {
		t.Fatalf("expected United Kingdom, got %q", iso2.Name)
	}

	iso3, err := GuessCountry("twn")
	if err != nil {
		t.Fatalf("expected ISO3 lookup to succeed: %v", err)
	}
	if iso3.Name != "Taiwan" {
		t.Fatalf("expected Taiwan, got %q", iso3.Name)
	}

	formal, err := GuessCountry("Czech Republic")
	if err != nil {
		t.Fatalf("expected formal-name lookup to succeed: %v", err)
	}
	if formal.Name != "Czechia" {
		t.Fatalf("expected Czechia, got %q", formal.Name)
	}
}

func TestGuessCountryAliases(t *testing.T) {
	country, err := GuessCountry("britain")
	if err != nil {
		t.Fatalf("expected alias lookup to succeed: %v", err)
	}
	if country.ISO2 != "GB" {
		t.Fatalf("expected GB, got %q", country.ISO2)
	}

	taiwan, err := GuessCountry("tw")
	if err != nil {
		t.Fatalf("expected alias lookup to succeed: %v", err)
	}
	if taiwan.Name != "Taiwan" {
		t.Fatalf("expected Taiwan, got %q", taiwan.Name)
	}
}

func TestGuessCountryNotFound(t *testing.T) {
	_, err := GuessCountry("no such country")
	if err == nil {
		t.Fatal("expected error for unknown country")
	}
}

func TestGuessCountryAttributeWithDefault(t *testing.T) {
	got := GuessCountryAttribute("PoRtUgAl", "iso2", "XX")
	iso2, ok := got.(string)
	if !ok {
		t.Fatalf("expected string attribute, got %T", got)
	}
	if iso2 != "PT" {
		t.Fatalf("expected PT, got %q", iso2)
	}

	fallback := GuessCountryAttribute("no such country", "iso2", "XX")
	if fallback != "XX" {
		t.Fatalf("expected default value XX, got %#v", fallback)
	}

	invalidAttr := GuessCountryAttribute("Portugal", "does_not_exist", "N/A")
	if invalidAttr != "N/A" {
		t.Fatalf("expected default for invalid attribute, got %#v", invalidAttr)
	}
}

func TestCountriesReturnsCopy(t *testing.T) {
	countries := Countries()
	if len(countries) == 0 {
		t.Fatal("expected non-empty country list")
	}

	original := goguesscountries[0].Name
	countries[0].Name = "Mutated"

	if goguesscountries[0].Name != original {
		t.Fatal("expected Countries() to return a copy, not mutate internal data")
	}
}
