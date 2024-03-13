package errors

type ServerError string

func (e ServerError) Error() string {
	return string(e)
}

type HandlerError string

func (e HandlerError) Error() string {
	return string(e)
}

type APIError string

func (e APIError) Error() string {
	return string(e)
}
