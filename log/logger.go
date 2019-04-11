package wxlog

import (
	"io"
	"fmt"
)

var msgWriter io.Writer

func SetLogger(logWriter io.Writer) {
	msgWriter = logWriter
}

func Logf(format string, a ...interface{}) (n int, err error) {
	if msgWriter != nil {
		return fmt.Fprintf(msgWriter, format, a...)
	}
	return 0, nil
}
