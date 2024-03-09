package customerr

type Error string

func (e Error) Error() string {
	return string(e)
}

// USE IN DIFF DIRS OR JUST HERE
const (
	UserNotOwner = Error("user is not owner of the place")
)
