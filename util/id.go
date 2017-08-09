package util

import (
	"log"
	"strconv"
	"strings"
)

func GetIdFrom(str, sep string) (uint64, error) {
	args := strings.Split(str, sep)
	var (
		uintId uint64
		e      error
	)
	if len(args) == 2 {
		uintId, e = strconv.ParseUint(args[1], 10, 64)
		if e != nil {
			log.Panic(e)
		}
	}
	return uintId, e
}

func GetStrIdFrom(str, sep string) (string, error) {
	args := strings.Split(str, sep)
	return args[1], nil
}
