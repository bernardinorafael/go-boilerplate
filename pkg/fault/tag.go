package fault

import "errors"

type Tag string

const (
	UNTAGGED              Tag = "UNTAGGED"
	BAD_REQUEST           Tag = "BAD_REQUEST"
	NOT_FOUND             Tag = "RESOURCE_NOT_FOUND"
	LOCKED_USER           Tag = "LOCKED_USER"
	UNAUTHORIZED          Tag = "UNAUTHORIZED"
	DISABLED_USER         Tag = "DISABLED_USER"
	INTERNAL_SERVER_ERROR Tag = "INTERNAL_SERVER_ERROR"
	UNPROCESSABLE_ENTITY  Tag = "UNPROCESSABLE_ENTITY"
	RESOURCE_TAKEN        Tag = "RESOURCE_ALREADY_TAKEN"
	CONFLICT              Tag = "CONFLICT"
)

func GetTag(err error) Tag {
	if err == nil {
		return UNTAGGED
	}

	for err != nil {
		if e, ok := err.(*Fault); ok {
			if ok && e.Tag != "" {
				return e.Tag
			}
		}
		err = errors.Unwrap(err)
	}

	return UNTAGGED
}
