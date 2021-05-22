package botgolang

type BotOption interface {
	Type() string
	Value() interface{}
}

type BotApiURL string

func (o BotApiURL) Type() string {
	return "api_url"
}

func (o BotApiURL) Value() interface{} {
	return string(o)
}

type BotDebug bool

func (o BotDebug) Type() string {
	return "debug"
}

func (o BotDebug) Value() interface{} {
	return bool(o)
}
