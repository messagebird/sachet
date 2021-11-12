package botgolang

var (
	voiceMessageSupportedExtensions = map[string]bool{
		".aac": true,
		".ogg": true,
		".m4a": true,
	}
)

const (
	voiceMessageLeadingRune = 'I'
)
