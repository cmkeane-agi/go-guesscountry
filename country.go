package goguessocuntry

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
)

const fuzzyCutoff = 0.8

// Country contains the countryguess data fields plus the legacy field names
// from this Go package.
type Country struct {
	Name       string  `json:"name,omitempty"`
	FormalName string  `json:"formal_name,omitempty"`
	Abbrev     string  `json:"abbrev,omitempty"`
	NameAlt    string  `json:"name_alt,omitempty"`
	Population int     `json:"population,omitempty"`
	GDP        int     `json:"gdp,omitempty"`
	ISO2       string  `json:"iso2"`
	ISO3       string  `json:"iso3"`
	Continent  string  `json:"continent"`
	RegionUN   string  `json:"region_un,omitempty"`
	Subregion  string  `json:"subregion,omitempty"`
	RegionWB   string  `json:"region_wb,omitempty"`
	LabelX     float64 `json:"label_x,omitempty"`
	LabelY     float64 `json:"label_y,omitempty"`

	NameShort      string `json:"name_short"`
	NameOfficial   string `json:"name_official"`
	Regex          string `json:"regex"`
	ISONumeric     string `json:"isonumeric"`
	UNCode         string `json:"uncode"`
	FAOCode        string `json:"faocode"`
	GBDCode        string `json:"gbdcode"`
	Continent7     string `json:"continent_7"`
	UNRegion       string `json:"unregion"`
	EXIO1          string `json:"exio1"`
	EXIO2          string `json:"exio2"`
	EXIO3          string `json:"exio3"`
	EXIO13L        string `json:"exio1_3l"`
	EXIO23L        string `json:"exio2_3l"`
	EXIO33L        string `json:"exio3_3l"`
	WIOD           string `json:"wiod"`
	Eora           string `json:"eora"`
	Message        string `json:"message"`
	Image          string `json:"image"`
	Remind         string `json:"remind"`
	OECD           string `json:"oecd"`
	EU             string `json:"eu"`
	EU28           string `json:"eu28"`
	EU27           string `json:"eu27"`
	EU272007       string `json:"eu27_2007"`
	EU25           string `json:"eu25"`
	EU15           string `json:"eu15"`
	EU12           string `json:"eu12"`
	EEA            string `json:"eea"`
	Schengen       string `json:"schengen"`
	Euro           string `json:"euro"`
	UN             string `json:"un"`
	UNMember       string `json:"unmember"`
	Obsolete       string `json:"obsolete"`
	Cecilia2050    string `json:"cecilia2050"`
	BRIC           string `json:"bric"`
	APEC           string `json:"apec"`
	BASIC          string `json:"basic"`
	CIS            string `json:"cis"`
	G7             string `json:"g7"`
	G20            string `json:"g20"`
	IEA            string `json:"iea"`
	IEAV2025       string `json:"iea_v2025"`
	DACCode        string `json:"daccode"`
	CCTLD          string `json:"cctld"`
	GWCode         string `json:"gwcode"`
	CC41           string `json:"cc41"`
	IOC            string `json:"ioc"`
	FIFA           string `json:"fifa"`
	GEONumeric     string `json:"geonumeric"`
	COE            string `json:"coe"`
	ISO4           string `json:"iso4,omitempty"`
	WithdrawalDate string `json:"withdrawal_date,omitempty"`
	Historical     bool   `json:"historical,omitempty"`
	Comment        string `json:"comment,omitempty"`
}

// CountryData mirrors the Python package's CountryData object for Go callers.
type CountryData struct {
	countries []Country
	once      sync.Once
	regexes   []compiledRegex
}

type compiledRegex struct {
	country Country
	re      *regexp.Regexp
}

// GuessOptions controls extended country lookup behavior.
type GuessOptions struct {
	Attribute string
	Default   *Country
	RegexMap  map[string]*regexp.Regexp
}

// NewCountryData returns a CountryData instance backed by the packaged data.
func NewCountryData() *CountryData {
	return &CountryData{countries: cloneCountries(goguesscountries)}
}

