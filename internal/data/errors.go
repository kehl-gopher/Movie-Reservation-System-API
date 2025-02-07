package data

type ServerError struct {
	Err error
}

type ValidationError struct {
	Err error
}

type BadRequestError struct {
	Err error
}

func (s *ServerError) Error() string {
	return s.Err.Error()
}

func (v *ValidationError) Error() string {
	return v.Err.Error()
}

func (b *BadRequestError) Error() string {
	return b.Err.Error()
}
