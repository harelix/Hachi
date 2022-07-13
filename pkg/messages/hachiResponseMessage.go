package messages

type DefaultResponseMessage struct {
	Error bool   `json:"error" xml:"error"`
	Data  string `json:"data" xml:"data"`
	Path  string `json:"path" xml:"path"`
}
