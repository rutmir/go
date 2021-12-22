package types

type ServiceError struct {
	Message string
	Code    int
	Errors  []interface{}
}

func (e *ServiceError) Error() string {
	return e.Message
}
