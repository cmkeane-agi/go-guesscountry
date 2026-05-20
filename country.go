package goguessocuntry

import (
"errors"
"regexp"
"slices"
"sort"
"strings"
"sync"
"unicode"
"unicode/utf8"

"github.com/lithammer/fuzzysearch/fuzzy"
)

// Country represents a country entry.
type Country struct {
Name       string  `json:"name"`
FormalName string  `json:"formal_name"`
Abbrev     string  `json:"abbrev"`
NameAlt    string  `json:"name_alt"`
Population int     `json:"population"`
GDP        int     `json:"gdp"`
ISO2       string  `json:"iso2"`
ISO3       string  `json:"iso3"`
Continent  string  `json:"continent"`
RegionUN   string  `json:"region_un"`
Subregion  string  `json:"subregion"`
RegionWB   string  `json:"region_wb"`
LabelX     float64 `json:"label_x"`
LabelY     float64 `json:"label_y"`
}

func Iso2(code string) (Country, error) {
code = strings.ToUpper(strings.TrimSpace(code))
for _, c := range goguesscountries {
if c.ISO2 == code {
return c, nil
}
}
return Country{}, errors.New("ISO2 not matched")
}

func Iso3(code string) (Country, error) {
code = strings.ToUpper(strings.TrimSpace(code))
for _, c := range goguesscountries {
if c.ISO3 == code {
return c, nil
}
}
return Country{}, errors.New("ISO3 not matched")
}

var (
compiledRegexISO2 map[string]*regexp.Regexp
compileRegexOnce  sync.Once
)

func getCompiledRegexISO2() map[string]*regexp.Regexp {
compileRegexOnce.Do(func() {
compiledRegexISO2 = make(map[string]*regexp.Regexp, len(countryRegexISO2))
for iso2, pattern := range countryRegexISO2 {
re, err := regexp.Compile("(?i)" + pattern)
if err == nil {
compiledRegexISO2[iso2] = re
}
}
})
return compiledRegexISO2
}

func normalizeLookup(s string) string {
s = strings.TrimSpace(strings.ToLower(s))
var b strings.Builder
b.Grow(len(s))
for _, r := range s {
if unicode.IsLetter(r) || unicode.IsDigit(r) {
b.WriteRune(r)
}
}
return b.String()
}

func countryLookupFields(c Country) []string {
return []string{c.Name, c.FormalName, c.Abbrev, c.NameAlt, c.ISO2, c.ISO3}
}

// Countries returns the complete built-in country dataset.
func Countries() []Country {
return slices.Clone(goguesscountries)
}

// CountryAttribute returns a specific attribute value from a Country.
func CountryAttribute(country Country, attribute string) (any, error) {
switch strings.ToLower(strings.TrimSpace(attribute)) {
case "name", "name_short":
return country.Name, nil
case "formal_name", "name_official":
return country.FormalName, nil
case "abbrev":
return country.Abbrev, nil
case "name_alt":
return country.NameAlt, nil
case "population":
return country.Population, nil
case "gdp":
return country.GDP, nil
case "iso2":
return country.ISO2, nil
case "iso3":
return country.ISO3, nil
case "continent":
return country.Continent, nil
case "region_un":
return country.RegionUN, nil
case "subregion":
return country.Subregion, nil
case "region_wb":
return country.RegionWB, nil
case "label_x":
return country.LabelX, nil
case "label_y":
return country.LabelY, nil
default:
return nil, errors.New("unknown attribute")
}
}

// GuessCountryAttribute looks up a country and returns an attribute value,
// falling back to defaultValue if no match or invalid attribute is found.
func GuessCountryAttribute(name string, attribute string, defaultValue any) any {
country, err := GuessCountry(name)
if err != nil {
return defaultValue
}
value, err := CountryAttribute(country, attribute)
if err != nil {
return defaultValue
}
return value
}

// GuessCountry gives a best-match country for input name/ISO/historical aliases.
func GuessCountry(name string) (Country, error) {
name = strings.TrimSpace(name)
if name == "" {
return Country{}, errors.New("country not matched")
}

normalized := normalizeLookup(name)
lowerTrimmed := strings.ToLower(name)

if iso2, ok := countryAliasISO2[normalized]; ok {
if c, err := Iso2(iso2); err == nil {
return c, nil
}
}
if iso2, ok := countryAliasISO2[lowerTrimmed]; ok {
if c, err := Iso2(iso2); err == nil {
return c, nil
}
}

switch utf8.RuneCountInString(name) {
case 2:
if c, err := Iso2(name); err == nil {
return c, nil
}
case 3:
if c, err := Iso3(name); err == nil {
return c, nil
}
}

regexMap := getCompiledRegexISO2()
codes := make([]string, 0, len(regexMap))
for code := range regexMap {
codes = append(codes, code)
}
sort.Strings(codes)
for _, iso2 := range codes {
if regexMap[iso2].MatchString(name) {
if c, err := Iso2(iso2); err == nil {
return c, nil
}
}
}

for _, c := range goguesscountries {
for _, field := range countryLookupFields(c) {
if field == "" {
continue
}
if strings.EqualFold(name, field) || normalized == normalizeLookup(field) {
return c, nil
}
}
}

bestIdx := -1
bestScore := 1.0
for i, c := range goguesscountries {
for _, field := range countryLookupFields(c) {
if field == "" {
continue
}
target := normalizeLookup(field)
if target == "" {
continue
}
maxLen := max(utf8.RuneCountInString(normalized), utf8.RuneCountInString(target))
if maxLen == 0 {
continue
}
score := float64(fuzzy.LevenshteinDistance(normalized, target)) / float64(maxLen)
if strings.Contains(target, normalized) || strings.Contains(normalized, target) {
score -= 0.15
}
if score < bestScore {
bestScore = score
bestIdx = i
}
}
}

if bestIdx >= 0 && bestScore <= 0.40 {
return goguesscountries[bestIdx], nil
}

return Country{}, errors.New("country not matched")
}
