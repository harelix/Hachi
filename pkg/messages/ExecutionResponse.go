package messages

type ExecutionResponseCode int

const (
	Ok ExecutionResponseCode = iota
	Failure
	Error
	Unknown
)

func (er ExecutionResponseCode) String() string {
	switch er {
	case Ok:
		return "ok"
	case Failure:
		return "failure"
	case Error:
		return "error"
	case Unknown:
		return "unknown"
	}
	return "unknown"
}

type ExecutionResponse struct {
	response ExecutionResponseCode
}
