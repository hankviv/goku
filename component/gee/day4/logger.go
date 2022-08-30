package gee

import (
	"fmt"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		fmt.Fprintf(c.Writer, "Part A \n")
		c.Next()
		fmt.Fprintf(c.Writer, "Part B \n")
	}
}
