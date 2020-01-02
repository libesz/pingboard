// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(http.Dir("static"), vfsgen.Options{
		VariableName: "static",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
