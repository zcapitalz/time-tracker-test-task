package domain

type UserNotFoundError struct {
	Message string
}

func (err UserNotFoundError) Error() string {
	return err.Message
}

type UserAlreadyExistsError struct {
	Message string
}

func (err UserAlreadyExistsError) Error() string {
	return err.Message
}

type TaskNotFoundError struct {
	Message string
}

func (err TaskNotFoundError) Error() string {
	return err.Message
}

type IncorrectPeriodError struct {
	Message string
}

func (err IncorrectPeriodError) Error() string {
	return err.Message
}
