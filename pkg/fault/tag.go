package fault

import (
	"errors"
)

type Tag string

const (
	Untagged            Tag = "UNTAGGED"
	BadRequest          Tag = "BAD_REQUEST_ERROR"
	NotFound            Tag = "NOT_FOUND_ERROR"
	InternalServerError Tag = "INTERNAL_SERVER_ERROR"
	Unauthorized        Tag = "UNAUTHORIZED_ERROR"
	Forbidden           Tag = "FORBIDDEN_ERROR"
	Conflict            Tag = "CONFLICT_ERROR"
	TooManyRequests     Tag = "TOO_MANY_REQUESTS_ERROR"
	ValidationError     Tag = "VALIDATION_ERROR"
	UnprocessableEntity Tag = "UNPROCESSABLE_ENTITY_ERROR"
	LockedUser          Tag = "LOCKED_USER_ERROR"
	DisabledUser        Tag = "DISABLED_USER_ERROR"
	InvalidEntity       Tag = "INVALID_ENTITY_ERROR"
	MailerError         Tag = "MAILER_ERROR"
	Expired             Tag = "EXPIRED_ERROR"
	CacheMiss           Tag = "CACHE_MISS_KEY_ERROR"
	DBTransaction       Tag = "DB_TRANSACTION_ERROR"
	InvalidBody         Tag = "INVALID_REQUEST_BODY"
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
		return Untagged
	}

	for err != nil {
		var f *Fault
		ok := errors.As(err, &f)
		if ok && f.Tag != "" {
			return f.Tag
		}
		err = errors.Unwrap(err)
	}

	return Untagged
}
