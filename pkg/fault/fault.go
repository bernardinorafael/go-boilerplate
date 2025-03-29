package fault

import "fmt"

type Fault struct {
	HTTPCode int    `json:"-"`
	Err      error  `json:"-"`
	Tag      Tag    `json:"tag"`
	Message  string `json:"message"`
}

func New(httpCode int, tag Tag, msg string, err error) Fault {
	return Fault{
		HTTPCode: httpCode,
		Err:      err,
		Tag:      tag,
		Message:  msg,
	}
}

func (f Fault) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", f.Tag, f.Message, f.Err)
	}
	return fmt.Sprintf("%s: %s", f.Tag, f.Message)
}

func (f Fault) StatusCode() int {
	return f.HTTPCode
}
