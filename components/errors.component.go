package components

type ComponentError string

func (e ComponentError) Error() string {
	return string(e)
}
