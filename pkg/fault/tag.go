package fault

import (
	"errors"
)

type Tag string

const (
	UNTAGGED              Tag = "UNTAGGED"
	BAD_REQUEST           Tag = "BAD_REQUEST"
	NOT_FOUND             Tag = "NOT_FOUND"
	INTERNAL_SERVER_ERROR Tag = "INTERNAL_SERVER_ERROR"
	UNAUTHORIZED          Tag = "UNAUTHORIZED"
	FORBIDDEN             Tag = "FORBIDDEN"
	CONFLICT              Tag = "CONFLICT"
	TOO_MANY_REQUESTS     Tag = "TOO_MANY_REQUESTS"
	UNPROCESSABLE_ENTITY  Tag = "UNPROCESSABLE_ENTITY"
	LOCKED_USER           Tag = "LOCKED_USER"
	DISABLED_USER         Tag = "DISABLED_USER"
	DB_RESOURCE_NOT_FOUND Tag = "DB_RESOURCE_NOT_FOUND"
	INVALID_ENTITY        Tag = "INVALID_ENTITY"
)

// GetTag returns the first tag of the error
//
// Example:
//
//	err := fault.NewBadRequest("invalid request")
//	tag := fault.GetTag(err)
//	fmt.Println(tag) // Output: BAD_REQUEST
//
// Example with switch:
//
//	switch fault.GetTag(err) {
//	case fault.BAD_REQUEST:
//		fmt.Println("bad request")
//	case fault.NOT_FOUND:
//		fmt.Println("not found")
//	default:
//		fmt.Println("unknown error")
//	}
func GetTag(err error) Tag {
	if err == nil {
		return UNTAGGED
	}

	for err != nil {
		e, ok := err.(*Fault)
		if ok && e.Tag != "" {
			return e.Tag
		}
		err = errors.Unwrap(err)
	}

	return UNTAGGED
}
