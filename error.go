package sqly

import "errors"

// errors
var (
	// ErrQueryFmt sql statement format error
	ErrQueryFmt = errors.New("query can't be formatted")

	// ErrArgType sql statement format type error
	ErrArgType = errors.New("invalid variable type for argument")

	// ErrStatement sql syntax error
	ErrStatement = errors.New("sql statement syntax error")

	// ErrContainer  container for results
	ErrContainer = errors.New("invalid container for scanning (struct pointer, not nil)")

	// ErrFieldsMatch fields not match
	ErrFieldsMatch = errors.New("queried fields not match with struct fields")

	// ErrMultiRes multi result for get
	ErrMultiRes = errors.New("get more than one results for get query")

	// ErrEmpty empty
	ErrEmpty = errors.New("no result for get query ")

	// ErrCapsule Invalid Capsule
	ErrCapsule = errors.New("query capsule is not available")

	ErrEmptyArrayInStatement = errors.New("has empty array in query arguments")

	// ErrNotSupportForThisDriver driver not support
	ErrNotSupportForThisDriver = errors.New("not support for this driver")
)