// NewCountryDataFromFile returns a CountryData instance backed by a JSON file
// containing a list of Country-compatible objects.
func NewCountryDataFromFile(filepath string) (*CountryData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	countries, err := unmarshalCountries(data)
	if err != nil {
		return nil, err
	}

	return &CountryData{countries: countries}, nil
}

// Countries returns a copy of the loaded country data.
func (d *CountryData) Countries() []Country {
	return cloneCountries(d.countries)
}

// CodesISO2 returns ISO 3166 alpha-2 codes in data order.
func (d *CountryData) CodesISO2() []string {
	return collectCountryStrings(d.countries, func(c Country) string { return c.ISO2 })
}

// CodesISO3 returns ISO 3166 alpha-3 codes in data order.
func (d *CountryData) CodesISO3() []string {
	return collectCountryStrings(d.countries, func(c Country) string { return c.ISO3 })
}

// NamesOfficial returns official country names in data order.
func (d *CountryData) NamesOfficial() []string {
	return collectCountryStrings(d.countries, func(c Country) string { return c.NameOfficial })
}

// NamesShort returns short country names in data order.
func (d *CountryData) NamesShort() []string {
	return collectCountryStrings(d.countries, func(c Country) string { return c.NameShort })
}

// Get returns the matching Country and true, or the supplied default/false when
// no match is found. Lookup order follows countryguess: ISO2, ISO3, custom
// regex, packaged regex, then fuzzy official/short names.
func (d *CountryData) Get(country string, options ...GuessOptions) (Country, bool) {
	var opts GuessOptions
	if len(options) > 0 {
		opts = options[0]
	}

	matched, ok := d.findCountry(country, opts.RegexMap)
	if ok {
		return matched, true
	}
	if opts.Default != nil {
		return *opts.Default, false
	}
	return Country{}, false
}

// MustGet is the indexing-style helper: it returns an error when no country is found.
func (d *CountryData) MustGet(country string, options ...GuessOptions) (Country, error) {
	if c, ok := d.Get(country, options...); ok {
		return c, nil
	}
	return Country{}, fmt.Errorf("country not matched: %s", country)
}

// Attribute returns one country attribute by its JSON key, for example
// "iso2", "name_official", or "continent".
func (d *CountryData) Attribute(country, attribute string, options ...GuessOptions) (string, error) {
	c, err := d.MustGet(country, options...)
	if err != nil {
		return "", err
	}
	return countryAttribute(c, attribute)
}

func (d *CountryData) NameOfficial(country string) (string, error) {
	return d.Attribute(country, "name_official")
}

func (d *CountryData) NameShort(country string) (string, error) {
	return d.Attribute(country, "name_short")
}

func (d *CountryData) ContinentName(country string) (string, error) {
	return d.Attribute(country, "continent")
}

func (d *CountryData) ISO2Code(country string) (string, error) {
	return d.Attribute(country, "iso2")
}

func (d *CountryData) ISO3Code(country string) (string, error) {
	return d.Attribute(country, "iso3")
}

func (d *CountryData) findCountry(country string, regexMap map[string]*regexp.Regexp) (Country, bool) {
	country = strings.TrimSpace(country)
	if country == "" {
		return Country{}, false
	}

	if utf8.RuneCountInString(country) == 2 {
		if c, ok := d.findCountryByCode(country, func(c Country) string { return c.ISO2 }); ok {
			return c, true
		}
	}
	if utf8.RuneCountInString(country) == 3 {
		if c, ok := d.findCountryByCode(country, func(c Country) string { return c.ISO3 }); ok {
			return c, true
		}
	}
	if utf8.RuneCountInString(country) == 4 {
		if c, ok := d.findCountryByCode(country, func(c Country) string { return c.ISO4 }); ok {
			return c, true
		}
	}

	if len(regexMap) > 0 {
		for iso2, re := range regexMap {
			if re != nil && re.MatchString(country) {
				return d.findCountryByCode(iso2, func(c Country) string { return c.ISO2 })
			}
		}
		if err := d.validateRegexMap(regexMap); err != nil {
			return Country{}, false
		}
	}

	if c, ok := d.findExactName(country); ok {
		return c, true
	}

	d.compileRegexes()
	for _, item := range d.regexes {
		if item.re.MatchString(country) {
			return item.country, true
		}
	}

	if c, ok := d.findCloseMatch(country, func(c Country) string { return c.NameOfficial }); ok {
		return c, true
	}
	return d.findCloseMatch(country, func(c Country) string { return c.NameShort })
}

