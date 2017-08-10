package util

import (
	"bytes"
	"fmt"
)

func Concate(sep string, args ...interface{}) string {
	var b bytes.Buffer
	for i, a := range args {
		b.WriteString(fmt.Sprint(a))
		if i != len(args)-1 {
			b.WriteString(sep)
		}
	}
	return b.String()
}
