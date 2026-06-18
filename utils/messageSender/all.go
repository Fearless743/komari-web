package messageSender

import (
	_ "github.com/Fearless743/komari/utils/messageSender/bark"
	_ "github.com/Fearless743/komari/utils/messageSender/email"
	_ "github.com/Fearless743/komari/utils/messageSender/empty"
	_ "github.com/Fearless743/komari/utils/messageSender/serverchan3"
	_ "github.com/Fearless743/komari/utils/messageSender/serverchanturbo"
	_ "github.com/Fearless743/komari/utils/messageSender/javascript"
	_ "github.com/Fearless743/komari/utils/messageSender/telegram"
	_ "github.com/Fearless743/komari/utils/messageSender/webhook"
)

func All() {
}
