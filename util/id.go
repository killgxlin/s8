package util

import (
	"strings"
)

func GetStrIdFrom(str, sep string) (string, error) {
	strs := strings.Split(str, sep)
	return strs[1], nil
}
