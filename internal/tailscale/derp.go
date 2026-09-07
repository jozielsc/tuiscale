package tailscale

// KnownDERPRegions mapeia IDs de regiões DERP conhecidas do Tailscale para cidades/países.
var KnownDERPRegions = map[int]struct {
	Code    string
	City    string
	Country string
}{
	1:  {Code: "nyc", City: "New York", Country: "EUA"},
	2:  {Code: "sfo", City: "San Francisco", Country: "EUA"},
	3:  {Code: "sin", City: "Singapore", Country: "Cingapura"},
	4:  {Code: "fra", City: "Frankfurt", Country: "Alemanha"},
	5:  {Code: "syd", City: "Sydney", Country: "Austrália"},
	6:  {Code: "blr", City: "Bangalore", Country: "Índia"},
	7:  {Code: "tok", City: "Tokyo", Country: "Japão"},
	8:  {Code: "lhr", City: "London", Country: "Reino Unido"},
	9:  {Code: "dfw", City: "Dallas", Country: "EUA"},
	10: {Code: "sea", City: "Seattle", Country: "EUA"},
	11: {Code: "sao", City: "São Paulo", Country: "Brasil"},
	12: {Code: "par", City: "Paris", Country: "França"},
	13: {Code: "den", City: "Denver", Country: "EUA"},
	14: {Code: "ord", City: "Chicago", Country: "EUA"},
	15: {Code: "jnb", City: "Johannesburg", Country: "África do Sul"},
	16: {Code: "mia", City: "Miami", Country: "EUA"},
	17: {Code: "was", City: "Washington DC", Country: "EUA"},
	18: {Code: "hnl", City: "Honolulu", Country: "EUA"},
	19: {Code: "arn", City: "Stockholm", Country: "Suécia"},
	20: {Code: "dub", City: "Dublin", Country: "Irlanda"},
	21: {Code: "mad", City: "Madrid", Country: "Espanha"},
	22: {Code: "vno", City: "Vilnius", Country: "Lituânia"},
	23: {Code: "hkg", City: "Hong Kong", Country: "Hong Kong"},
	24: {Code: "waw", City: "Warsaw", Country: "Polônia"},
	25: {Code: "dxb", City: "Dubai", Country: "Emirados Árabes"},
	26: {Code: "nbo", City: "Nairobi", Country: "Quênia"},
	27: {Code: "tor", City: "Toronto", Country: "Canadá"},
	28: {Code: "atl", City: "Atlanta", Country: "EUA"},
}

// GetDERPLocation retorna o código, cidade e país de uma região DERP.
func GetDERPLocation(regionID int) (code, city, country string) {
	if info, ok := KnownDERPRegions[regionID]; ok {
		return info.Code, info.City, info.Country
	}
	return "derp", "Região DERP", "Global"
}

// GetDERPLocationByCode retorna a cidade e o país de uma região DERP pelo código textual.
func GetDERPLocationByCode(code string) (city, country string) {
	for _, info := range KnownDERPRegions {
		if info.Code == code {
			return info.City, info.Country
		}
	}
	return "Servidor DERP", "Global"
}