func (d *CountryData) compileRegexes() {
	d.once.Do(func() {
		for _, country := range d.countries {
			if country.Regex == "" {
				continue
			}
			re, err := regexp.Compile("(?i)" + country.Regex)
			if err != nil {
				re, err = regexp.Compile("(?i)" + simplifyPythonRegex(country.Regex))
				if err != nil {
					continue
				}
			}
			d.regexes = append(d.regexes, compiledRegex{country: country, re: re})
		}
	})
}

func simplifyPythonRegex(pattern string) string {
	replacements := []*regexp.Regexp{
		regexp.MustCompile(`\(\?=[^)]*\)`),
		regexp.MustCompile(`\(\?![^)]*\)`),
	}
	for _, re := range replacements {
		pattern = re.ReplaceAllString(pattern, "")
	}
	return pattern
}

func (d *CountryData) validateRegexMap(regexMap map[string]*regexp.Regexp) error {
	for iso2, re := range regexMap {
		if re == nil {
			return fmt.Errorf("not a regular expression: %s", iso2)
		}
		if _, ok := d.findCountryByCode(iso2, func(c Country) string { return c.ISO2 }); !ok {
			return fmt.Errorf("not an ISO 3166-1 alpha-2 country code: %s", iso2)
		}
	}
	return nil
}

func (d *CountryData) findCountryByCode(code string, field func(Country) string) (Country, bool) {
	code = strings.ToUpper(strings.TrimSpace(code))
	for _, c := range d.countries {
		if strings.EqualFold(field(c), code) {
			return c, true
		}
	}
	return Country{}, false
}

func (d *CountryData) findExactName(country string) (Country, bool) {
	for _, c := range d.countries {
		if strings.EqualFold(country, c.NameShort) || strings.EqualFold(country, c.NameOfficial) ||
			strings.EqualFold(country, c.Name) || strings.EqualFold(country, c.FormalName) {
			return c, true
		}
	}
	return Country{}, false
}

func (d *CountryData) findCloseMatch(country string, field func(Country) string) (Country, bool) {
	bestRatio := fuzzyCutoff
	bestIndex := -1
	for i, c := range d.countries {
		name := field(c)
		if name == "" {
			continue
		}
		if ratio := sequenceRatio(strings.ToLower(country), strings.ToLower(name)); ratio >= bestRatio {
			bestRatio = ratio
			bestIndex = i
		}
	}
	if bestIndex == -1 {
		return Country{}, false
	}
	return d.countries[bestIndex], true
}

// Iso2 returns the country matching an ISO 3166 alpha-2 code.
func Iso2(code string) (Country, error) {
	if c, ok := defaultCountryData.Get(code); ok && strings.EqualFold(c.ISO2, strings.TrimSpace(code)) {
		return c, nil
	}
	return Country{}, errors.New("ISO2 not matched")
}

// Iso3 returns the country matching an ISO 3166 alpha-3 code.
func Iso3(code string) (Country, error) {
	if c, ok := defaultCountryData.Get(code); ok && strings.EqualFold(c.ISO3, strings.TrimSpace(code)) {
		return c, nil
	}
	return Country{}, errors.New("ISO3 not matched")
}

// Iso4 returns the country matching an ISO 3166-3 alpha-4 historical code.
func Iso4(code string) (Country, error) {
	if c, ok := defaultCountryData.findCountryByCode(code, func(c Country) string { return c.ISO4 }); ok {
		return c, nil
	}
	return Country{}, errors.New("ISO4 not matched")
}

