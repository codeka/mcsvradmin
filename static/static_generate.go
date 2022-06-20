// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(http.Dir("./files"), vfsgen.Options{
		PackageName:  "static",
		BuildTags:    "!dev",
		VariableName: "Assets",
		Filename:     "static_vsfgen.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
