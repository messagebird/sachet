package messagebird

// Error holds details including error code, human readable description and optional parameter that is related to the error.
type Error struct {
	Code        int
	Description string
	Parameter   string
}
