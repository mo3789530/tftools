package hcl

import (
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func Open(path string) *hcl.File {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	result, diags := hclsyntax.ParseConfig(data, path, hcl.Pos{})
	if diags.HasErrors() {
		panic(diags.Error())
	}
	return result
}
