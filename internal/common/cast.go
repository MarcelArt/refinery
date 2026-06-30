package common

import (
	"github.com/devfeel/mapper"
)

func Cast[T any](src any) (T, error) {
	var dst T

	err := mapper.AutoMapper(&src, &dst)
	return dst, err
}
