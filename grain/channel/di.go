package channel

type Observe func(touser string, channel, user, msg, enter, quit string)

func RegisterObserver(observe Observe) {
	observeFunc = observe
}

func Publish(touser string, channel, user, msg, enter, quit string) {
	observeFunc(touser, channel, user, msg, enter, quit)
}

var (
	observeFunc Observe
)
