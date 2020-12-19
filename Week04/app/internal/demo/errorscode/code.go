package errorscode

type ErrorCode struct {
	Code           int
	Message        string
	HTTPStatusCode int
}

func (e ErrorCode) Error() string {
	return e.Message
}

func NewErrorCode(code int, message string, status int) error {
	return ErrorCode{
		Code:           code,
		Message:        message,
		HTTPStatusCode: status,
	}
}
