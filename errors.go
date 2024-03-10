package main

type HandlerError string

func (e HandlerError) Error() string {
	return string(e)
}

type APIError string

func (e APIError) Error() string {
	return string(e)
}

type ComponentError string

func (e ComponentError) Error() string {
	return string(e)
}

type ServerError string

func (e ServerError) Error() string {
	return string(e)
}
