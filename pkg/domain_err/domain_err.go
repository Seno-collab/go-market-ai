package domainerr

type AppError struct {
	Msg    string
	Status int
}

func (e AppError) Error() string {
	return e.Msg
}

func New(status int, msg string) AppError {
	return AppError{
		Status: status,
		Msg:    msg,
	}
}
