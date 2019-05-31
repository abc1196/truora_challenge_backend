package utils

// Exception data type
type Exception struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ThrowException - creates and returns a new Exception
func ThrowException(code int, message string) *Exception {
	e := new(Exception)
	e.Code = code
	e.Message = message
	return e
}
