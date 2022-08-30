package main

import (
	"fmt"
	gee "github.com/hankviv/goku/component/gee/day4"
	"log"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		name := c.Query("name")
		c.Data(200, []byte(fmt.Sprintf("hello world  Query Name:%s\n", name)))
	})
	r.Use(gee.Logger())
	r.Use(gee.Recover())
	log.Fatal(r.Run(":9999"))
}
