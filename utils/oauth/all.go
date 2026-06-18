package oauth

import (
	_ "github.com/Fearless743/komari/utils/oauth/cloudflare"
	_ "github.com/Fearless743/komari/utils/oauth/factory"
	_ "github.com/Fearless743/komari/utils/oauth/generic"
	_ "github.com/Fearless743/komari/utils/oauth/github"
	_ "github.com/Fearless743/komari/utils/oauth/qq"
)

func All() {
	//empty function to ensure all OIDC providers are registered
}