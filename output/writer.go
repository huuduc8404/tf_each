package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"tf_each/refactor"
)

func WriteFiles(resType string, refactored refactor.RefactoredBlock, tfvars refactor.TFVars) error {
	tfPath := filepath.Join("convert", fmt.Sprintf("%s.tf", resType))
	varPath := filepath.Join("convert", fmt.Sprintf("%s.tfvars", resType))

	err := os.WriteFile(tfPath, []byte(refactored.HCLContent), 0644)
	if err != nil {
		return err
	}

	varContent, err := json.MarshalIndent(tfvars, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(varPath, []byte(fmt.Sprintf("%s =  %s", resType, varContent)), 0644)
}
