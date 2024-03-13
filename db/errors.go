package db

type DBError string

func (e DBError) Error() string {
	return string(e)
}
