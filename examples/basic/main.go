package main

import (
	"fmt"
	"log"
	"regexp"

	guesscountry "go-guesscountry"
)

func main() {
	country, err := guesscountry.GuessCountry("britain")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("britain -> %s (%s/%s)\n", country.NameShort, country.ISO2, country.ISO3)

	iso2, err := guesscountry.GuessCountryAttribute("PoRtUgAl", "iso2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PoRtUgAl iso2 -> %s\n", iso2)

	historical, err := guesscountry.Iso4("SUHH")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SUHH -> %s, withdrawn %s\n", historical.NameShort, historical.WithdrawalDate)

	countries := guesscountry.NewCountryData()
	custom, ok := countries.Get("Mongol Uls", guesscountry.GuessOptions{
		RegexMap: map[string]*regexp.Regexp{
			"MN": regexp.MustCompile(`(?i)^mongol\s+uls$`),
		},
	})
	if !ok {
		log.Fatal("custom regex did not match")
	}
	fmt.Printf("Mongol Uls -> %s\n", custom.NameShort)

	if _, ok := countries.Get("no such country"); !ok {
		fmt.Println("no such country -> no match")
	}
}
