package fault

import (
	"errors"
)

type Tag string

const (
	Untagged            Tag = "untagged"
	BadRequest          Tag = "bad_request_error"
	NotFound            Tag = "not_found_error"
	InternalServerError Tag = "internal_server_error"
	Unauthorized        Tag = "unauthorized_error"
	Forbidden           Tag = "forbidden_error"
	Conflict            Tag = "conflict_error"
	TooManyRequests     Tag = "too_many_requests_error"
	ValidationError     Tag = "validation_error"
	UnprocessableEntity Tag = "unprocessable_entity_error"
	LockedUser          Tag = "locked_user_error"
	DisabledUser        Tag = "disabled_user_error"
	DBResourceNotFound  Tag = "db_resource_not_found_error"
	InvalidEntity       Tag = "invalid_entity_error"
	MailerError         Tag = "mailer_error"
	Expired             Tag = "expired_error"
	CacheMiss           Tag = "cache_miss_key_error"
	DBTransaction       Tag = "db_transaction_error"
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
