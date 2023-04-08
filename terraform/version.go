package terraform

import (
	"fmt"
	"log"
	"os"

	myhcl "github.com/mo3789530/tftools/hcl"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// https://github.com/orgrim/carcass/blob/fdbbcb85bf4e7628c7dd2c6f5a2b3814726722a6/terraform/config.go
func GetVersions(filepath string) (string, error) {
	src, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("error read file: %s", err)
		return "", err
	}

	file, diags := hclwrite.ParseConfig(src, filepath, hcl.InitialPos)
	if diags.HasErrors() {
		log.Printf("error hcl parse : %s", diags.Errs()[0])
		return "", diags.Errs()[0]
	}

	for _, tf := range myhcl.FindMatchingBlocks(file.Body(), "terraform", []string{}) {
		if tf.Body().GetAttribute("required_version") != nil {
			v := tf.Body().GetAttribute("required_version").Expr().BuildTokens(nil).Bytes()
			return string(v), nil
		}
	}
	return "", fmt.Errorf("error not found terrform block")
}
