package tailscale

import "testing"

func TestGetDERPLocation(t *testing.T) {
	code, city, country := GetDERPLocation(11)
	if code != "sao" || city != "São Paulo" || country != "Brasil" {
		t.Errorf("GetDERPLocation(11) = (%s, %s, %s); want (sao, São Paulo, Brasil)", code, city, country)
	}

	code, city, country = GetDERPLocation(9999)
	if code != "derp" || country != "Global" {
		t.Errorf("GetDERPLocation(9999) fallback failed, got code=%s, country=%s", code, country)
	}
}

func TestGetDERPLocationByCode(t *testing.T) {
	city, country := GetDERPLocationByCode("sao")
	if city != "São Paulo" || country != "Brasil" {
		t.Errorf("GetDERPLocationByCode(sao) = (%s, %s); want (São Paulo, Brasil)", city, country)
	}

	city, country = GetDERPLocationByCode("invalid_code")
	if city != "Servidor DERP" || country != "Global" {
		t.Errorf("GetDERPLocationByCode fallback failed, got (%s, %s)", city, country)
	}
}
