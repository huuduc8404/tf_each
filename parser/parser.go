package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Resource struct {
	Type    string
	Name    string
	Content *hclsyntax.Body
}

func ExtractResources(path string) ([]Resource, error) {
	parser := hclparse.NewParser()

	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("HCL parsing errors: %s", diags.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("unexpected body type")
	}

	var resources []Resource
	for _, block := range body.Blocks {
		if block.Type == "resource" && len(block.Labels) == 2 {
			resType := block.Labels[0]
			resName := block.Labels[1]
			resources = append(resources, Resource{
				Type:    resType,
				Name:    resName,
				Content: block.Body,
			})
		}
	}
	return resources, nil
}
