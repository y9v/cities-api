package ds

import (
	"fmt"
)

type InvalidDataError struct {
	BucketName []byte
	key        string
	val        string
}

func (e InvalidDataError) Error() string {
	return fmt.Sprintf(
		"Invalid data in bucket %v at key %v: %v",
		string(e.BucketName), e.key, e.val,
	)
}
