package ddns

import (
	_ "github.com/Fearless743/komari/utils/ddns/cloudflare"
	_ "github.com/Fearless743/komari/utils/ddns/empty"
	_ "github.com/Fearless743/komari/utils/ddns/webhook"
)

func All() {
	// empty function to ensure all DDNS providers are registered
}
