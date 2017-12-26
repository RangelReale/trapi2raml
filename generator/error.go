package trapi2ramlgen

type ErrWarning struct {
	Message string
}

func NewErrWarning(message string) *ErrWarning {
	return &ErrWarning{
		Message: message,
	}
}

func (ew *ErrWarning) Error() string {
	return ew.Message
}
