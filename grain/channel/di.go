package channel

type EventHandler func(touser string, channel, user, msg, enter, quit string)

func RegisterEventHandler(h EventHandler) {
	handler = h
}

func OnEvent(touser string, channel, user, msg, enter, quit string) {
	handler(touser, channel, user, msg, enter, quit)
}

var (
	handler EventHandler
)
