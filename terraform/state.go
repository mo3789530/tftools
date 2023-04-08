package terraform

import (
	"encoding/json"
	"fmt"
	"log"

	myhcl "github.com/mo3789530/tftools/hcl"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func ShowStateFileRaw(state *tfjson.State) string {
	file := hclwrite.NewEmptyFile()
	body := file.Body()

	writeBodyHcl(body, state.Values.RootModule.Resources)
	for _, child := range state.Values.RootModule.ChildModules {
		for _, grandChild := range child.ChildModules {
			writeBodyHcl(body, grandChild.Resources)
		}
		writeBodyHcl(body, child.Resources)
	}
	// fmt.Println(string(file.Bytes()))
	return string(file.Bytes())
}

func writeBodyHcl(body *hclwrite.Body, resources []*tfjson.StateResource) {
	for _, v := range resources {
		// write comment
		commentToken := hclwrite.Token{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte(fmt.Sprintf("# %v", v.Address)),
		}
		body.AppendUnstructuredTokens(hclwrite.Tokens{&commentToken})
		body.AppendNewline()
		// write resource
		block := body.AppendNewBlock("resource", []string{v.Type, v.Name})
		blockBody := block.Body()
		for k, attr := range v.AttributeValues {
			if attr == nil {
				continue
			}
			switch v := attr.(type) {
			case []interface{}:
				myhcl.WriteBodyBodyHcl(blockBody, k, v)
			case map[string]interface{}:
				if len(v) < 1 {
					// empty objects
					blockBody.SetAttributeRaw(k, hclwrite.Tokens{
						&hclwrite.Token{
							Bytes: []byte(fmt.Sprintf("%v", "{}")),
							Type:  hclsyntax.TokenAnd,
						},
					})
				} else {
					/*
						sample output
						tags_all = {
						    aaa = "bbb"
						}
					*/
					for kk, vv := range v {
						blockBody.SetAttributeValue(k, cty.ObjectVal(map[string]cty.Value{
							kk: cty.StringVal(fmt.Sprintf("%v", vv)),
						}))
					}
				}

			default:
				if v != "" {
					blockBody.SetAttributeRaw(k, hclwrite.TokensForValue(cty.StringVal(fmt.Sprintf("%v", v))))
				}
			}

		}
		body.AppendNewline()
	}
}

func ShowStateFileJson(state *tfjson.State) string {
	jsonByte, err := json.Marshal(state)
	if err != nil {
		log.Printf("error convert to json %s", err)

	}
	// fmt.Println(string(jsonByte))

	return string(jsonByte)
}