// GuessCountry returns the best matching country, or an error when no match exists.
func GuessCountry(name string) (Country, error) {
	if c, ok := defaultCountryData.Get(name); ok {
		return c, nil
	}
	return Country{}, errors.New("country not matched")
}

// GuessCountryWithOptions exposes countryguess-like defaults, attribute lookup,
// and custom regex matching through one Go-friendly call.
func GuessCountryWithOptions(name string, options GuessOptions) (Country, error) {
	if c, ok := defaultCountryData.Get(name, options); ok || options.Default != nil {
		return c, nil
	}
	return Country{}, errors.New("country not matched")
}

// GuessCountryAttribute returns a single attribute from the matched country.
func GuessCountryAttribute(name, attribute string) (string, error) {
	return defaultCountryData.Attribute(name, attribute)
}

func makeISO2LUT() (map[string]string, map[string]string) {
	iso2lut := make(map[string]string, len(goguesscountries))
	countrylut := make(map[string]string, len(goguesscountries))
	for _, c := range goguesscountries {
		if c.ISO2 == "" {
			continue
		}
		iso2 := strings.ToLower(c.ISO2)
		country := c.Name
		iso2lut[iso2] = country
		countrylut[country] = iso2
	}
	return countrylut, iso2lut
}

func countryAttribute(country Country, attribute string) (string, error) {
	attribute = strings.ToLower(attribute)
	value := reflect.ValueOf(country)
	typeOfCountry := value.Type()
	for i := 0; i < typeOfCountry.NumField(); i++ {
		field := typeOfCountry.Field(i)
		jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonName == attribute || strings.EqualFold(field.Name, attribute) {
			fieldValue := value.Field(i)
			switch fieldValue.Kind() {
			case reflect.String:
				return fieldValue.String(), nil
			case reflect.Bool:
				if fieldValue.Bool() {
					return "true", nil
				}
				return "false", nil
			case reflect.Int, reflect.Int64, reflect.Int32:
				return fmt.Sprint(fieldValue.Int()), nil
			case reflect.Float64, reflect.Float32:
				return fmt.Sprint(fieldValue.Float()), nil
			}
		}
	}
	return "", fmt.Errorf("unknown country attribute: %s", attribute)
}

func sequenceRatio(a, b string) float64 {
	ar := []rune(a)
	br := []rune(b)
	if len(ar) == 0 && len(br) == 0 {
		return 1
	}
	matches := longestCommonSubsequence(ar, br)
	return (2.0 * float64(matches)) / float64(len(ar)+len(br))
}

func longestCommonSubsequence(a, b []rune) int {
	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				current[j] = previous[j-1] + 1
			} else if previous[j] > current[j-1] {
				current[j] = previous[j]
			} else {
				current[j] = current[j-1]
			}
		}
		previous, current = current, previous
		for j := range current {
			current[j] = 0
		}
	}
	return previous[len(b)]
}

func collectCountryStrings(countries []Country, field func(Country) string) []string {
	out := make([]string, 0, len(countries))
	for _, c := range countries {
		out = append(out, field(c))
	}
	return out
}

func cloneCountries(countries []Country) []Country {
	out := make([]Country, len(countries))
	copy(out, countries)
	return out
}

func unmarshalCountries(data []byte) ([]Country, error) {
	var countries []Country
	if err := json.Unmarshal(data, &countries); err != nil {
		return nil, err
	}
	normalizeCountries(countries)
	return countries, nil
}

func normalizeCountries(countries []Country) {
	for i := range countries {
		c := &countries[i]
		if c.NameShort == "" {
			c.NameShort = c.Name
		}
		if c.NameOfficial == "" {
			c.NameOfficial = c.FormalName
		}
		if c.Name == "" {
			c.Name = c.NameShort
		}
		if c.FormalName == "" {
			c.FormalName = c.NameOfficial
		}
		if c.RegionUN == "" {
			c.RegionUN = c.UNRegion
		}
	}
}

func mustLoadCountries(raw string) []Country {
	countries, err := unmarshalCountries([]byte(raw))
	if err != nil {
		panic(err)
	}
	return countries
}
