package hcl

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/hashicorp/hcl/v2/hclwrite"
	hcl2json "github.com/tmccombs/hcl2json/convert"
	"github.com/zclconf/go-cty/cty"
)

func HclBytesToJson(hclBytes []byte) {

	res, err := hcl2json.Bytes(hclBytes, "", hcl2json.Options{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(res))

}

func HclStringToJson(hcl string) {
	HclBytesToJson([]byte(hcl))
}

func JsonToHcl(jsonStr string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f := hclwrite.NewFile()
	rootBody := f.Body()
	for k, v := range data {
		WriteBodyBodyHcl(rootBody, k, v.([]interface{}))
	}

	// fmt.Println(string(f.Bytes()))
	return string(f.Bytes())
}

func FindMatchingBlocks(b *hclwrite.Body, name string, labels []string) []*hclwrite.Block {
	var matched []*hclwrite.Block
	for _, block := range b.Blocks() {
		if name == block.Type() {
			labelsName := block.Labels()
			if len(labels) == 0 && len(labelsName) == 0 {
				matched = append(matched, block)
				continue
			}
			if reflect.DeepEqual(labels, name) {
				matched = append(matched, block)
			}
		}
	}
	return matched
}

// write block
func WriteBodyBodyHcl(body *hclwrite.Body, key string, values []interface{}) {
	if len(values) < 1 {
		return
	}

	block := body.AppendNewBlock(key, []string{})
	blockBody := block.Body()
	for _, value := range values {
		WriteAttribute(blockBody, value)
	}

}

// write attribute
func WriteAttribute(blockBody *hclwrite.Body, value interface{}) {
	switch vt := value.(type) {
	case map[string]interface{}:
		for k, v := range vt {
			switch vv := v.(type) {
			case []interface{}:
				WriteBodyBodyHcl(blockBody, k, vv)
			default:
				blockBody.SetAttributeRaw(k, hclwrite.TokensForValue(cty.StringVal(fmt.Sprintf("%v", vv))))
			}
		}
	case []interface{}:
		for _, v := range vt {
			WriteAttribute(blockBody, v)
		}
	default:
	}
}
