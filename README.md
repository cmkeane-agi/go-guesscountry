# go-guesscountry
A Go implementation that mimics the Python `countryguess` package.

```go
country, err := goguessocuntry.GuessCountry("britain")
iso2, err := goguessocuntry.GuessCountryAttribute("PoRtUgAl", "iso2")

countries := goguessocuntry.NewCountryData()
officialName, err := countries.NameOfficial("TW")
```

Lookups support ISO 3166-1 alpha-2, ISO 3166-1 alpha-3, ISO 3166-3 alpha-4
historical codes, packaged regex aliases, custom regex maps, and fuzzy matching
against official and short names.

Run the example program:

```sh
go run ./examples/basic
```
