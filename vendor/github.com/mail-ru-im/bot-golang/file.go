package botgolang

//go:generate easyjson -all file.go

type File struct {
	// Id of the file
	ID string `json:"fileId"`

	// Type of the file
	Type string `json:"type"`

	// Size in bytes
	Size uint64 `json:"size"`

	// Name of file
	Name string `json:"filename"`

	// URL to the file
	URL string `json:"url"`
}
