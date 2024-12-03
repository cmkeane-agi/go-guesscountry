package goguessocuntry

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// CountryList represents a list of countries
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

var goguesscountries = []Country{
	{Name: "Costa Rica", FormalName: "Republic of Costa Rica", Abbrev: "C.R.", NameAlt: "", Population: 5047561, GDP: 61801, ISO2: "CR", ISO3: "CRI", Continent: "North America", RegionUN: "Americas", Subregion: "Central America", RegionWB: "Latin America & Caribbean", LabelX: -84.077922, LabelY: 10.0651},
	{Name: "Nicaragua", FormalName: "Republic of Nicaragua", Abbrev: "Nic.", NameAlt: "", Population: 6545502, GDP: 12520, ISO2: "NI", ISO3: "NIC", Continent: "North America", RegionUN: "Americas", Subregion: "Central America", RegionWB: "Latin America & Caribbean", LabelX: -85.069347, LabelY: 12.670697},
	{Name: "Haiti", FormalName: "Republic of Haiti", Abbrev: "Haiti", NameAlt: "", Population: 11263077, GDP: 14332, ISO2: "HT", ISO3: "HTI", Continent: "North America", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -72.224051, LabelY: 19.263784},
	{Name: "Dominican Rep.", FormalName: "Dominican Republic", Abbrev: "Dom. Rep.", NameAlt: "", Population: 10738958, GDP: 88941, ISO2: "DO", ISO3: "DOM", Continent: "North America", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -70.653998, LabelY: 19.104137},
	{Name: "El Salvador", FormalName: "Republic of El Salvador", Abbrev: "El. S.", NameAlt: "", Population: 6453553, GDP: 27022, ISO2: "SV", ISO3: "SLV", Continent: "North America", RegionUN: "Americas", Subregion: "Central America", RegionWB: "Latin America & Caribbean", LabelX: -88.890124, LabelY: 13.685371},
	{Name: "Guatemala", FormalName: "Republic of Guatemala", Abbrev: "Guat.", NameAlt: "", Population: 16604026, GDP: 76710, ISO2: "GT", ISO3: "GTM", Continent: "North America", RegionUN: "Americas", Subregion: "Central America", RegionWB: "Latin America & Caribbean", LabelX: -90.497134, LabelY: 14.982133},
	{Name: "Cuba", FormalName: "Republic of Cuba", Abbrev: "Cuba", NameAlt: "", Population: 11333483, GDP: 100023, ISO2: "CU", ISO3: "CUB", Continent: "North America", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -77.975855, LabelY: 21.334024},
	{Name: "Honduras", FormalName: "Republic of Honduras", Abbrev: "Hond.", NameAlt: "", Population: 9746117, GDP: 25095, ISO2: "HN", ISO3: "HND", Continent: "North America", RegionUN: "Americas", Subregion: "Central America", RegionWB: "Latin America & Caribbean", LabelX: -86.887604, LabelY: 14.794801},
	{Name: "United States", FormalName: "United States of America", Abbrev: "U.S.A.", NameAlt: "", Population: 328239523, GDP: 21433226, ISO2: "US", ISO3: "USA", Continent: "North America", RegionUN: "Americas", Subregion: "Northern America", RegionWB: "North America", LabelX: -97.482602, LabelY: 39.538479},
	{Name: "Canada", FormalName: "Canada", Abbrev: "Can.", NameAlt: "", Population: 37589262, GDP: 1736425, ISO2: "CA", ISO3: "CAN", Continent: "North America", RegionUN: "Americas", Subregion: "Northern America", RegionWB: "North America", LabelX: -101.9107, LabelY: 60.324287},
	{Name: "Mexico", FormalName: "United Mexican States", Abbrev: "Mex.", NameAlt: "", Population: 127575529, GDP: 1268870, ISO2: "MX", ISO3: "MEX", Continent: "North America", RegionUN: "Americas", Subregion: "Central America", RegionWB: "Latin America & Caribbean", LabelX: -102.289448, LabelY: 23.919988},
	{Name: "Belize", FormalName: "Belize", Abbrev: "Belize", NameAlt: "", Population: 390353, GDP: 1879, ISO2: "BZ", ISO3: "BLZ", Continent: "North America", RegionUN: "Americas", Subregion: "Central America", RegionWB: "Latin America & Caribbean", LabelX: -88.712962, LabelY: 17.202068},
	{Name: "Panama", FormalName: "Republic of Panama", Abbrev: "Pan.", NameAlt: "", Population: 4246439, GDP: 66800, ISO2: "PA", ISO3: "PAN", Continent: "North America", RegionUN: "Americas", Subregion: "Central America", RegionWB: "Latin America & Caribbean", LabelX: -80.352106, LabelY: 8.72198},
	{Name: "Greenland", FormalName: "Greenland", Abbrev: "Grlnd.", NameAlt: "", Population: 56225, GDP: 3051, ISO2: "GL", ISO3: "GRL", Continent: "North America", RegionUN: "Americas", Subregion: "Northern America", RegionWB: "Europe & Central Asia", LabelX: -39.335251, LabelY: 74.319387},
	{Name: "Bahamas", FormalName: "Commonwealth of the Bahamas", Abbrev: "Bhs.", NameAlt: "", Population: 389482, GDP: 13578, ISO2: "BS", ISO3: "BHS", Continent: "North America", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -77.146688, LabelY: 26.401789},
	{Name: "Trinidad and Tobago", FormalName: "Republic of Trinidad and Tobago", Abbrev: "Tr.T.", NameAlt: "", Population: 1394973, GDP: 24269, ISO2: "TT", ISO3: "TTO", Continent: "North America", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -60.9184, LabelY: 10.9989},
	{Name: "Puerto Rico", FormalName: "Commonwealth of Puerto Rico", Abbrev: "P.R.", NameAlt: "", Population: 3193694, GDP: 104988, ISO2: "PR", ISO3: "PRI", Continent: "North America", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -66.481065, LabelY: 18.234668},
	{Name: "Jamaica", FormalName: "Jamaica", Abbrev: "Jam.", NameAlt: "", Population: 2948279, GDP: 16458, ISO2: "JM", ISO3: "JAM", Continent: "North America", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -77.318767, LabelY: 18.137124},
	{Name: "Indonesia", FormalName: "Republic of Indonesia", Abbrev: "Indo.", NameAlt: "", Population: 270625568, GDP: 1119190, ISO2: "ID", ISO3: "IDN", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 101.892949, LabelY: -0.954404},
	{Name: "Malaysia", FormalName: "Malaysia", Abbrev: "Malay.", NameAlt: "", Population: 31949777, GDP: 364681, ISO2: "MY", ISO3: "MYS", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 113.83708, LabelY: 2.528667},
	{Name: "Cyprus", FormalName: "Republic of Cyprus", Abbrev: "Cyp.", NameAlt: "", Population: 1198575, GDP: 24948, ISO2: "CY", ISO3: "CYP", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Europe & Central Asia", LabelX: 33.084182, LabelY: 34.913329},
	{Name: "India", FormalName: "Republic of India", Abbrev: "India", NameAlt: "", Population: 1366417754, GDP: 2868929, ISO2: "IN", ISO3: "IND", Continent: "Asia", RegionUN: "Asia", Subregion: "Southern Asia", RegionWB: "South Asia", LabelX: 79.358105, LabelY: 22.686852},
	{Name: "China", FormalName: "People's Republic of China", Abbrev: "China", NameAlt: "", Population: 1397715000, GDP: 14342903, ISO2: "CN", ISO3: "CHN", Continent: "Asia", RegionUN: "Asia", Subregion: "Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 106.337289, LabelY: 32.498178},
	{Name: "Israel", FormalName: "State of Israel", Abbrev: "Isr.", NameAlt: "", Population: 9053300, GDP: 394652, ISO2: "IL", ISO3: "ISR", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 34.847915, LabelY: 30.911148},
	{Name: "Palestine", FormalName: "West Bank and Gaza", Abbrev: "Pal.", NameAlt: "", Population: 4685306, GDP: 16276, ISO2: "PS", ISO3: "PSE", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 35.291341, LabelY: 32.047431},
	{Name: "Lebanon", FormalName: "Lebanese Republic", Abbrev: "Leb.", NameAlt: "", Population: 6855713, GDP: 51991, ISO2: "LB", ISO3: "LBN", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 35.992892, LabelY: 34.133368},
	{Name: "Syria", FormalName: "Syrian Arab Republic", Abbrev: "Syria", NameAlt: "", Population: 17070135, GDP: 98830, ISO2: "SY", ISO3: "SYR", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 38.277783, LabelY: 35.006636},
	{Name: "South Korea", FormalName: "Republic of Korea", Abbrev: "S.K.", NameAlt: "", Population: 51709098, GDP: 1646739, ISO2: "KR", ISO3: "KOR", Continent: "Asia", RegionUN: "Asia", Subregion: "Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 128.129504, LabelY: 36.384924},
	{Name: "North Korea", FormalName: "Democratic People's Republic of Korea", Abbrev: "N.K.", NameAlt: "", Population: 25666161, GDP: 40000, ISO2: "KP", ISO3: "PRK", Continent: "Asia", RegionUN: "Asia", Subregion: "Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 126.444516, LabelY: 39.885252},
	{Name: "Bhutan", FormalName: "Kingdom of Bhutan", Abbrev: "Bhutan", NameAlt: "", Population: 763092, GDP: 2530, ISO2: "BT", ISO3: "BTN", Continent: "Asia", RegionUN: "Asia", Subregion: "Southern Asia", RegionWB: "South Asia", LabelX: 90.040294, LabelY: 27.536685},
	{Name: "Oman", FormalName: "Sultanate of Oman", Abbrev: "Oman", NameAlt: "", Population: 4974986, GDP: 76331, ISO2: "OM", ISO3: "OMN", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 57.336553, LabelY: 22.120427},
	{Name: "Uzbekistan", FormalName: "Republic of Uzbekistan", Abbrev: "Uzb.", NameAlt: "", Population: 33580650, GDP: 57921, ISO2: "UZ", ISO3: "UZB", Continent: "Asia", RegionUN: "Asia", Subregion: "Central Asia", RegionWB: "Europe & Central Asia", LabelX: 64.005429, LabelY: 41.693603},
	{Name: "Kazakhstan", FormalName: "Republic of Kazakhstan", Abbrev: "Kaz.", NameAlt: "", Population: 18513930, GDP: 181665, ISO2: "KZ", ISO3: "KAZ", Continent: "Asia", RegionUN: "Asia", Subregion: "Central Asia", RegionWB: "Europe & Central Asia", LabelX: 68.685548, LabelY: 49.054149},
	{Name: "Tajikistan", FormalName: "Republic of Tajikistan", Abbrev: "Tjk.", NameAlt: "", Population: 9321018, GDP: 8116, ISO2: "TJ", ISO3: "TJK", Continent: "Asia", RegionUN: "Asia", Subregion: "Central Asia", RegionWB: "Europe & Central Asia", LabelX: 72.587276, LabelY: 38.199835},
	{Name: "Mongolia", FormalName: "Mongolia", Abbrev: "Mong.", NameAlt: "", Population: 3225167, GDP: 13996, ISO2: "MN", ISO3: "MNG", Continent: "Asia", RegionUN: "Asia", Subregion: "Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 104.150405, LabelY: 45.997488},
	{Name: "Vietnam", FormalName: "Socialist Republic of Vietnam", Abbrev: "Viet.", NameAlt: "", Population: 96462106, GDP: 261921, ISO2: "VN", ISO3: "VNM", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 105.387292, LabelY: 21.715416},
	{Name: "Cambodia", FormalName: "Kingdom of Cambodia", Abbrev: "Camb.", NameAlt: "", Population: 16486542, GDP: 27089, ISO2: "KH", ISO3: "KHM", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 104.50487, LabelY: 12.647584},
	{Name: "United Arab Emirates", FormalName: "United Arab Emirates", Abbrev: "U.A.E.", NameAlt: "", Population: 9770529, GDP: 421142, ISO2: "AE", ISO3: "ARE", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 54.547256, LabelY: 23.466285},
	{Name: "Georgia", FormalName: "Georgia", Abbrev: "Geo.", NameAlt: "", Population: 3720382, GDP: 17477, ISO2: "GE", ISO3: "GEO", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Europe & Central Asia", LabelX: 43.735724, LabelY: 41.870087},
	{Name: "Azerbaijan", FormalName: "Republic of Azerbaijan", Abbrev: "Aze.", NameAlt: "", Population: 10023318, GDP: 48047, ISO2: "AZ", ISO3: "AZE", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Europe & Central Asia", LabelX: 47.210994, LabelY: 40.402387},
	{Name: "Turkey", FormalName: "Republic of Turkey", Abbrev: "Tur.", NameAlt: "", Population: 83429615, GDP: 761425, ISO2: "TR", ISO3: "TUR", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Europe & Central Asia", LabelX: 34.508268, LabelY: 39.345388},
	{Name: "Laos", FormalName: "Lao People's Democratic Republic", Abbrev: "Laos", NameAlt: "", Population: 7169455, GDP: 18173, ISO2: "LA", ISO3: "LAO", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 102.533912, LabelY: 19.431821},
	{Name: "Kyrgyzstan", FormalName: "Kyrgyz Republic", Abbrev: "Kgz.", NameAlt: "", Population: 6456900, GDP: 8454, ISO2: "KG", ISO3: "KGZ", Continent: "Asia", RegionUN: "Asia", Subregion: "Central Asia", RegionWB: "Europe & Central Asia", LabelX: 74.532637, LabelY: 41.66854},
	{Name: "Armenia", FormalName: "Republic of Armenia", Abbrev: "Arm.", NameAlt: "", Population: 2957731, GDP: 13672, ISO2: "AM", ISO3: "ARM", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Europe & Central Asia", LabelX: 44.800564, LabelY: 40.459077},
	{Name: "Iraq", FormalName: "Republic of Iraq", Abbrev: "Iraq", NameAlt: "", Population: 39309783, GDP: 234094, ISO2: "IQ", ISO3: "IRQ", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 43.26181, LabelY: 33.09403},
	{Name: "Iran", FormalName: "Islamic Republic of Iran", Abbrev: "Iran", NameAlt: "", Population: 82913906, GDP: 453996, ISO2: "IR", ISO3: "IRN", Continent: "Asia", RegionUN: "Asia", Subregion: "Southern Asia", RegionWB: "Middle East & North Africa", LabelX: 54.931495, LabelY: 32.166225},
	{Name: "Qatar", FormalName: "State of Qatar", Abbrev: "Qatar", NameAlt: "", Population: 2832067, GDP: 175837, ISO2: "QA", ISO3: "QAT", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 51.143509, LabelY: 25.237383},
	{Name: "Saudi Arabia", FormalName: "Kingdom of Saudi Arabia", Abbrev: "Saud.", NameAlt: "", Population: 34268528, GDP: 792966, ISO2: "SA", ISO3: "SAU", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 44.6996, LabelY: 23.806908},
	{Name: "Pakistan", FormalName: "Islamic Republic of Pakistan", Abbrev: "Pak.", NameAlt: "", Population: 216565318, GDP: 278221, ISO2: "PK", ISO3: "PAK", Continent: "Asia", RegionUN: "Asia", Subregion: "Southern Asia", RegionWB: "South Asia", LabelX: 68.545632, LabelY: 29.328389},
	{Name: "Thailand", FormalName: "Kingdom of Thailand", Abbrev: "Thai.", NameAlt: "", Population: 69625582, GDP: 543548, ISO2: "TH", ISO3: "THA", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 101.073198, LabelY: 15.45974},
	{Name: "Kuwait", FormalName: "State of Kuwait", Abbrev: "Kwt.", NameAlt: "", Population: 4207083, GDP: 134628, ISO2: "KW", ISO3: "KWT", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 47.313999, LabelY: 29.413628},
	{Name: "Timor-Leste", FormalName: "Democratic Republic of Timor-Leste", Abbrev: "T.L.", NameAlt: "East Timor", Population: 1293119, GDP: 2017, ISO2: "TL", ISO3: "TLS", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 125.854679, LabelY: -8.803705},
	{Name: "Brunei", FormalName: "Negara Brunei Darussalam", Abbrev: "Brunei", NameAlt: "", Population: 433285, GDP: 13469, ISO2: "BN", ISO3: "BRN", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 114.551943, LabelY: 4.448298},
	{Name: "Myanmar", FormalName: "Republic of the Union of Myanmar", Abbrev: "Myan.", NameAlt: "", Population: 54045420, GDP: 76085, ISO2: "MM", ISO3: "MMR", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 95.804497, LabelY: 21.573855},
	{Name: "Bangladesh", FormalName: "People's Republic of Bangladesh", Abbrev: "Bang.", NameAlt: "", Population: 163046161, GDP: 302571, ISO2: "BD", ISO3: "BGD", Continent: "Asia", RegionUN: "Asia", Subregion: "Southern Asia", RegionWB: "South Asia", LabelX: 89.684963, LabelY: 24.214956},
	{Name: "Afghanistan", FormalName: "Islamic State of Afghanistan", Abbrev: "Afg.", NameAlt: "", Population: 38041754, GDP: 19291, ISO2: "AF", ISO3: "AFG", Continent: "Asia", RegionUN: "Asia", Subregion: "Southern Asia", RegionWB: "South Asia", LabelX: 66.496586, LabelY: 34.164262},
	{Name: "Turkmenistan", FormalName: "Turkmenistan", Abbrev: "Turkm.", NameAlt: "", Population: 5942089, GDP: 40761, ISO2: "TM", ISO3: "TKM", Continent: "Asia", RegionUN: "Asia", Subregion: "Central Asia", RegionWB: "Europe & Central Asia", LabelX: 58.676647, LabelY: 39.855246},
	{Name: "Jordan", FormalName: "Hashemite Kingdom of Jordan", Abbrev: "Jord.", NameAlt: "", Population: 10101694, GDP: 44502, ISO2: "JO", ISO3: "JOR", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 36.375991, LabelY: 30.805025},
	{Name: "Nepal", FormalName: "Nepal", Abbrev: "Nepal", NameAlt: "", Population: 28608710, GDP: 30641, ISO2: "NP", ISO3: "NPL", Continent: "Asia", RegionUN: "Asia", Subregion: "Southern Asia", RegionWB: "South Asia", LabelX: 83.639914, LabelY: 28.297925},
	{Name: "Yemen", FormalName: "Republic of Yemen", Abbrev: "Yem.", NameAlt: "", Population: 29161922, GDP: 22581, ISO2: "YE", ISO3: "YEM", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 45.874383, LabelY: 15.328226},
	{Name: "N. Cyprus", FormalName: "Turkish Republic of Northern Cyprus", Abbrev: "N. Cy.", NameAlt: "", Population: 326000, GDP: 3600, ISO2: "-99", ISO3: "-99", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Europe & Central Asia", LabelX: 33.692434, LabelY: 35.216071},
	{Name: "Philippines", FormalName: "Republic of the Philippines", Abbrev: "Phil.", NameAlt: "", Population: 108116615, GDP: 376795, ISO2: "PH", ISO3: "PHL", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 122.465, LabelY: 11.198},
	{Name: "Sri Lanka", FormalName: "Democratic Socialist Republic of Sri Lanka", Abbrev: "Sri L.", NameAlt: "", Population: 21803000, GDP: 84008, ISO2: "LK", ISO3: "LKA", Continent: "Asia", RegionUN: "Asia", Subregion: "Southern Asia", RegionWB: "South Asia", LabelX: 80.704823, LabelY: 7.581097},
	{Name: "Taiwan", FormalName: "", Abbrev: "Taiwan", NameAlt: "", Population: 23568378, GDP: 1127000, ISO2: "CN-TW", ISO3: "TWN", Continent: "Asia", RegionUN: "Asia", Subregion: "Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 120.868204, LabelY: 23.652408},
	{Name: "Japan", FormalName: "Japan", Abbrev: "Japan", NameAlt: "", Population: 126264931, GDP: 5081769, ISO2: "JP", ISO3: "JPN", Continent: "Asia", RegionUN: "Asia", Subregion: "Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 138.44217, LabelY: 36.142538},
	{Name: "Chile", FormalName: "Republic of Chile", Abbrev: "Chile", NameAlt: "", Population: 18952038, GDP: 282318, ISO2: "CL", ISO3: "CHL", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -72.318871, LabelY: -38.151771},
	{Name: "Bolivia", FormalName: "Plurinational State of Bolivia", Abbrev: "Bolivia", NameAlt: "", Population: 11513100, GDP: 40895, ISO2: "BO", ISO3: "BOL", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -64.593433, LabelY: -16.666015},
	{Name: "Peru", FormalName: "Republic of Peru", Abbrev: "Peru", NameAlt: "", Population: 32510453, GDP: 226848, ISO2: "PE", ISO3: "PER", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -72.90016, LabelY: -12.976679},
	{Name: "Argentina", FormalName: "Argentine Republic", Abbrev: "Arg.", NameAlt: "", Population: 44938712, GDP: 445445, ISO2: "AR", ISO3: "ARG", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -64.173331, LabelY: -33.501159},
	{Name: "Suriname", FormalName: "Republic of Suriname", Abbrev: "Sur.", NameAlt: "", Population: 581363, GDP: 3697, ISO2: "SR", ISO3: "SUR", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -55.91094, LabelY: 4.143987},
	{Name: "Guyana", FormalName: "Co-operative Republic of Guyana", Abbrev: "Guy.", NameAlt: "", Population: 782766, GDP: 5173, ISO2: "GY", ISO3: "GUY", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -58.942643, LabelY: 5.124317},
	{Name: "Brazil", FormalName: "Federative Republic of Brazil", Abbrev: "Brazil", NameAlt: "", Population: 211049527, GDP: 1839758, ISO2: "BR", ISO3: "BRA", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -49.55945, LabelY: -12.098687},
	{Name: "Uruguay", FormalName: "Oriental Republic of Uruguay", Abbrev: "Ury.", NameAlt: "", Population: 3461734, GDP: 56045, ISO2: "UY", ISO3: "URY", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -55.966942, LabelY: -32.961127},
	{Name: "Ecuador", FormalName: "Republic of Ecuador", Abbrev: "Ecu.", NameAlt: "", Population: 17373662, GDP: 107435, ISO2: "EC", ISO3: "ECU", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -78.188375, LabelY: -1.259076},
	{Name: "Colombia", FormalName: "Republic of Colombia", Abbrev: "Col.", NameAlt: "", Population: 50339443, GDP: 323615, ISO2: "CO", ISO3: "COL", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -73.174347, LabelY: 3.373111},
	{Name: "Paraguay", FormalName: "Republic of Paraguay", Abbrev: "Para.", NameAlt: "", Population: 7044636, GDP: 38145, ISO2: "PY", ISO3: "PRY", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -60.146394, LabelY: -21.674509},
	{Name: "Venezuela", FormalName: "Bolivarian Republic of Venezuela", Abbrev: "Ven.", NameAlt: "", Population: 28515829, GDP: 482359, ISO2: "VE", ISO3: "VEN", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -64.599381, LabelY: 7.182476},
	{Name: "Falkland Is.", FormalName: "Falkland Islands", Abbrev: "Flk. Is.", NameAlt: "Islas Malvinas", Population: 3398, GDP: 282, ISO2: "FK", ISO3: "FLK", Continent: "South America", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -58.738602, LabelY: -51.608913},
	{Name: "Ethiopia", FormalName: "Federal Democratic Republic of Ethiopia", Abbrev: "Eth.", NameAlt: "", Population: 112078730, GDP: 95912, ISO2: "ET", ISO3: "ETH", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 39.0886, LabelY: 8.032795},
	{Name: "S. Sudan", FormalName: "Republic of South Sudan", Abbrev: "S. Sud.", NameAlt: "", Population: 11062113, GDP: 11998, ISO2: "SS", ISO3: "SSD", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 30.390151, LabelY: 7.230477},
	{Name: "Somalia", FormalName: "Federal Republic of Somalia", Abbrev: "Som.", NameAlt: "", Population: 10192317, GDP: 4719, ISO2: "SO", ISO3: "SOM", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 45.19238, LabelY: 3.568925},
	{Name: "Kenya", FormalName: "Republic of Kenya", Abbrev: "Ken.", NameAlt: "", Population: 52573973, GDP: 95503, ISO2: "KE", ISO3: "KEN", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 37.907632, LabelY: 0.549043},
	{Name: "Malawi", FormalName: "Republic of Malawi", Abbrev: "Mal.", NameAlt: "", Population: 18628747, GDP: 7666, ISO2: "MW", ISO3: "MWI", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 33.608082, LabelY: -13.386737},
	{Name: "Tanzania", FormalName: "United Republic of Tanzania", Abbrev: "Tanz.", NameAlt: "", Population: 58005463, GDP: 63177, ISO2: "TZ", ISO3: "TZA", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 34.959183, LabelY: -6.051866},
	{Name: "Somaliland", FormalName: "Republic of Somaliland", Abbrev: "Solnd.", NameAlt: "", Population: 5096159, GDP: 17836, ISO2: "-99", ISO3: "-99", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 46.731595, LabelY: 9.443889},
	{Name: "Morocco", FormalName: "Kingdom of Morocco", Abbrev: "Mor.", NameAlt: "", Population: 36471769, GDP: 119700, ISO2: "MA", ISO3: "MAR", Continent: "Africa", RegionUN: "Africa", Subregion: "Northern Africa", RegionWB: "Middle East & North Africa", LabelX: -7.187296, LabelY: 31.650723},
	{Name: "W. Sahara", FormalName: "Sahrawi Arab Democratic Republic", Abbrev: "W. Sah.", NameAlt: "", Population: 603253, GDP: 907, ISO2: "EH", ISO3: "ESH", Continent: "Africa", RegionUN: "Africa", Subregion: "Northern Africa", RegionWB: "Middle East & North Africa", LabelX: -12.630304, LabelY: 23.967592},
	{Name: "Congo", FormalName: "Republic of the Congo", Abbrev: "Rep. Congo", NameAlt: "", Population: 5380508, GDP: 12267, ISO2: "CG", ISO3: "COG", Continent: "Africa", RegionUN: "Africa", Subregion: "Middle Africa", RegionWB: "Sub-Saharan Africa", LabelX: 15.9005, LabelY: 0.142331},
	{Name: "Dem. Rep. Congo", FormalName: "Democratic Republic of the Congo", Abbrev: "D.R.C.", NameAlt: "", Population: 86790567, GDP: 50400, ISO2: "CD", ISO3: "COD", Continent: "Africa", RegionUN: "Africa", Subregion: "Middle Africa", RegionWB: "Sub-Saharan Africa", LabelX: 23.458829, LabelY: -1.858167},
	{Name: "Namibia", FormalName: "Republic of Namibia", Abbrev: "Nam.", NameAlt: "", Population: 2494530, GDP: 12366, ISO2: "NA", ISO3: "NAM", Continent: "Africa", RegionUN: "Africa", Subregion: "Southern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 17.108166, LabelY: -20.575298},
	{Name: "South Africa", FormalName: "Republic of South Africa", Abbrev: "S.Af.", NameAlt: "", Population: 58558270, GDP: 351431, ISO2: "ZA", ISO3: "ZAF", Continent: "Africa", RegionUN: "Africa", Subregion: "Southern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 23.665734, LabelY: -29.708776},
	{Name: "Libya", FormalName: "Libya", Abbrev: "Libya", NameAlt: "", Population: 6777452, GDP: 52091, ISO2: "LY", ISO3: "LBY", Continent: "Africa", RegionUN: "Africa", Subregion: "Northern Africa", RegionWB: "Middle East & North Africa", LabelX: 18.011015, LabelY: 26.638944},
	{Name: "Tunisia", FormalName: "Republic of Tunisia", Abbrev: "Tun.", NameAlt: "", Population: 11694719, GDP: 38796, ISO2: "TN", ISO3: "TUN", Continent: "Africa", RegionUN: "Africa", Subregion: "Northern Africa", RegionWB: "Middle East & North Africa", LabelX: 9.007881, LabelY: 33.687263},
	{Name: "Zambia", FormalName: "Republic of Zambia", Abbrev: "Zambia", NameAlt: "", Population: 17861030, GDP: 23309, ISO2: "ZM", ISO3: "ZMB", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 26.395298, LabelY: -14.660804},
	{Name: "Sierra Leone", FormalName: "Republic of Sierra Leone", Abbrev: "S.L.", NameAlt: "", Population: 7813215, GDP: 4121, ISO2: "SL", ISO3: "SLE", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -11.763677, LabelY: 8.617449},
	{Name: "Guinea", FormalName: "Republic of Guinea", Abbrev: "Gin.", NameAlt: "", Population: 12771246, GDP: 12296, ISO2: "GN", ISO3: "GIN", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -10.016402, LabelY: 10.618516},
	{Name: "Liberia", FormalName: "Republic of Liberia", Abbrev: "Liberia", NameAlt: "", Population: 4937374, GDP: 3070, ISO2: "LR", ISO3: "LBR", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -9.460379, LabelY: 6.447177},
	{Name: "Central African Rep.", FormalName: "Central African Republic", Abbrev: "C.A.R.", NameAlt: "", Population: 4745185, GDP: 2220, ISO2: "CF", ISO3: "CAF", Continent: "Africa", RegionUN: "Africa", Subregion: "Middle Africa", RegionWB: "Sub-Saharan Africa", LabelX: 20.906897, LabelY: 6.989681},
	{Name: "Sudan", FormalName: "Republic of the Sudan", Abbrev: "Sudan", NameAlt: "", Population: 42813238, GDP: 30513, ISO2: "SD", ISO3: "SDN", Continent: "Africa", RegionUN: "Africa", Subregion: "Northern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 29.260657, LabelY: 16.330746},
	{Name: "Djibouti", FormalName: "Republic of Djibouti", Abbrev: "Dji.", NameAlt: "", Population: 973560, GDP: 3324, ISO2: "DJ", ISO3: "DJI", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Middle East & North Africa", LabelX: 42.498825, LabelY: 11.976343},
	{Name: "Eritrea", FormalName: "State of Eritrea", Abbrev: "Erit.", NameAlt: "", Population: 6081196, GDP: 2065, ISO2: "ER", ISO3: "ERI", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 38.285566, LabelY: 15.787401},
	{Name: "Côte d'Ivoire", FormalName: "Republic of Ivory Coast", Abbrev: "I.C.", NameAlt: "", Population: 25716544, GDP: 58539, ISO2: "CI", ISO3: "CIV", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -5.568618, LabelY: 7.49139},
	{Name: "Mali", FormalName: "Republic of Mali", Abbrev: "Mali", NameAlt: "", Population: 19658031, GDP: 17279, ISO2: "ML", ISO3: "MLI", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -2.038455, LabelY: 18.692713},
	{Name: "Senegal", FormalName: "Republic of Senegal", Abbrev: "Sen.", NameAlt: "", Population: 16296364, GDP: 23578, ISO2: "SN", ISO3: "SEN", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -14.778586, LabelY: 15.138125},
	{Name: "Nigeria", FormalName: "Federal Republic of Nigeria", Abbrev: "Nigeria", NameAlt: "", Population: 200963599, GDP: 448120, ISO2: "NG", ISO3: "NGA", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: 7.50322, LabelY: 9.439799},
	{Name: "Benin", FormalName: "Republic of Benin", Abbrev: "Benin", NameAlt: "", Population: 11801151, GDP: 14390, ISO2: "BJ", ISO3: "BEN", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: 2.352018, LabelY: 10.324775},
	{Name: "Angola", FormalName: "People's Republic of Angola", Abbrev: "Ang.", NameAlt: "", Population: 31825295, GDP: 88815, ISO2: "AO", ISO3: "AGO", Continent: "Africa", RegionUN: "Africa", Subregion: "Middle Africa", RegionWB: "Sub-Saharan Africa", LabelX: 17.984249, LabelY: -12.182762},
	{Name: "Botswana", FormalName: "Republic of Botswana", Abbrev: "Bwa.", NameAlt: "", Population: 2303697, GDP: 18340, ISO2: "BW", ISO3: "BWA", Continent: "Africa", RegionUN: "Africa", Subregion: "Southern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 24.179216, LabelY: -22.102634},
	{Name: "Zimbabwe", FormalName: "Republic of Zimbabwe", Abbrev: "Zimb.", NameAlt: "", Population: 14645468, GDP: 21440, ISO2: "ZW", ISO3: "ZWE", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 29.925444, LabelY: -18.91164},
	{Name: "Chad", FormalName: "Republic of Chad", Abbrev: "Chad", NameAlt: "", Population: 15946876, GDP: 11314, ISO2: "TD", ISO3: "TCD", Continent: "Africa", RegionUN: "Africa", Subregion: "Middle Africa", RegionWB: "Sub-Saharan Africa", LabelX: 18.645041, LabelY: 15.142959},
	{Name: "Algeria", FormalName: "People's Democratic Republic of Algeria", Abbrev: "Alg.", NameAlt: "", Population: 43053054, GDP: 171091, ISO2: "DZ", ISO3: "DZA", Continent: "Africa", RegionUN: "Africa", Subregion: "Northern Africa", RegionWB: "Middle East & North Africa", LabelX: 2.808241, LabelY: 27.397406},
	{Name: "Mozambique", FormalName: "Republic of Mozambique", Abbrev: "Moz.", NameAlt: "", Population: 30366036, GDP: 15291, ISO2: "MZ", ISO3: "MOZ", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 37.83789, LabelY: -13.94323},
	{Name: "eSwatini", FormalName: "Kingdom of eSwatini", Abbrev: "eSw.", NameAlt: "Swaziland", Population: 1148130, GDP: 4471, ISO2: "SZ", ISO3: "SWZ", Continent: "Africa", RegionUN: "Africa", Subregion: "Southern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 31.467264, LabelY: -26.533676},
	{Name: "Burundi", FormalName: "Republic of Burundi", Abbrev: "Bur.", NameAlt: "", Population: 11530580, GDP: 3012, ISO2: "BI", ISO3: "BDI", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 29.917086, LabelY: -3.332836},
	{Name: "Rwanda", FormalName: "Republic of Rwanda", Abbrev: "Rwa.", NameAlt: "", Population: 12626950, GDP: 10354, ISO2: "RW", ISO3: "RWA", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 30.103894, LabelY: -1.897196},
	{Name: "Uganda", FormalName: "Republic of Uganda", Abbrev: "Uga.", NameAlt: "", Population: 44269594, GDP: 35165, ISO2: "UG", ISO3: "UGA", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 32.948555, LabelY: 1.972589},
	{Name: "Lesotho", FormalName: "Kingdom of Lesotho", Abbrev: "Les.", NameAlt: "", Population: 2125268, GDP: 2376, ISO2: "LS", ISO3: "LSO", Continent: "Africa", RegionUN: "Africa", Subregion: "Southern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 28.246639, LabelY: -29.480158},
	{Name: "Cameroon", FormalName: "Republic of Cameroon", Abbrev: "Cam.", NameAlt: "", Population: 25876380, GDP: 39007, ISO2: "CM", ISO3: "CMR", Continent: "Africa", RegionUN: "Africa", Subregion: "Middle Africa", RegionWB: "Sub-Saharan Africa", LabelX: 12.473488, LabelY: 4.585041},
	{Name: "Gabon", FormalName: "Gabonese Republic", Abbrev: "Gabon", NameAlt: "", Population: 2172579, GDP: 16874, ISO2: "GA", ISO3: "GAB", Continent: "Africa", RegionUN: "Africa", Subregion: "Middle Africa", RegionWB: "Sub-Saharan Africa", LabelX: 11.835939, LabelY: -0.437739},
	{Name: "Niger", FormalName: "Republic of Niger", Abbrev: "Niger", NameAlt: "", Population: 23310715, GDP: 12911, ISO2: "NE", ISO3: "NER", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: 9.504356, LabelY: 17.446195},
	{Name: "Burkina Faso", FormalName: "Burkina Faso", Abbrev: "B.F.", NameAlt: "", Population: 20321378, GDP: 15990, ISO2: "BF", ISO3: "BFA", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -1.36388, LabelY: 12.673048},
	{Name: "Togo", FormalName: "Togolese Republic", Abbrev: "Togo", NameAlt: "", Population: 8082366, GDP: 5490, ISO2: "TG", ISO3: "TGO", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: 1.058113, LabelY: 8.80722},
	{Name: "Ghana", FormalName: "Republic of Ghana", Abbrev: "Ghana", NameAlt: "", Population: 30417856, GDP: 66983, ISO2: "GH", ISO3: "GHA", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -1.036941, LabelY: 7.717639},
	{Name: "Guinea-Bissau", FormalName: "Republic of Guinea-Bissau", Abbrev: "GnB.", NameAlt: "", Population: 1920922, GDP: 1339, ISO2: "GW", ISO3: "GNB", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -14.52413, LabelY: 12.163712},
	{Name: "Egypt", FormalName: "Arab Republic of Egypt", Abbrev: "Egypt", NameAlt: "", Population: 100388073, GDP: 303092, ISO2: "EG", ISO3: "EGY", Continent: "Africa", RegionUN: "Africa", Subregion: "Northern Africa", RegionWB: "Middle East & North Africa", LabelX: 29.445837, LabelY: 26.186173},
	{Name: "Mauritania", FormalName: "Islamic Republic of Mauritania", Abbrev: "Mrt.", NameAlt: "", Population: 4525696, GDP: 7600, ISO2: "MR", ISO3: "MRT", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -9.740299, LabelY: 19.587062},
	{Name: "Eq. Guinea", FormalName: "Republic of Equatorial Guinea", Abbrev: "Eq. G.", NameAlt: "", Population: 1355986, GDP: 11026, ISO2: "GQ", ISO3: "GNQ", Continent: "Africa", RegionUN: "Africa", Subregion: "Middle Africa", RegionWB: "Sub-Saharan Africa", LabelX: 8.9902, LabelY: 2.333},
	{Name: "Gambia", FormalName: "Republic of the Gambia", Abbrev: "Gambia", NameAlt: "", Population: 2347706, GDP: 1826, ISO2: "GM", ISO3: "GMB", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -14.998318, LabelY: 13.641721},
	{Name: "Madagascar", FormalName: "Republic of Madagascar", Abbrev: "Mad.", NameAlt: "", Population: 26969307, GDP: 14114, ISO2: "MG", ISO3: "MDG", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 46.704241, LabelY: -18.628288},
	{Name: "France", FormalName: "French Republic", Abbrev: "Fr.", NameAlt: "", Population: 67059887, GDP: 2715518, ISO2: "FR", ISO3: "FRA", Continent: "Europe", RegionUN: "Europe", Subregion: "Western Europe", RegionWB: "Europe & Central Asia", LabelX: 2.552275, LabelY: 46.696113},
	{Name: "Ukraine", FormalName: "Ukraine", Abbrev: "Ukr.", NameAlt: "", Population: 44385155, GDP: 153781, ISO2: "UA", ISO3: "UKR", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 32.140865, LabelY: 49.724739},
	{Name: "Belarus", FormalName: "Republic of Belarus", Abbrev: "Bela.", NameAlt: "", Population: 9466856, GDP: 63080, ISO2: "BY", ISO3: "BLR", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 28.417701, LabelY: 53.821888},
	{Name: "Lithuania", FormalName: "Republic of Lithuania", Abbrev: "Lith.", NameAlt: "", Population: 2786844, GDP: 54627, ISO2: "LT", ISO3: "LTU", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: 24.089932, LabelY: 55.103703},
	{Name: "Russia", FormalName: "Russian Federation", Abbrev: "Rus.", NameAlt: "", Population: 144373535, GDP: 1699876, ISO2: "RU", ISO3: "RUS", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 44.686469, LabelY: 58.249357},
	{Name: "Czechia", FormalName: "Czech Republic", Abbrev: "Cz.", NameAlt: "Česko", Population: 10669709, GDP: 250680, ISO2: "CZ", ISO3: "CZE", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 15.377555, LabelY: 49.882364},
	{Name: "Germany", FormalName: "Federal Republic of Germany", Abbrev: "Ger.", NameAlt: "", Population: 83132799, GDP: 3861123, ISO2: "DE", ISO3: "DEU", Continent: "Europe", RegionUN: "Europe", Subregion: "Western Europe", RegionWB: "Europe & Central Asia", LabelX: 9.678348, LabelY: 50.961733},
	{Name: "Estonia", FormalName: "Republic of Estonia", Abbrev: "Est.", NameAlt: "", Population: 1326590, GDP: 31471, ISO2: "EE", ISO3: "EST", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: 25.867126, LabelY: 58.724865},
	{Name: "Latvia", FormalName: "Republic of Latvia", Abbrev: "Lat.", NameAlt: "", Population: 1912789, GDP: 34102, ISO2: "LV", ISO3: "LVA", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: 25.458723, LabelY: 57.066872},
	{Name: "Norway", FormalName: "Kingdom of Norway", Abbrev: "Nor.", NameAlt: "", Population: 5347896, GDP: 403336, ISO2: "NO", ISO3: "NOR", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: 9.679975, LabelY: 61.357092},
	{Name: "Sweden", FormalName: "Kingdom of Sweden", Abbrev: "Swe.", NameAlt: "", Population: 10285453, GDP: 530883, ISO2: "SE", ISO3: "SWE", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: 19.01705, LabelY: 65.85918},
	{Name: "Finland", FormalName: "Republic of Finland", Abbrev: "Fin.", NameAlt: "", Population: 5520314, GDP: 269296, ISO2: "FI", ISO3: "FIN", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: 27.276449, LabelY: 63.252361},
	{Name: "Luxembourg", FormalName: "Grand Duchy of Luxembourg", Abbrev: "Lux.", NameAlt: "", Population: 619896, GDP: 71104, ISO2: "LU", ISO3: "LUX", Continent: "Europe", RegionUN: "Europe", Subregion: "Western Europe", RegionWB: "Europe & Central Asia", LabelX: 6.07762, LabelY: 49.733732},
	{Name: "Belgium", FormalName: "Kingdom of Belgium", Abbrev: "Belg.", NameAlt: "", Population: 11484055, GDP: 533097, ISO2: "BE", ISO3: "BEL", Continent: "Europe", RegionUN: "Europe", Subregion: "Western Europe", RegionWB: "Europe & Central Asia", LabelX: 4.800448, LabelY: 50.785392},
	{Name: "North Macedonia", FormalName: "Republic of North Macedonia", Abbrev: "N. Mac.", NameAlt: "", Population: 2083459, GDP: 12547, ISO2: "MK", ISO3: "MKD", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 21.555839, LabelY: 41.558223},
	{Name: "Albania", FormalName: "Republic of Albania", Abbrev: "Alb.", NameAlt: "", Population: 2854191, GDP: 15279, ISO2: "AL", ISO3: "ALB", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 20.11384, LabelY: 40.654855},
	{Name: "Kosovo", FormalName: "Republic of Kosovo", Abbrev: "Kos.", NameAlt: "", Population: 1794248, GDP: 7926, ISO2: "-99", ISO3: "-99", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 20.860719, LabelY: 42.593587},
	{Name: "Spain", FormalName: "Kingdom of Spain", Abbrev: "Sp.", NameAlt: "", Population: 47076781, GDP: 1393490, ISO2: "ES", ISO3: "ESP", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: -3.464718, LabelY: 40.090953},
	{Name: "Denmark", FormalName: "Kingdom of Denmark", Abbrev: "Den.", NameAlt: "", Population: 5818553, GDP: 350104, ISO2: "DK", ISO3: "DNK", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: 9.018163, LabelY: 55.966965},
	{Name: "Romania", FormalName: "Romania", Abbrev: "Rom.", NameAlt: "", Population: 19356544, GDP: 250077, ISO2: "RO", ISO3: "ROU", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 24.972624, LabelY: 45.733237},
	{Name: "Hungary", FormalName: "Republic of Hungary", Abbrev: "Hun.", NameAlt: "", Population: 9769949, GDP: 163469, ISO2: "HU", ISO3: "HUN", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 19.447867, LabelY: 47.086841},
	{Name: "Slovakia", FormalName: "Slovak Republic", Abbrev: "Svk.", NameAlt: "", Population: 5454073, GDP: 105079, ISO2: "SK", ISO3: "SVK", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 19.049868, LabelY: 48.734044},
	{Name: "Poland", FormalName: "Republic of Poland", Abbrev: "Pol.", NameAlt: "", Population: 37970874, GDP: 595858, ISO2: "PL", ISO3: "POL", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 19.490468, LabelY: 51.990316},
	{Name: "Ireland", FormalName: "Ireland", Abbrev: "Ire.", NameAlt: "", Population: 4941444, GDP: 388698, ISO2: "IE", ISO3: "IRL", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: -7.798588, LabelY: 53.078726},
	{Name: "United Kingdom", FormalName: "United Kingdom of Great Britain and Northern Ireland", Abbrev: "U.K.", NameAlt: "", Population: 66834405, GDP: 2829108, ISO2: "GB", ISO3: "GBR", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: -2.116346, LabelY: 54.402739},
	{Name: "Greece", FormalName: "Hellenic Republic", Abbrev: "Greece", NameAlt: "", Population: 10716322, GDP: 209852, ISO2: "GR", ISO3: "GRC", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 21.72568, LabelY: 39.492763},
	{Name: "Austria", FormalName: "Republic of Austria", Abbrev: "Aust.", NameAlt: "", Population: 8877067, GDP: 445075, ISO2: "AT", ISO3: "AUT", Continent: "Europe", RegionUN: "Europe", Subregion: "Western Europe", RegionWB: "Europe & Central Asia", LabelX: 14.130515, LabelY: 47.518859},
	{Name: "Italy", FormalName: "Italian Republic", Abbrev: "Italy", NameAlt: "", Population: 60297396, GDP: 2003576, ISO2: "IT", ISO3: "ITA", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 11.076907, LabelY: 44.732482},
	{Name: "Switzerland", FormalName: "Swiss Confederation", Abbrev: "Switz.", NameAlt: "", Population: 8574832, GDP: 703082, ISO2: "CH", ISO3: "CHE", Continent: "Europe", RegionUN: "Europe", Subregion: "Western Europe", RegionWB: "Europe & Central Asia", LabelX: 7.463965, LabelY: 46.719114},
	{Name: "Netherlands", FormalName: "Kingdom of the Netherlands", Abbrev: "Neth.", NameAlt: "", Population: 17332850, GDP: 907050, ISO2: "NL", ISO3: "NLD", Continent: "Europe", RegionUN: "Europe", Subregion: "Western Europe", RegionWB: "Europe & Central Asia", LabelX: 5.61144, LabelY: 52.422211},
	{Name: "Serbia", FormalName: "Republic of Serbia", Abbrev: "Serb.", NameAlt: "", Population: 6944975, GDP: 51475, ISO2: "RS", ISO3: "SRB", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 20.787989, LabelY: 44.189919},
	{Name: "Croatia", FormalName: "Republic of Croatia", Abbrev: "Cro.", NameAlt: "", Population: 4067500, GDP: 60752, ISO2: "HR", ISO3: "HRV", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 16.37241, LabelY: 45.805799},
	{Name: "Slovenia", FormalName: "Republic of Slovenia", Abbrev: "Slo.", NameAlt: "", Population: 2087946, GDP: 54174, ISO2: "SI", ISO3: "SVN", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 14.915312, LabelY: 46.06076},
	{Name: "Bulgaria", FormalName: "Republic of Bulgaria", Abbrev: "Bulg.", NameAlt: "", Population: 6975761, GDP: 68558, ISO2: "BG", ISO3: "BGR", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 25.15709, LabelY: 42.508785},
	{Name: "Montenegro", FormalName: "Montenegro", Abbrev: "Mont.", NameAlt: "", Population: 622137, GDP: 5542, ISO2: "ME", ISO3: "MNE", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 19.143727, LabelY: 42.803101},
	{Name: "Bosnia and Herz.", FormalName: "Bosnia and Herzegovina", Abbrev: "B.H.", NameAlt: "", Population: 3301000, GDP: 20164, ISO2: "BA", ISO3: "BIH", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 18.06841, LabelY: 44.091051},
	{Name: "Portugal", FormalName: "Portuguese Republic", Abbrev: "Port.", NameAlt: "", Population: 10269417, GDP: 238785, ISO2: "PT", ISO3: "PRT", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: -8.271754, LabelY: 39.606675},
	{Name: "Moldova", FormalName: "Republic of Moldova", Abbrev: "Mda.", NameAlt: "", Population: 2657637, GDP: 11968, ISO2: "MD", ISO3: "MDA", Continent: "Europe", RegionUN: "Europe", Subregion: "Eastern Europe", RegionWB: "Europe & Central Asia", LabelX: 28.487904, LabelY: 47.434999},
	{Name: "Iceland", FormalName: "Republic of Iceland", Abbrev: "Iceland", NameAlt: "", Population: 361313, GDP: 24188, ISO2: "IS", ISO3: "ISL", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: -18.673711, LabelY: 64.779286},
	{Name: "Papua New Guinea", FormalName: "Independent State of Papua New Guinea", Abbrev: "P.N.G.", NameAlt: "", Population: 8776109, GDP: 24829, ISO2: "PG", ISO3: "PNG", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Melanesia", RegionWB: "East Asia & Pacific", LabelX: 143.910216, LabelY: -5.695285},
	{Name: "Australia", FormalName: "Commonwealth of Australia", Abbrev: "Auz.", NameAlt: "", Population: 25364307, GDP: 1396567, ISO2: "AU", ISO3: "AUS", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Australia and New Zealand", RegionWB: "East Asia & Pacific", LabelX: 134.04972, LabelY: -24.129522},
	{Name: "Fiji", FormalName: "Republic of Fiji", Abbrev: "Fiji", NameAlt: "", Population: 889953, GDP: 5496, ISO2: "FJ", ISO3: "FJI", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Melanesia", RegionWB: "East Asia & Pacific", LabelX: 177.975427, LabelY: -17.826099},
	{Name: "New Zealand", FormalName: "New Zealand", Abbrev: "N.Z.", NameAlt: "", Population: 4917000, GDP: 206928, ISO2: "NZ", ISO3: "NZL", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Australia and New Zealand", RegionWB: "East Asia & Pacific", LabelX: 172.787, LabelY: -39.759},
	{Name: "New Caledonia", FormalName: "New Caledonia", Abbrev: "New C.", NameAlt: "", Population: 287800, GDP: 10770, ISO2: "NC", ISO3: "NCL", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Melanesia", RegionWB: "East Asia & Pacific", LabelX: 165.084004, LabelY: -21.064697},
	{Name: "Solomon Is.", FormalName: "", Abbrev: "S. Is.", NameAlt: "", Population: 669823, GDP: 1589, ISO2: "SB", ISO3: "SLB", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Melanesia", RegionWB: "East Asia & Pacific", LabelX: 159.170468, LabelY: -8.029548},
	{Name: "Vanuatu", FormalName: "Republic of Vanuatu", Abbrev: "Van.", NameAlt: "", Population: 299882, GDP: 934, ISO2: "VU", ISO3: "VUT", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Melanesia", RegionWB: "East Asia & Pacific", LabelX: 166.908762, LabelY: -15.37153},
	{Name: "Antarctica", FormalName: "", Abbrev: "Ant.", NameAlt: "", Population: 4490, GDP: 898, ISO2: "AQ", ISO3: "ATA", Continent: "Antarctica", RegionUN: "Antarctica", Subregion: "Antarctica", RegionWB: "Antarctica", LabelX: 35.885455, LabelY: -79.843222},
	{Name: "Fr. S. Antarctic Lands", FormalName: "Territory of the French Southern and Antarctic Lands", Abbrev: "Fr. S.A.L.", NameAlt: "", Population: 140, GDP: 16, ISO2: "TF", ISO3: "ATF", Continent: "Seven seas (open ocean)", RegionUN: "Africa", Subregion: "Seven seas (open ocean)", RegionWB: "Sub-Saharan Africa", LabelX: 69.122136, LabelY: -49.303721},
	{Name: "Rest of World", FormalName: "Rest of World", Abbrev: "RoW", NameAlt: "", Population: 0, GDP: 0, ISO2: "ZZ", ISO3: "ZZZ", Continent: "Seven seas (open ocean)", RegionUN: "Seven seas (open ocean)", Subregion: "Seven seas (open ocean)", RegionWB: "Sub-Saharan Africa", LabelX: 0, LabelY: 0},
	{Name: "Singapore", FormalName: "Republic of Singapore", Abbrev: "Sing.", NameAlt: "", Population: 5850342, GDP: 364134, ISO2: "SG", ISO3: "SGP", Continent: "Asia", RegionUN: "Asia", Subregion: "South-Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 103.819836, LabelY: 1.352083},
	{Name: "Mauritius", FormalName: "Republic of Mauritius", Abbrev: "Maur.", NameAlt: "", Population: 1265985, GDP: 14100, ISO2: "MU", ISO3: "MUS", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 57.552152, LabelY: -20.348404},
	{Name: "Malta", FormalName: "Republic of Malta", Abbrev: "Malta", NameAlt: "", Population: 514564, GDP: 14495, ISO2: "MT", ISO3: "MLT", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 14.375416, LabelY: 35.937496},
	{Name: "Cabo Verde", FormalName: "Republic of Cabo Verde", Abbrev: "C.V.", NameAlt: "", Population: 555987, GDP: 1683, ISO2: "CV", ISO3: "CPV", Continent: "Africa", RegionUN: "Africa", Subregion: "Western Africa", RegionWB: "Sub-Saharan Africa", LabelX: -23.0418, LabelY: 16.5388},
	{Name: "Hong Kong", FormalName: "Hong Kong", Abbrev: "H.K.", NameAlt: "", Population: 7451000, GDP: 366152, ISO2: "HK", ISO3: "HKG", Continent: "Asia", RegionUN: "Asia", Subregion: "Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 114.109497, LabelY: 22.396428},
	{Name: "Andorra", FormalName: "Principality of Andorra", Abbrev: "And.", NameAlt: "", Population: 77265, GDP: 3013, ISO2: "AD", ISO3: "AND", Continent: "Europe", RegionUN: "Europe", Subregion: "Southern Europe", RegionWB: "Europe & Central Asia", LabelX: 1.521801, LabelY: 42.546245},
	{Name: "Anguilla", FormalName: "Anguilla", Abbrev: "Ang.", NameAlt: "", Population: 15002, GDP: 337, ISO2: "AI", ISO3: "AIA", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -63.068615, LabelY: 18.220554},
	{Name: "Sao Tome and Principe", FormalName: "Democratic Republic of Sao Tome and Principe", Abbrev: "S.T.P.", NameAlt: "", Population: 211028, GDP: 372, ISO2: "ST", ISO3: "STP", Continent: "Africa", RegionUN: "Africa", Subregion: "Middle Africa", RegionWB: "Sub-Saharan Africa", LabelX: 6.613081, LabelY: 0.18636},
	{Name: "Netherlands Antilles", FormalName: "Netherlands Antilles", Abbrev: "Neth. A.", NameAlt: "", Population: 0, GDP: 0, ISO2: "AN", ISO3: "ANT", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -68.971535, LabelY: 12.52088},
	{Name: "Antigua and Barbuda", FormalName: "Antigua and Barbuda", Abbrev: "A.B.", NameAlt: "", Population: 97929, GDP: 1531, ISO2: "AG", ISO3: "ATG", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -61.796428, LabelY: 17.060816},
	{Name: "Aruba", FormalName: "Aruba", Abbrev: "Aruba", NameAlt: "", Population: 106766, GDP: 2930, ISO2: "AW", ISO3: "ABW", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -69.968338, LabelY: 12.52111},
	{Name: "Barbados", FormalName: "Barbados", Abbrev: "Barb.", NameAlt: "", Population: 287025, GDP: 4713, ISO2: "BB", ISO3: "BRB", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -59.559797, LabelY: 13.193887},
	{Name: "Bermuda", FormalName: "Bermuda", Abbrev: "Berm.", NameAlt: "", Population: 63918, GDP: 6267, ISO2: "BM", ISO3: "BMU", Continent: "Americas", RegionUN: "Americas", Subregion: "Northern America", RegionWB: "North America", LabelX: -64.754559, LabelY: 32.313678},
	{Name: "Bahrain", FormalName: "Kingdom of Bahrain", Abbrev: "Bahr.", NameAlt: "", Population: 1641172, GDP: 35289, ISO2: "BH", ISO3: "BHR", Continent: "Asia", RegionUN: "Asia", Subregion: "Western Asia", RegionWB: "Middle East & North Africa", LabelX: 50.541969, LabelY: 26.0667},
	{Name: "Cook Islands", FormalName: "Cook Islands", Abbrev: "C. Is.", NameAlt: "", Population: 17548, GDP: 311, ISO2: "CK", ISO3: "COK", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Polynesia", RegionWB: "East Asia & Pacific", LabelX: -159.777671, LabelY: -21.236736},
	{Name: "Dominica", FormalName: "Commonwealth of Dominica", Abbrev: "Dom.", NameAlt: "", Population: 71986, GDP: 520, ISO2: "DM", ISO3: "DMA", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -61.370976, LabelY: 15.414999},
	{Name: "Federated States of Micronesia", FormalName: "Federated States of Micronesia", Abbrev: "F.S.M.", NameAlt: "", Population: 113815, GDP: 333, ISO2: "FM", ISO3: "FSM", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Micronesia", RegionWB: "East Asia & Pacific", LabelX: 158.215071, LabelY: 6.924757},
	{Name: "Grenada", FormalName: "Grenada", Abbrev: "Gren.", NameAlt: "", Population: 112003, GDP: 1111, ISO2: "GD", ISO3: "GRD", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -61.604171, LabelY: 12.1165},
	{Name: "French Guiana", FormalName: "French Guiana", Abbrev: "Fr. G.", NameAlt: "", Population: 290691, GDP: 0, ISO2: "GF", ISO3: "GUF", Continent: "Americas", RegionUN: "Americas", Subregion: "South America", RegionWB: "Latin America & Caribbean", LabelX: -53.125782, LabelY: 3.933889},
	{Name: "Guadeloupe", FormalName: "Guadeloupe", Abbrev: "Guad.", NameAlt: "", Population: 395700, GDP: 0, ISO2: "GP", ISO3: "GLP", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -61.534042, LabelY: 16.264727},
	{Name: "Kiribati", FormalName: "Republic of Kiribati", Abbrev: "Kir.", NameAlt: "", Population: 119449, GDP: 196, ISO2: "KI", ISO3: "KIR", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Micronesia", RegionWB: "East Asia & Pacific", LabelX: -45.611105, LabelY: 0.860015},
	{Name: "Comoros", FormalName: "Union of the Comoros", Abbrev: "Com.", NameAlt: "", Population: 850886, GDP: 658, ISO2: "KM", ISO3: "COM", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 43.872219, LabelY: -11.875001},
	{Name: "St. Kitts and Nevis", FormalName: "Federation of Saint Kitts and Nevis", Abbrev: "St. K.N.", NameAlt: "", Population: 53199, GDP: 932, ISO2: "KN", ISO3: "KNA", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -62.782998, LabelY: 17.357822},
	{Name: "Cayman Is.", FormalName: "Cayman Islands", Abbrev: "Cay. Is.", NameAlt: "", Population: 65722, GDP: 3116, ISO2: "KY", ISO3: "CYM", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -81.2546, LabelY: 19.3133},
	{Name: "St. Lucia", FormalName: "Saint Lucia", Abbrev: "St. L.", NameAlt: "", Population: 182790, GDP: 1549, ISO2: "LC", ISO3: "LCA", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -60.978893, LabelY: 13.909444},
	{Name: "Macau", FormalName: "Macau", Abbrev: "Macau", NameAlt: "", Population: 631636, GDP: 0, ISO2: "MO", ISO3: "MAC", Continent: "Asia", RegionUN: "Asia", Subregion: "Eastern Asia", RegionWB: "East Asia & Pacific", LabelX: 113.543873, LabelY: 22.198745},
	{Name: "Martinique", FormalName: "Martinique", Abbrev: "Mart.", NameAlt: "", Population: 375265, GDP: 0, ISO2: "MQ", ISO3: "MTQ", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -61.024174, LabelY: 14.641528},
	{Name: "Montserrat", FormalName: "Montserrat", Abbrev: "Monts.", NameAlt: "", Population: 5900, GDP: 0, ISO2: "MS", ISO3: "MSR", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -62.187366, LabelY: 16.742498},
	{Name: "The Maldives", FormalName: "Republic of Maldives", Abbrev: "Mald.", NameAlt: "", Population: 530953, GDP: 0, ISO2: "MV", ISO3: "MDV", Continent: "Asia", RegionUN: "Asia", Subregion: "Southern Asia", RegionWB: "South Asia", LabelX: 73.22068, LabelY: 3.202778},
	{Name: "French Polynesia", FormalName: "French Polynesia", Abbrev: "Fr. P.", NameAlt: "", Population: 280208, GDP: 0, ISO2: "PF", ISO3: "PYF", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Polynesia", RegionWB: "East Asia & Pacific", LabelX: -149.406843, LabelY: -17.679742},
	{Name: "Palau", FormalName: "Republic of Palau", Abbrev: "Pal.", NameAlt: "", Population: 17907, GDP: 0, ISO2: "PW", ISO3: "PLW", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Micronesia", RegionWB: "East Asia & Pacific", LabelX: 134.58252, LabelY: 7.51498},
	{Name: "Reunion", FormalName: "Reunion", Abbrev: "Reun.", NameAlt: "", Population: 859959, GDP: 0, ISO2: "RE", ISO3: "REU", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 55.536384, LabelY: -21.115141},
	{Name: "Seycheles", FormalName: "Republic of Seychelles", Abbrev: "Sey.", NameAlt: "", Population: 96762, GDP: 0, ISO2: "SC", ISO3: "SYC", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 55.491977, LabelY: -4.679574},
	{Name: "Turks and Caicos Is.", FormalName: "Turks and Caicos Islands", Abbrev: "T.C. Is.", NameAlt: "", Population: 35446, GDP: 0, ISO2: "TC", ISO3: "TCA", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -71.797928, LabelY: 21.694025},
	{Name: "Tuvalu", FormalName: "Tuvalu", Abbrev: "Tuva.", NameAlt: "", Population: 11646, GDP: 0, ISO2: "TV", ISO3: "TUV", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Polynesia", RegionWB: "East Asia & Pacific", LabelX: 179.198566, LabelY: -7.109535},
	{Name: "St. Vincent and the Grenadines", FormalName: "Saint Vincent and the Grenadines", Abbrev: "St. V.G.", NameAlt: "", Population: 110211, GDP: 0, ISO2: "VC", ISO3: "VCT", Continent: "Americas", RegionUN: "Americas", Subregion: "Caribbean", RegionWB: "Latin America & Caribbean", LabelX: -61.287228, LabelY: 12.984305},
	{Name: "Wallis and Futuna Is.", FormalName: "Wallis and Futuna Islands", Abbrev: "W.F. Is.", NameAlt: "", Population: 15289, GDP: 0, ISO2: "WF", ISO3: "WLF", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Polynesia", RegionWB: "East Asia & Pacific", LabelX: -177.156097, LabelY: -13.768752},
	{Name: "Samoa", FormalName: "Independent State of Samoa", Abbrev: "Samoa", NameAlt: "", Population: 196440, GDP: 0, ISO2: "WS", ISO3: "WSM", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Polynesia", RegionWB: "East Asia & Pacific", LabelX: -172.104629, LabelY: -13.759029},
	{Name: "Mayotte", FormalName: "Mayotte", Abbrev: "May.", NameAlt: "", Population: 270372, GDP: 0, ISO2: "YT", ISO3: "MYT", Continent: "Africa", RegionUN: "Africa", Subregion: "Eastern Africa", RegionWB: "Sub-Saharan Africa", LabelX: 45.166244, LabelY: -12.8275},
	{Name: "Faroe Is.", FormalName: "Faroe Islands", Abbrev: "F. Is.", NameAlt: "", Population: 48678, GDP: 0, ISO2: "FO", ISO3: "FRO", Continent: "Europe", RegionUN: "Europe", Subregion: "Northern Europe", RegionWB: "Europe & Central Asia", LabelX: -6.911806, LabelY: 61.892635},
	{Name: "Tonga", FormalName: "Kingdom of Tonga", Abbrev: "Tonga", NameAlt: "", Population: 106501, GDP: 0, ISO2: "TO", ISO3: "TON", Continent: "Oceania", RegionUN: "Oceania", Subregion: "Polynesia", RegionWB: "East Asia & Pacific", LabelX: -175.198242, LabelY: -21.178986}}

var goguessaltnames = []string{"East Timor", "Islas Malvinas", "Swaziland", "\u010cesko"}

var goguessiso2 = []string{"CR", "NI", "HT", "DO", "SV", "GT", "CU", "HN", "US", "CA", "MX", "BZ", "PA", "GL", "BS", "TT", "PR", "JM", "ID", "MY", "CY", "IN", "CN", "IL", "PS", "LB", "SY", "KR", "KP", "BT", "OM", "UZ",
	"KZ", "TJ", "MN", "VN", "KH", "AE", "GE", "AZ", "TR", "LA", "KG", "AM", "IQ", "IR", "QA", "SA", "PK", "TH",
	"KW", "TL", "BN", "MM", "BD", "AF", "TM", "JO", "NP", "YE", "-99", "PH", "LK", "CN-TW", "JP", "CL", "BO", "PE",
	"AR", "SR", "GY", "BR", "UY", "EC", "CO", "PY", "VE", "FK", "ET", "SS", "SO", "KE", "MW", "TZ", "-99", "MA", "EH",
	"CG", "CD", "NA", "ZA", "LY", "TN", "ZM", "SL", "GN", "LR", "CF", "SD", "DJ", "ER", "CI", "ML", "SN", "NG", "BJ",
	"AO", "BW", "ZW", "TD", "DZ", "MZ", "SZ", "BI", "RW", "UG", "LS", "CM", "GA", "NE", "BF", "TG", "GH", "GW", "EG",
	"MR", "GQ", "GM", "MG", "-99", "UA", "BY", "LT", "RU", "CZ", "DE", "EE", "LV", "-99", "SE", "FI", "LU", "BE", "MK",
	"AL", "-99", "ES", "DK", "RO", "HU", "SK", "PL", "IE", "GB", "GR", "AT", "IT", "CH", "NL", "RS", "HR", "SI", "BG",
	"ME", "BA", "PT", "MD", "IS", "PG", "AU", "FJ", "NZ", "NC", "SB", "VU", "AQ", "TF", "ZZ", "SG", "MU", "MT", "CV", "HK", "AD", "AI",
	"ST"}

var goguessiso3 = []string{"CRI", "NIC", "HTI", "DOM", "SLV", "GTM", "CUB", "HND", "USA", "CAN", "MEX", "BLZ", "PAN",
	"GRL", "BHS", "TTO", "PRI", "JAM", "IDN", "MYS", "CYP", "IND", "CHN", "ISR", "PSE", "LBN", "SYR", "KOR", "PRK", "BTN", "OMN",
	"UZB", "KAZ", "TJK", "MNG", "VNM", "KHM", "ARE", "GEO", "AZE", "TUR", "LAO", "KGZ", "ARM", "IRQ", "IRN", "QAT", "SAU", "PAK",
	"THA", "KWT", "TLS", "BRN", "MMR", "BGD", "AFG", "TKM", "JOR", "NPL", "YEM", "-99", "PHL", "LKA", "TWN", "JPN", "CHL", "BOL",
	"PER", "ARG", "SUR", "GUY", "BRA", "URY", "ECU", "COL", "PRY", "VEN", "FLK", "ETH", "SSD", "SOM", "KEN", "MWI", "TZA", "-99",
	"MAR", "ESH", "COG", "COD", "NAM", "ZAF", "LBY", "TUN", "ZMB", "SLE", "GIN", "LBR", "CAF", "SDN", "DJI", "ERI", "CIV", "MLI",
	"SEN", "NGA", "BEN", "AGO", "BWA", "ZWE", "TCD", "DZA", "MOZ", "SWZ", "BDI", "RWA", "UGA", "LSO", "CMR", "GAB", "NER", "BFA",
	"TGO", "GHA", "GNB", "EGY", "MRT", "GNQ", "GMB", "MDG", "-99", "UKR", "BLR", "LTU", "RUS", "CZE", "DEU", "EST", "LVA", "-99",
	"SWE", "FIN", "LUX", "BEL", "MKD", "ALB", "-99", "ESP", "DNK", "ROU", "HUN", "SVK", "POL", "IRL", "GBR", "GRC", "AUT", "ITA",
	"CHE", "NLD", "SRB", "HRV", "SVN", "BGR", "MNE", "BIH", "PRT", "MDA", "ISL", "PNG", "AUS", "FJI", "NZL", "NCL", "SLB", "VUT",
	"ATA", "ATF", "ZZZ", "SGP", "MUS", "MLT", "CPV", "HKG", "AND", "AIA", "STP"}

var goguessFullnames = []string{"Republic of Costa Rica", "Republic of Nicaragua", "Republic of Haiti", "Dominican Republic",
	"Republic of El Salvador", "Republic of Guatemala", "Republic of Cuba", "Republic of Honduras",
	"United States of America", "Canada", "United Mexican States", "Belize", "Republic of Panama",
	"Greenland", "Commonwealth of the Bahamas", "Republic of Trinidad and Tobago", "Commonwealth of Puerto Rico",
	"Jamaica", "Republic of Indonesia", "Malaysia", "Republic of Cyprus", "Republic of India",
	"People's Republic of China", "State of Israel", "West Bank and Gaza", "Lebanese Republic", "Syrian Arab Republic",
	"Republic of Korea", "Democratic People's Republic of Korea", "Kingdom of Bhutan", "Sultanate of Oman",
	"Republic of Uzbekistan", "Republic of Kazakhstan", "Republic of Tajikistan", "Mongolia", "Socialist Republic of Vietnam",
	"Kingdom of Cambodia", "United Arab Emirates", "Georgia", "Republic of Azerbaijan", "Republic of Turkey",
	"Lao People's Democratic Republic", "Kyrgyz Republic", "Republic of Armenia", "Republic of Iraq", "Islamic Republic of Iran",
	"State of Qatar", "Kingdom of Saudi Arabia", "Islamic Republic of Pakistan", "Kingdom of Thailand", "State of Kuwait",
	"Democratic Republic of Timor-Leste", "Negara Brunei Darussalam", "Republic of the Union of Myanmar", "People's Republic of Bangladesh",
	"Islamic State of Afghanistan", "Turkmenistan", "Hashemite Kingdom of Jordan", "Nepal", "Republic of Yemen", "Turkish Republic of Northern Cyprus",
	"Republic of the Philippines", "Democratic Socialist Republic of Sri Lanka", "", "Japan", "Republic of Chile", "Plurinational State of Bolivia",
	"Republic of Peru", "Argentine Republic", "Republic of Suriname", "Co-operative Republic of Guyana", "Federative Republic of Brazil",
	"Oriental Republic of Uruguay", "Republic of Ecuador", "Republic of Colombia", "Republic of Paraguay", "Bolivarian Republic of Venezuela",
	"Falkland Islands", "Federal Democratic Republic of Ethiopia", "Republic of South Sudan", "Federal Republic of Somalia", "Republic of Kenya",
	"Republic of Malawi", "United Republic of Tanzania", "Republic of Somaliland", "Kingdom of Morocco", "Sahrawi Arab Democratic Republic",
	"Republic of the Congo", "Democratic Republic of the Congo", "Republic of Namibia", "Republic of South Africa", "Libya", "Republic of Tunisia",
	"Republic of Zambia", "Republic of Sierra Leone", "Republic of Guinea", "Republic of Liberia", "Central African Republic", "Republic of the Sudan",
	"Republic of Djibouti", "State of Eritrea", "Republic of Ivory Coast", "Republic of Mali", "Republic of Senegal", "Federal Republic of Nigeria",
	"Republic of Benin", "People's Republic of Angola", "Republic of Botswana", "Republic of Zimbabwe", "Republic of Chad", "People's Democratic Republic of Algeria",
	"Republic of Mozambique", "Kingdom of eSwatini", "Republic of Burundi", "Republic of Rwanda", "Republic of Uganda", "Kingdom of Lesotho",
	"Republic of Cameroon", "Gabonese Republic", "Republic of Niger", "Burkina Faso", "Togolese Republic", "Republic of Ghana", "Republic of Guinea-Bissau",
	"Arab Republic of Egypt", "Islamic Republic of Mauritania", "Republic of Equatorial Guinea", "Republic of the Gambia", "Republic of Madagascar",
	"French Republic", "Ukraine", "Republic of Belarus", "Republic of Lithuania", "Russian Federation", "Czech Republic", "Federal Republic of Germany",
	"Republic of Estonia", "Republic of Latvia", "Kingdom of Norway", "Kingdom of Sweden", "Republic of Finland", "Grand Duchy of Luxembourg",
	"Kingdom of Belgium", "Republic of North Macedonia", "Republic of Albania", "Republic of Kosovo", "Kingdom of Spain", "Kingdom of Denmark",
	"Romania", "Republic of Hungary", "Slovak Republic", "Republic of Poland", "Ireland", "United Kingdom of Great Britain and Northern Ireland",
	"Hellenic Republic", "Republic of Austria", "Italian Republic", "Swiss Confederation", "Kingdom of the Netherlands", "Republic of Serbia",
	"Republic of Croatia", "Republic of Slovenia", "Republic of Bulgaria", "Montenegro", "Bosnia and Herzegovina", "Portuguese Republic",
	"Republic of Moldova", "Republic of Iceland", "Independent State of Papua New Guinea", "Commonwealth of Australia", "Republic of Fiji",
	"New Zealand", "New Caledonia", "", "Republic of Vanuatu", "", "Territory of the French Southern and Antarctic Lands", "Rest of World", "Republic of Singapore", "Republic of Mauritius"}

var goguessCountrylist = []string{"Costa Rica", "Nicaragua", "Haiti", "Dominican Rep.", "El Salvador", "Guatemala",
	"Cuba", "Honduras", "United States", "Canada", "Mexico", "Belize", "Panama",
	"Greenland", "Bahamas", "Trinidad and Tobago", "Puerto Rico", "Jamaica",
	"Indonesia", "Malaysia", "Cyprus", "India", "China", "Israel", "Palestine",
	"Lebanon", "Syria", "South Korea", "North Korea", "Bhutan", "Oman",
	"Uzbekistan", "Kazakhstan", "Tajikistan", "Mongolia", "Vietnam", "Cambodia",
	"United Arab Emirates", "Georgia", "Azerbaijan", "Turkey", "Laos", "Kyrgyzstan",
	"Armenia", "Iraq", "Iran", "Qatar", "Saudi Arabia", "Pakistan", "Thailand",
	"Kuwait", "Timor-Leste", "Brunei", "Myanmar", "Bangladesh", "Afghanistan",
	"Turkmenistan", "Jordan", "Nepal", "Yemen", "N. Cyprus", "Philippines",
	"Sri Lanka", "Taiwan", "Japan", "Chile", "Bolivia", "Peru", "Argentina",
	"Suriname", "Guyana", "Brazil", "Uruguay", "Ecuador", "Colombia", "Paraguay",
	"Venezuela", "Falkland Is.", "Ethiopia", "S. Sudan", "Somalia", "Kenya",
	"Malawi", "Tanzania", "Somaliland", "Morocco", "W. Sahara", "Congo",
	"Dem. Rep. Congo", "Namibia", "South Africa", "Libya", "Tunisia", "Zambia",
	"Sierra Leone", "Guinea", "Liberia", "Central African Rep.", "Sudan",
	"Djibouti", "Eritrea", "Côte d'Ivoire", "Mali", "Senegal", "Nigeria",
	"Benin", "Angola", "Botswana", "Zimbabwe", "Chad", "Algeria", "Mozambique",
	"eSwatini", "Burundi", "Rwanda", "Uganda", "Lesotho", "Cameroon", "Gabon",
	"Niger", "Burkina Faso", "Togo", "Ghana", "Guinea-Bissau", "Egypt",
	"Mauritania", "Eq. Guinea", "Gambia", "Madagascar", "France", "Ukraine",
	"Belarus", "Lithuania", "Russia", "Czechia", "Germany", "Estonia", "Latvia",
	"Norway", "Sweden", "Finland", "Luxembourg", "Belgium", "North Macedonia",
	"Albania", "Kosovo", "Spain", "Denmark", "Romania", "Hungary", "Slovakia",
	"Poland", "Ireland", "United Kingdom", "Greece", "Austria", "Italy",
	"Switzerland", "Netherlands", "Serbia", "Croatia", "Slovenia", "Bulgaria",
	"Montenegro", "Bosnia and Herz.", "Portugal", "Moldova", "Iceland",
	"Papua New Guinea", "Australia", "Fiji", "New Zealand", "New Caledonia",
	"Solomon Is.", "Vanuatu", "Antarctica", "Fr. S. Antarctic Lands", "Rest of World",
	"Singapore", "Mauritius", "Malta", "Cabo Verde", "Andorra", "Anguilla", "Sao Tome and Principe"}

// -------------------------------------------------------------------------
// Create the quick iso2 lookup
// -------------------------------------------------------------------------
func makeISO2LUT() (map[string]string, map[string]string) {

	iso2lut := make(map[string]string, 0)
	countrylut := make(map[string]string, 0)

	for _, c := range goguesscountries {

		iso2 := strings.ToLower(c.ISO2)
		country := c.Name

		iso2lut[iso2] = country
		countrylut[country] = iso2
	}

	return countrylut, iso2lut

}

// -------------------------------------------------------------------------
func Iso2(code string) (Country, error) {

	var outcountry Country
	code = strings.ToUpper(code)

	for _, c := range goguesscountries {
		if c.ISO2 == code {
			return c, nil
		}
	}
	return outcountry, errors.New("ISO2 not matched")
}

func Iso3(code string) (Country, error) {

	var outcountry Country
	code = strings.ToUpper(code)

	for _, c := range goguesscountries {
		if c.ISO3 == code {
			return c, nil
		}
	}
	return outcountry, errors.New("ISO3 not matched")
}

// ##########################################################################
// The GuessCountry function - give it a string and it will return the best matched struct
// ##########################################################################
func GuessCountry(name string) (Country, error) {

	var outcountry Country
	stringlen := utf8.RuneCountInString(name)

	// check to see if string is 2 characters long
	if stringlen == 2 { // is it an iso2?

		name = strings.ToUpper(name)
		if slices.Contains(goguessiso2, name) {
			outcountry, err := Iso2(name)
			if err == nil {
				return outcountry, nil
			}
		}
	}

	// check to see if string is 3 characters long
	if stringlen == 3 { // is it an iso3?
		name = strings.ToUpper(name)
		if slices.Contains(goguessiso3, name) {
			outcountry, err := Iso3(name)
			if err == nil {
				return outcountry, nil
			}
		}
	}

	// need to do some fuzzy searches

	namematch := fuzzy.RankFindFold(name, goguessCountrylist)
	sort.Sort(namematch)
	// check to see if the top match is perfect
	for _, n := range namematch {
		fmt.Println(n.Source, n.Target, n.Distance)
	}

	return outcountry, nil
}
