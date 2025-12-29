package core

// CountryNames contains name datasets for different countries
var CountryNames = map[string][]string{
	// Middle East
	"iran":           {"ali", "reza", "mohammad", "hossein", "mahdi", "amir", "sara", "maryam"},
	"turkey":         {"ahmet", "mehmet", "mustafa", "ali", "murat", "ayse", "zeynep"},
	"saudi_arabia":   {"mohammed", "abdullah", "ahmed", "khaled", "omar", "fatima"},
	"uae":            {"mohammed", "ahmed", "rashid", "abdulrahman", "aisha"},
	"egypt":          {"ahmed", "mohamed", "mahmoud", "islam", "sara", "noor"},
	"jordan":         {"omar", "ahmad", "yusuf", "khaled", "rana"},
	"iraq":           {"ali", "hassan", "hussein", "mustafa", "zahra"},

	// Europe
	"united_states":  {"john", "michael", "david", "james", "robert", "emily", "sarah", "jessica"},
	"germany":        {"max", "lukas", "leon", "paul", "tim", "jan", "anna", "sophie"},
	"france":         {"pierre", "jean", "lucas", "thomas", "nicolas", "marie", "camille"},
	"united_kingdom": {"jack", "oliver", "harry", "george", "james", "emma", "sophie"},
	"italy":          {"luca", "marco", "andrea", "francesco", "giulia", "sara"},
	"spain":          {"juan", "carlos", "miguel", "jose", "antonio", "maria", "laura"},
	"portugal":       {"joao", "miguel", "tiago", "pedro", "ana", "sofia"},
	"netherlands":    {"jan", "piet", "tom", "lucas", "emma", "lisa"},
	"belgium":        {"luc", "thomas", "nicolas", "pierre", "marie", "anne"},
	"russia":         {"alexey", "dmitry", "ivan", "sergey", "anna", "maria"},
	"ukraine":        {"andriy", "oleksandr", "dmytro", "olena", "iryna"},
	"poland":         {"piotr", "pawel", "lukasz", "mateusz", "anna"},
	"czech_republic": {"jan", "petr", "tomas", "martin", "eva"},
	"romania":        {"andrei", "mihai", "alexandru", "ioana"},
	"sweden":         {"erik", "karl", "johan", "lina", "emma"},
	"norway":         {"ole", "lars", "magnus", "anna", "ida"},
	"finland":        {"mikko", "jari", "sami", "laura", "emilia"},

	// Asia
	"china":          {"wei", "jun", "ming", "lei", "yan", "li"},
	"japan":          {"hiroshi", "takashi", "kenji", "yuki", "sakura"},
	"south_korea":    {"minjun", "jihun", "seojun", "jisoo", "hana"},
	"india":          {"rahul", "amit", "rohit", "vikram", "priya", "anjali"},
	"pakistan":       {"ali", "ahmed", "usman", "imran", "fatima"},
	"indonesia":      {"agus", "budi", "andi", "rina", "sari"},
	"malaysia":       {"ahmad", "faizal", "hafiz", "nur", "aisyah"},
	"thailand":       {"somchai", "anan", "chai", "malee"},
	"vietnam":        {"nguyen", "minh", "tuan", "anh", "linh"},

	// Americas
	"brazil":         {"joao", "pedro", "lucas", "gabriel", "ana", "maria"},
	"argentina":      {"juan", "martin", "nicolas", "diego", "lucia"},
	"chile":          {"carlos", "felipe", "diego", "francisco", "camila"},
	"colombia":       {"andres", "juan", "carlos", "felipe", "laura"},
	"mexico":         {"juan", "jose", "luis", "miguel", "ana"},
	"canada":         {"daniel", "matthew", "joshua", "emily", "sophia"},
	"australia":      {"jack", "oliver", "noah", "liam", "emma"},

	// Africa
	"south_africa":   {"thabo", "sipho", "andile", "nomsa"},
	"nigeria":        {"emeka", "chinedu", "ibrahim", "musa", "fatima"},
	"kenya":          {"john", "peter", "james", "mary", "grace"},
}

// GetCountryNames returns names for selected countries
func GetCountryNames(countries []string) []string {
	result := []string{}
	seen := make(map[string]bool)

	for _, country := range countries {
		if names, ok := CountryNames[country]; ok {
			for _, name := range names {
				// Deduplicate at this level
				if !seen[name] {
					result = append(result, name)
					seen[name] = true
				}
			}
		}
	}

	return result
}

// GetAllCountries returns list of all available country codes
func GetAllCountries() []string {
	countries := make([]string, 0, len(CountryNames))
	for country := range CountryNames {
		countries = append(countries, country)
	}
	return countries
}

// GetCountryDisplayName returns a human-readable name for country code
func GetCountryDisplayName(countryCode string) string {
	displayNames := map[string]string{
		"iran":           "🇮🇷 Iran",
		"united_states":  "🇺🇸 United States",
		"germany":        "🇩🇪 Germany",
		"france":         "🇫🇷 France",
		"united_kingdom": "🇬🇧 United Kingdom",
		"italy":          "🇮🇹 Italy",
		"spain":          "🇪🇸 Spain",
		"portugal":       "🇵🇹 Portugal",
		"netherlands":    "🇳🇱 Netherlands",
		"belgium":        "🇧🇪 Belgium",
		"turkey":         "🇹🇷 Turkey",
		"saudi_arabia":   "🇸🇦 Saudi Arabia",
		"uae":            "🇦🇪 UAE",
		"egypt":          "🇪🇬 Egypt",
		"jordan":         "🇯🇴 Jordan",
		"iraq":           "🇮🇶 Iraq",
		"russia":         "🇷🇺 Russia",
		"ukraine":        "🇺🇦 Ukraine",
		"poland":         "🇵🇱 Poland",
		"czech_republic": "🇨🇿 Czech Republic",
		"romania":        "🇷🇴 Romania",
		"china":          "🇨🇳 China",
		"japan":          "🇯🇵 Japan",
		"south_korea":    "🇰🇷 South Korea",
		"india":          "🇮🇳 India",
		"pakistan":       "🇵🇰 Pakistan",
		"indonesia":      "🇮🇩 Indonesia",
		"malaysia":       "🇲🇾 Malaysia",
		"thailand":       "🇹🇭 Thailand",
		"vietnam":        "🇻🇳 Vietnam",
		"brazil":         "🇧🇷 Brazil",
		"argentina":      "🇦🇷 Argentina",
		"chile":          "🇨🇱 Chile",
		"colombia":       "🇨🇴 Colombia",
		"mexico":         "🇲🇽 Mexico",
		"south_africa":   "🇿🇦 South Africa",
		"nigeria":        "🇳🇬 Nigeria",
		"kenya":          "🇰🇪 Kenya",
		"australia":      "🇦🇺 Australia",
		"canada":         "🇨🇦 Canada",
		"sweden":         "🇸🇪 Sweden",
		"norway":         "🇳🇴 Norway",
		"finland":        "🇫🇮 Finland",
	}

	if name, ok := displayNames[countryCode]; ok {
		return name
	}
	return countryCode
}

