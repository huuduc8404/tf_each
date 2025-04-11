package refactor

import (
	"fmt"
	"tf_each/parser"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type RefactoredBlock struct {
	HCLContent string
}

type TFVars map[string]interface{}

func GroupResourcesByType(resources []parser.Resource) map[string][]parser.Resource {
	grouped := make(map[string][]parser.Resource)
	for _, res := range resources {
		grouped[res.Type] = append(grouped[res.Type], res)
	}
	return grouped
}

func RefactorGroup(resType string, group []parser.Resource) (RefactoredBlock, TFVars) {
	var hcl string
	vars := make(TFVars)
	hcl += fmt.Sprintf("variable \"%s\" {\n  type = map(object({\n", resType)

	example := group[0]
	// Collect keys/types from the first resource as structure template
	for attrName, attr := range example.Content.Attributes {
		hcl += fmt.Sprintf("    %s = %s\n", attrName, inferType(attr.Expr))
	}

	hcl += "  }))\n}\n\n"

	hcl += fmt.Sprintf("resource \"%s\" \"generated\" {\n", resType)
	hcl += fmt.Sprintf("  for_each = var.%s\n", resType)
	for attrName := range example.Content.Attributes {
		hcl += fmt.Sprintf("  %s = each.value.%s\n", attrName, attrName)
	}
	hcl += "}\n"

	for _, res := range group {
		attrMap := make(map[string]interface{})
		for attrName, attr := range res.Content.Attributes {
			val := exprToLiteral(attr.Expr)
			// fmt.Printf("%s: %v\n", attrName, val)
			attrMap[attrName] = val
		}
		vars[res.Name] = attrMap
	}

	return RefactoredBlock{HCLContent: hcl}, vars
}

func inferType(expr hclsyntax.Expression) string {
	switch expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		typ := expr.(*hclsyntax.LiteralValueExpr).Val.Type().FriendlyName()
		return typ
	default:
		return "string"
	}
}

var evalCtx = &hcl.EvalContext{}

func exprToLiteral(expr hclsyntax.Expression) interface{} {
	val, diag := expr.Value(evalCtx)
	if diag.HasErrors() {
		return fmt.Sprintf("${%s}", expr.Range().String())
	}
	return ctyToGoValue(val)
}
func ctyToGoValue(val cty.Value) interface{} {
	switch val.Type() {
	case cty.String:
		return val.AsString()
	case cty.Number:
		n, _ := val.AsBigFloat().Float64()
		return n
	case cty.Bool:
		return val.True()
	default:
		if val.Type().IsMapType() || val.Type().IsObjectType() {
			m := make(map[string]interface{})
			for k, v := range val.AsValueMap() {
				m[k] = ctyToGoValue(v)
			}
			return m
		} else if val.Type().IsTupleType() || val.Type().IsListType() {
			list := val.AsValueSlice()
			arr := make([]interface{}, len(list))
			for i, v := range list {
				arr[i] = ctyToGoValue(v)
			}
			return arr
		}
		return val.GoString()
	}
}
