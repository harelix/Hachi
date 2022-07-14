package messages

type DefaultResponseMessage struct {
	Error bool `json:"error" xml:"error"`
	Data  any  `json:"data" xml:"data"`
	Path  any  `json:"path" xml:"path"`
}
