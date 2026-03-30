package api

import (
//"encoding/json"
)

type resource[T any] struct {
	inner *T
}

func newResource[T any](inner *T) resource[T] {
	return resource[T]{
		inner: inner,
	}
}
