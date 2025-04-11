package main

import (
	"fmt"
	"os"
	"tf_each/output"
	"tf_each/parser"
	"tf_each/refactor"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: tf_each <input.tf>")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Printf("File does not exist: %s\n", inputPath)
		os.Exit(1)
	}

	resources, err := parser.ExtractResources(inputPath)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	grouped := refactor.GroupResourcesByType(resources)

	err = os.MkdirAll("convert", os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create convert/ directory: %v\n", err)
		os.Exit(1)
	}

	for resType, group := range grouped {
		refactored, tfvars := refactor.RefactorGroup(resType, group)
		err := output.WriteFiles(resType, refactored, tfvars)
		if err != nil {
			fmt.Printf("Failed writing output for %s: %v\n", resType, err)
		}
	}

	fmt.Println("Conversion completed. Files generated in convert/")
}
